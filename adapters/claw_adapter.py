import ctypes
import os
import re

class VajraClawAdapter:
    def __init__(self, core_dll_path: str, static_rule_path: str):
        # 載入 C-Shared 二進位晶片
        try:
            self.lib = ctypes.CDLL(core_dll_path)
            print(f"[ClawAdapter] 成功掛載二進位晶片: {core_dll_path}")
        except OSError:
            print(f"[ClawAdapter] 警告：無法載入 {core_dll_path}，請確認是否已使用 go build -buildmode=c-shared 編譯。")
            self.lib = None
            return

        # 定義 C-FFI 接口簽名
        self.lib.init_static_vajra.argtypes = [ctypes.c_char_p]
        self.lib.init_static_vajra.restype = ctypes.c_int
        
        self.lib.inject_ephemeral_rule.argtypes = [ctypes.c_char_p]
        self.lib.inject_ephemeral_rule.restype = ctypes.c_int
        
        self.lib.match_token_stream.argtypes = [ctypes.c_char_p]
        self.lib.match_token_stream.restype = ctypes.c_int
        
        self.lib.clear_ephemeral_rules.restype = None

        # 啟動時立即結晶化常駐鐵律
        self._init_static(static_rule_path)

    def _init_static(self, static_rule_path: str):
        if not self.lib: return
        path_bytes = static_rule_path.encode('utf-8')
        result = self.lib.init_static_vajra(path_bytes)
        if result != 1:
            raise RuntimeError("Vajra 晶片啟動失敗！")

    def intercept_and_inject(self, user_prompt: str) -> str:
        """
        過電第一、二步：攔截與注入
        使用正則表達式，提取出使用者強制指定的邊界限制（例如：「不要提到...」、「嚴禁...」）
        """
        if not self.lib: return user_prompt

        print("[ClawAdapter] 🔍 攔截器開始掃描使用者 Prompt...")
        
        # 簡單正則：捕捉使用者要求「嚴禁...」的句子
        match = re.search(r'(嚴禁|不要提到|禁止)(.+?)(，|。|$)', user_prompt)
        
        if match:
            restriction = match.group(2).strip()
            print(f"[ClawAdapter] ⚡ 抓取到動態指令限制：禁止出現『{restriction}』")
            # 打入底層晶片
            self.lib.inject_ephemeral_rule(restriction.encode('utf-8'))
            
            # 從原 prompt 中移除這段限制，避免干擾 LLM 本身的回答邏輯
            # 我們讓底層防禦晶片來負責阻擋，LLM 只需要專心回答
            user_prompt = user_prompt.replace(match.group(0), "")
        
        return user_prompt.strip()

    def stream_monitor(self, llm_output_token: str):
        """
        過電第三、四步：放行 LLM 與 執行期盯防
        攔截 LLM 吐出的 Token 或 Tool Call 請求。
        """
        if not self.lib: return

        # 呼叫 C-FFI 進行 O(1) 雙軌比對
        token_bytes = llm_output_token.encode('utf-8')
        is_safe = self.lib.match_token_stream(token_bytes)
        
        if is_safe == 0:
            # 物理熔斷
            self.cleanup()
            raise PermissionError(f"[VajraClaw-Exception] 🛑 物理熔斷！LLM 嘗試輸出違規字串：{llm_output_token}")

    def cleanup(self):
        """
        任務結束，清理 JIT 動態指針
        """
        if self.lib:
            self.lib.clear_ephemeral_rules()

# ==========================================
# 示範使用 (Mock Execution)
# ==========================================
if __name__ == "__main__":
    dll_path = os.path.abspath(os.path.join(os.path.dirname(__file__), "../core/vajra_claw.dll"))
    rules_path = os.path.abspath(os.path.join(os.path.dirname(__file__), "../rules/Vajra.md"))
    
    # 建立防禦陣列
    claw = VajraClawAdapter(dll_path, rules_path)
    
    # 模擬使用者輸入
    raw_input = "請解釋唯識學的八識。另外，嚴禁提到股票與期貨。"
    
    print("\n--- 執行期開始 ---")
    
    # 1. 攔截並注入，回傳清洗過的乾淨 prompt 給 LLM
    clean_prompt = claw.intercept_and_inject(raw_input)
    print(f"發送給 LLM 的乾淨 Prompt: {clean_prompt}")
    
    # 2. 模擬 LLM 吐出 Token (執行期盯防)
    llm_tokens = ["唯識宗", "認為", "前六識", "負責感知", "而股票", "投資", "是..."]
    
    try:
        for token in llm_tokens:
            print(f"LLM 輸出: {token}")
            claw.stream_monitor(token)
    except PermissionError as e:
        print(e)
    finally:
        # 3. 蒸發動態規則
        claw.cleanup()
