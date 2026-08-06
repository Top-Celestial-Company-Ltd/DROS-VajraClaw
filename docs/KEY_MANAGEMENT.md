# DROS VajraClaw 密碼學金鑰管理與備份指南 (KEY_MANAGEMENT.md)

本文件定義 DROS / VajraClaw 安全治理政策二進位合約（`policy.bin`）在開發期、編譯期與部署期的密鑰生命週期管理規範。所有系統管理員與安全工程師必須遵循本指南以確保零信任邊界不被攻破。

> [!CAUTION]
> **核心安全聲明：原廠無後門與金鑰救援機制（No Backdoors Warning）**
> * **聯絡原廠也無法救援**：DROS / VajraClaw 嚴格遵循零信任安全原則，系統中**不包含任何原廠萬能金鑰或後門**。如果您遺失了私鑰種子（Seed Hex），**即使聯絡原廠也完全無法幫您復原私鑰，亦無法替您簽署新的 `policy.bin`**。
> * **唯一救贖路徑**：您必須手動執行「重建信任根」程序（詳見第 3 節情境 B）。這要求您必須仍保有對伺服器的最高管理員權限（SSH/Root），重新生成金鑰對並逐一更新所有執行端節點的公鑰配置。如果同時遺失私鑰與伺服器控制權，系統將永久鎖死，無法升級。

---

## 1. 核心概念：臨時金鑰與靜態金鑰

VajraClaw 編譯器（`cli.py`）支援以下兩種簽署模式：

| 模式 | 金鑰生成與生命週期 | 生產環境適用度 | 潛在安全性與運維風險 |
| :--- | :--- | :--- | :--- |
| **臨時金鑰 (Ephemeral Key)** | 每次編譯時於記憶體內隨機生成，編譯結束即刻銷毀。 | **🚨 嚴禁使用** | 每次編譯產出的驗證公鑰不同，會導致執行端 GuardVM 驗證失敗而觸發物理熔斷，除非手動重新部署 VM 公鑰錨點。 |
| **靜態金鑰 (Static Key)** | 使用固定的 Ed25519 32 位元組（256-bit）私鑰種子（Seed Hex）進行簽章。 | **🟢 生產推薦** | 公鑰固定部署於 VM 端，合約可平滑升級，僅需嚴密保護與備份該私鑰種子。 |

---

## 2. 靜態金鑰生成與備份方案

### 2.1 密鑰對生成 (Python)
請在離線、安全的管理員電腦上執行以下腳本生成 Ed25519 金鑰對：
```python
import nacl.signing
import binascii

# 生成隨機私鑰種子
signing_key = nacl.signing.SigningKey.generate()
seed_hex = binascii.hexlify(signing_key.encode()).decode('utf-8')
pub_hex = binascii.hexlify(signing_key.verify_key.encode()).decode('utf-8')

print(f"【請備份】私鑰種子 (32-byte Seed Hex): {seed_hex}")
print(f"【請配置】驗證公鑰 (PubKey Hex)      : {pub_hex}")
```

### 2.2 密鑰種子備份最佳實踐
1. **禁止提交至 Git**：確保任何包含 `seed_hex` 的設定檔或腳本已登載於 `.gitignore`，嚴禁推送到公共或私有代碼庫。
2. **冷備份儲存**：將私鑰種子作為「安全備忘錄（Secure Note）」存放在企業級加密密碼庫中（如 1Password, KeePass 或 Bitwarden）。
3. **金鑰分離部署**：
   * **私鑰種子**：存放在安全的簽署發佈機上，僅用於編譯與發行 `policy.bin`。
   * **驗證公鑰**：寫入執行端 VM / 邊緣設備的啟動配置，公鑰丟失不影響安全性。

---

## 3. 災難復原流程 (Disaster Recovery SOP)

### 🚨 情境 A：編譯主機損毀，但私鑰種子有備份
此情境下，執行端公鑰錨點無需變更，合約升級完全無感。

1. **環境重置**：在新編編譯機上安裝依賴：
   ```bash
   pip install pynacl pyyaml
   ```
2. **還原編譯**：自加密庫提取備份的 `Seed Hex`，執行確定性編譯：
   ```bash
   python cli.py build rules/policy_young.yaml -o policy.bin --key <YOUR_BACKED_UP_SEED_HEX>
   ```
3. **合約替換**：將新產出的 `policy.bin` 拷貝至執行端對應目錄。微內核會無縫通過校驗並載入新版合約。

### 🚨 情境 B：私鑰種子徹底遺失 (災難性重建)
此情境下必須重新生成金鑰，並手動更新所有 VM 執行端。

1. **生成新金鑰**：依據 2.1 節重新生成一組 `NEW_SEED` 與 `NEW_PUBKEY`。
2. **重編合約**：使用新私鑰種子編譯政策：
   ```bash
   python cli.py build rules/policy_young.yaml -o policy.bin --key <NEW_SEED>
   ```
3. **更新執行端公鑰**：
   * 登入 VM，將 `gemini_proxy.py` 或微內核設定檔中的舊驗證公鑰替換為 `NEW_PUBKEY`。
   * 將 `NEW_SEED` 存入加密備份庫。
4. **服務重啟與部署**：
   * 上傳新版 `policy.bin` 到執行端.
   * 重啟微內核/代理服務使新公鑰生效（如：`sudo systemctl restart dros-proxy.service`）。
