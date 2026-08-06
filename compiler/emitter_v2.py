"""
DROS V2 (Project Aegis) — Vajra Compiler V2 DCT Emitter Extension

此模組擴充了 V1 emitter.py 的核心能力，在保留 V1 Binary 格式向後相容性的前提下，
新增對 Dynamic Capability Token (DCT) 前綴規則的編譯期解析與二進位封裝支援。

## 架構設計決策

### 設計決策 1：DSL 限制為 Prefix/Offset Matching Only

V2 強制限制策略語言的表達能力，僅允許：
```yaml
resources:
  - prefix: /data/tenant_{jwt.tenant_id}
  - prefix: /logs/user_{jwt.sub}
```

嚴格禁止任何執行期表達式（理由見下方）：
```yaml
# 禁止：這會讓決策時間從 O(1) 退化成不確定
conditions:
  - jwt.role == "manager"     # DENY
  - resource.cost < 100       # DENY
  - A && B                    # DENY（引入邏輯求值器）
  - A || B                    # DENY（引入邏輯求值器）
```

核心原因：DROS 的商業價值在於 Compile-Time → Bitmap → O(1) → Deterministic 的物理鏈。
引入任何「執行期求值」都會斷裂這條物理鏈，讓 DROS 退化成普通的策略引擎。

### 設計決策 2：不引入 JSON 解析器

DCT 在 Runtime 不解析 JSON。編譯器預先計算 JWT Claims 的「屬性偏移量表」，
Runtime 只執行「偏移量查表 + 記憶體切片比較」（等價於一個 u32 陣列的索引）。

## V2 policy.bin 附加區段格式

在 V1 的二進位格式後方，追加 DCT 擴充區段：
```
[V1 Binary Block (完整保留)] ← 向後相容，V1 Runtime 可忽略後方區段
[4 bytes]  DCT Magic: b'DCTX'
[4 bytes]  DCT Rule Count (uint32, big-endian)
[N * DCT_RULE_SIZE bytes]  DCT Rule Records
```

每條 DCT Rule Record 格式 (DCT_RULE_SIZE = 72 bytes)：
```
[64 bytes]  prefix (null-padded, 64-byte aligned)
[4 bytes]   attr_offset (uint32, big-endian)
[1 byte]    effect (0=DENY, 1=ALLOW)
[3 bytes]   padding
```
"""

import struct
import hashlib
from typing import List, Dict, Tuple, Optional

# --- V2 DSL 支援的動態佔位符前綴 ---
# 這些是允許在 DSL 中使用的 JWT Claims 佔位符
# 索引對應到 ClaimsTable 中的 attr_offset
KNOWN_JWT_PLACEHOLDERS: Dict[str, int] = {
    "jwt.tenant_id": 1,
    "jwt.sub":       2,
    "jwt.user_id":   3,
    "jwt.group_id":  4,
    "jwt.org_id":    5,
    "jwt.scope":     6,
}

DCT_MAGIC = b'DCTX'
DCT_RULE_SIZE = 72  # 64 (prefix) + 4 (attr_offset) + 1 (effect) + 3 (padding)


class DctLinterError(Exception):
    """DCT 策略語法錯誤"""
    pass


class DctLinterWarning:
    """DCT 策略潛在風險警告"""
    def __init__(self, msg: str):
        self.msg = msg


def parse_dct_prefix(resource_pattern: str) -> Tuple[bytes, int]:
    """
    解析 DSL 中的動態資源前綴宣告，提取靜態前綴與屬性偏移量。
    
    輸入範例:
        "/data/tenant_{jwt.tenant_id}/*"
    
    輸出:
        (b"/data/tenant_", 1)  # (靜態前綴, attr_offset)
    
    Args:
        resource_pattern: 包含 {jwt.xxx} 佔位符的資源路徑模式
    
    Returns:
        (prefix_bytes, attr_offset)
        - prefix_bytes: 前綴的 UTF-8 位元組，最長 64 bytes
        - attr_offset: 對應 JWT Claim 的偏移量索引（0 表示無動態驗證）
    
    Raises:
        DctLinterError: 若語法不符合 V2 DSL 限制
    """
    # 檢查是否包含任何被禁止的邏輯運算子
    FORBIDDEN_OPERATORS = ["&&", "||", "==", "!=", "<", ">", "<=", ">="]
    for op in FORBIDDEN_OPERATORS:
        if op in resource_pattern:
            raise DctLinterError(
                f"[DCT FATAL] 禁止在資源前綴中使用邏輯運算子 '{op}'。"
                f"\n  錯誤輸入: {resource_pattern!r}"
                f"\n  V2 只允許 Prefix/Offset Matching。請參閱 Vajra DSL Spec V2 第 2.5 節。"
            )

    # 尋找佔位符 {jwt.xxx}
    import re
    placeholder_pattern = re.compile(r'\{([^}]+)\}')
    matches = list(placeholder_pattern.finditer(resource_pattern))
    
    if not matches:
        # 無動態佔位符，純靜態前綴
        prefix_raw = resource_pattern.rstrip('/*').rstrip('/')
        prefix_bytes = prefix_raw.encode('utf-8')[:64]
        return prefix_bytes, 0
    
    if len(matches) > 1:
        raise DctLinterError(
            f"[DCT FATAL] V2 每條資源規則只允許一個動態佔位符，"
            f"但發現 {len(matches)} 個: {resource_pattern!r}"
        )
    
    placeholder = matches[0].group(1)  # e.g. "jwt.tenant_id"
    placeholder_start = matches[0].start()  # 佔位符在字串中的起始位置
    
    if placeholder not in KNOWN_JWT_PLACEHOLDERS:
        raise DctLinterError(
            f"[DCT FATAL] 未知的 JWT 佔位符: '{{{placeholder}}}'。"
            f"\n  允許的佔位符: {list(KNOWN_JWT_PLACEHOLDERS.keys())}"
        )
    
    attr_offset = KNOWN_JWT_PLACEHOLDERS[placeholder]
    
    # 靜態前綴 = 佔位符之前的部分（不含 {）
    static_prefix = resource_pattern[:placeholder_start]
    prefix_bytes = static_prefix.encode('utf-8')[:64]
    
    return prefix_bytes, attr_offset


def build_dct_rules(v2_resources: List[Dict]) -> Tuple[List[Dict], List[str]]:
    """
    從 V2 DSL 的 resources 區塊建立 DCT 規則列表。
    
    V2 DSL 資源區塊格式：
    ```yaml
    resources:
      - prefix: /data/tenant_{jwt.tenant_id}
        effect: ALLOW
      - prefix: /admin/
        effect: DENY
    ```
    
    Returns:
        (dct_rules, linter_warnings)
        - dct_rules: [{"prefix_bytes": bytes, "attr_offset": int, "effect": int}]
        - linter_warnings: 警告訊息列表
    """
    dct_rules = []
    warnings = []
    
    for i, resource in enumerate(v2_resources):
        pattern = resource.get("prefix", "")
        effect_str = resource.get("effect", "DENY").upper()
        
        if not pattern:
            warnings.append(f"[DCT WARN] resources[{i}] 缺少 'prefix' 欄位，已跳過。")
            continue
        
        try:
            prefix_bytes, attr_offset = parse_dct_prefix(pattern)
        except DctLinterError as e:
            # 嚴重錯誤：中止編譯
            raise
        
        if len(prefix_bytes) == 0:
            warnings.append(f"[DCT WARN] resources[{i}] 前綴 {pattern!r} 解析後為空，已跳過。")
            continue
        
        effect_byte = 1 if effect_str == "ALLOW" else 0
        
        dct_rules.append({
            "prefix_bytes": prefix_bytes,
            "attr_offset": attr_offset,
            "effect": effect_byte,
        })
    
    return dct_rules, warnings


def pack_dct_rule(rule: Dict) -> bytes:
    """
    將單條 DCT 規則序列化為 72-byte 二進位記錄。
    
    格式：
    - 64 bytes: 前綴（UTF-8，不足部分補 \\x00）
    - 4 bytes: attr_offset (uint32 big-endian)
    - 1 byte: effect (0=DENY, 1=ALLOW)
    - 3 bytes: padding \\x00
    """
    prefix_bytes = rule["prefix_bytes"]
    # 固定 64 bytes，不足補零
    padded_prefix = prefix_bytes[:64].ljust(64, b'\x00')
    attr_offset = struct.pack('>I', rule["attr_offset"])
    effect = struct.pack('>B', rule["effect"])
    padding = b'\x00' * 3
    
    result = padded_prefix + attr_offset + effect + padding
    assert len(result) == DCT_RULE_SIZE, f"DCT rule size mismatch: {len(result)} != {DCT_RULE_SIZE}"
    return result


def emit_dct_extension_block(dct_rules: List[Dict]) -> bytes:
    """
    生成 V2 的 DCT 擴充二進位區塊，附加在 V1 policy.bin 之後。
    
    格式：
    - 4 bytes: 'DCTX' Magic
    - 4 bytes: Rule Count (uint32 big-endian)
    - N * 72 bytes: DCT Rule Records
    """
    if not dct_rules:
        return b''
    
    block = bytearray()
    block.extend(DCT_MAGIC)
    block.extend(struct.pack('>I', len(dct_rules)))
    
    for rule in dct_rules:
        block.extend(pack_dct_rule(rule))
    
    return bytes(block)


def parse_v2_extensions(policy: dict) -> Tuple[List[Dict], List[str], List[str]]:
    """
    從 YAML 策略文件中提取 V2 DCT 資源宣告。
    
    Returns:
        (dct_rules, errors, warnings)
    """
    errors = []
    warnings = []
    dct_rules = []
    
    # V2 擴充欄位: resources (動態前綴資源規則)
    v2_resources = policy.get("resources", [])
    
    if not v2_resources:
        # 沒有 resources 區塊是允許的（V1 相容模式）
        return [], [], []
    
    try:
        dct_rules, warn = build_dct_rules(v2_resources)
        warnings.extend(warn)
    except DctLinterError as e:
        errors.append(str(e))
    
    return dct_rules, errors, warnings


def emit_binary_v2_extended(
    v1_binary: bytes,
    dct_rules: List[Dict]
) -> bytes:
    """
    在 V1 二進位之後附加 DCT 擴充區段，生成 V2 複合二進位。
    
    V2 的核心向後相容設計：
    - V1 Runtime 讀取二進位，到達 V1 payload 末尾後停止，不感知後方的 DCTX 區塊
    - V2 Runtime 繼續讀取 DCTX Magic，載入 DCT 規則表
    - 這確保了 V1/V2 Runtime 可以讀取同一份 policy.bin
    
    Args:
        v1_binary: 由 emit_binary_v3() 生成的 V1 格式二進位
        dct_rules: 由 parse_v2_extensions() 解析出的 DCT 規則列表
    
    Returns:
        v1_binary + dct_extension_block（若無 DCT 規則則直接回傳 v1_binary）
    """
    if not dct_rules:
        return v1_binary
    
    dct_block = emit_dct_extension_block(dct_rules)
    return v1_binary + dct_block


def decode_dct_block(binary: bytes, v1_payload_end: int) -> Optional[List[Dict]]:
    """
    從 V2 複合二進位中解析 DCT 擴充區塊（用於驗證與 Runtime 載入）。
    
    Args:
        binary: 完整的 V2 policy.bin 二進位
        v1_payload_end: V1 Payload 結束的位元組偏移量
    
    Returns:
        DCT 規則列表，若無 DCTX 區塊則回傳 None
    """
    if len(binary) <= v1_payload_end + 4:
        return None  # 沒有足夠空間容納 DCTX Magic
    
    magic = binary[v1_payload_end:v1_payload_end + 4]
    if magic != DCT_MAGIC:
        return None  # 不是 DCTX 區塊
    
    offset = v1_payload_end + 4
    rule_count = struct.unpack_from('>I', binary, offset)[0]
    offset += 4
    
    rules = []
    for _ in range(rule_count):
        if offset + DCT_RULE_SIZE > len(binary):
            break
        prefix_bytes = binary[offset:offset + 64].rstrip(b'\x00')
        attr_offset = struct.unpack_from('>I', binary, offset + 64)[0]
        effect = binary[offset + 68]
        rules.append({
            "prefix_bytes": prefix_bytes,
            "attr_offset": attr_offset,
            "effect": effect,
        })
        offset += DCT_RULE_SIZE
    
    return rules
