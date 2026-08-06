# 🏯 DROS-VajraClaw 商用核心 Agent 指引 (DROS-VajraClaw/.agent/README.md)

> [!CAUTION]
> **本地區域規範**
> 任何 AI Agent 在操作或修改本資料夾（`DROS-VajraClaw/`）底下的商用金剛執行期原始碼時，必須嚴格遵守本安全與連線規範。

---

## 🛑 本地區域開發戒律

### 1. 🈲 嚴禁提交敏感金鑰與憑證 (No Secrets in Git)
*   **安全防禦**：本目錄下涉及的任何 API Key、授權驗證私鑰、資料庫連線密碼，**一律嚴禁以明文形式寫入程式碼中**。
*   **引用途徑**：代碼中必須透過環境變數或專案根目錄的 `.env_secrets` 檔案進行動態載入。

### 2. 🈲 嚴禁放行 Proxy 監聽位址 (Localhost Binding Mandate)
*   安全閘道代理（`gemini_proxy.py`）作為網關，**必須強制綁定在 `127.0.0.1` (localhost)**。
*   **嚴禁**將 Host 修改為 `0.0.0.0`，防範外部網絡未經安全代理直接調用。

### 3. ⚖️ 確保推理模型與溫度的相容性 (Reasoning vs Standard)
*   修改 `gemini_proxy.py` 處理 OpenAI/Gemini 請求的模組時，必須保留對推理型模型（如 `deepseek-v4-pro`）的相容邏輯：
    *   **非推理模型**（標準聊天）：允許 DROS 設定極低溫度 `0.05` 或 `0.2`。
    *   **推理型模型**：如果調用的是思考型模型，應忽略或排除低溫度參數（設定 `temperature: 1` 或不傳送），防範 API 報 400 錯誤。

### 4. 🗄️ 資料庫連線合規
*   技術客服（`query_db_for_email`）的資料庫查詢，必須經由 Tailscale 連線至 MariaDB 的 3308 埠。所有資料庫連線在結束後**必須在 `finally` 區塊中執行 `conn.close()`**，避免連接池溢出。

---
*DROS-VajraClaw 商業安全守護 ── 金剛護盾，密鑰防漏。* 🛡️🔒⚙️
