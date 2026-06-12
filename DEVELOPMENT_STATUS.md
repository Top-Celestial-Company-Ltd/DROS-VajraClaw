# VajraClaw Core (Execution Governance Standard) - 開發紀錄與規格現況

> **最後更新時間**: 2026-05-31
> **對齊目標**: DROS V3 (Category Definition: Execution OS Layer)

## 📌 最新技術規格 (Specifications)
1. **核心語言與防禦機制**: 
   - 使用 Go 與 Zig C-FFI 開發，徹底放棄動態 JSON 解析，全面採用 **O(1) Bitmap Runtime Engine**。
   - 💡 **[ADR: 為什麼不留在 JSON？]** Agentic AI 需要絕對的執行確定性 (Determinism)。JSON 解析帶來的時間抖動會導致 TOCTOU 漏洞。O(1) Bitmap 讓我們能宣告 "Time-to-first-deny: < 1ms"。
2. **跨語言生態 (Ecosystem)**: 
   - **Python SDK**: 透過 `vajraclaw` 模組，無需理解底層指標，無縫包裹 `ProtectedTool`。
   - **Mobile SDK**: 透過 `gomobile`，輸出 Android `.aar` 確保邊緣 Edge AI 亦能遵守安全憲法。
3. **安全憲法 (Strict Fail-Closed)**: 
   - 不認得的 Tool -> DENY。
   - 憑證過期或簽章錯誤 -> 觸發致命中斷 (Panic) 並強制攔截。
   - 💡 **[ADR: 為什麼採用 Fail-Closed？]** 寧可讓正常的 Agent 卡住，也絕不放行一次危險的 Prompt Injection 攻擊。
4. **企業級合規與稽核 (Cryptographic Audit)**: 
   - 支援 `policy.bin` 格式 (Ed25519 簽章保護)。
   - 每一筆 `AuditTrace` 皆強制綁定 `PolicyHash` (SHA-256) 與 `CompilerVersion`。

## 🚀 開發進度現況 (Current Development Status)
- [x] **DROS Compiler & DSL v1 完成**: `vajra build` 決定性編譯 (Deterministic Build) 上線，相同的 `Vajra.md` 必然產出相同 Hash 的 `.bin`。
- [x] **Vajra Doctor & Linter 完工**: `[INFO]` 到 `[FATAL]` 的 5 級防護，將資安防線推向 Compile-time。
- [x] **Python C-FFI 與 Mobile JNI 完備**: 所有跨語言介面已全數打通。
- [x] **Free-Trial 機制 (Timebomb)**: 支援離線 30 天授權控制，為 PLG 成長鋪路。

## 📝 戰略定位與未來規劃 (Next Steps)
- DROS 已經從單一的 Security Tool，升級為 **"Execution Governance Standard"**。
- `dros-spec` (開源標準庫) 與 `dros-engine` (商業核心引擎) 的雙軌戰略已全面確立。
- 準備發布 DROS v1 Ecosystem Launch Pack (GitHub, Hacker News, Medium 白皮書)。
