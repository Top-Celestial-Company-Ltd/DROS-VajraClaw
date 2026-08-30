# 🛡️ DROS - 自主型 AI Agent 系統執行治理標準
## (DROS: Deterministic Runtime Operating System)

> **企業級自主 AI 系統不可或缺的底層執行治理架構。**

[English](README.md) | [繁體中文](README_zh.md)

DROS (確定性運行期作業系統 / Deterministic Runtime Operating System) 是一套針對 Agentic AI 的隱形、高可靠執行治理基礎設施。它完全獨立於 LLM 語意推理空間之外，物理性地駐留在 Agent 工具調用輸出與企業核心作業系統/資料庫之間。

---

## 1. 核心問題：機率型 AI 資安的必然潰敗

缺乏執行期硬控制的 AI Agent，猶如一把隨時可能走火的武器。傳統依賴「LLM-as-a-judge」或自然語言 Prompt Guard 的資安手法無法防禦企業風險：
*   **提示注入 (Prompt Injection)**：零日語意攻擊能輕易穿透自然語言防護層。
*   **TOCTOU 與延遲**：在運行期解析 JSON 或調用第三方模型會產生不可預測的延遲與競爭漏洞。
*   **缺乏可審計性**：無法從密碼學上證明 *為什麼* 一個提示工程封裝允許了某項危險特權操作。

---

## 2. 開發者體驗：策略即代碼 (Policy-as-Code via AI)

告別複雜繁瑣的 YAML/JSON 手工配置。透過 DROS，企業 CISO 或資安工程師只需以自然語言結合 Markdown 撰寫策略 (`Vajra.md`)。您甚至可以使用 ChatGPT 或 Claude 根據企業資安守則直接生成 `Vajra.md`。

一旦撰寫完成，**DROS Compiler** 會將人類可讀的 Markdown 轉譯為高度最佳化、具備 Ed25519 密碼學簽名的唯讀二進位工件 (`policy.bin`)。

---

## 3. DROS 解決方案：確定性作業系統層 (Deterministic OS Layer)

DROS 將安全智慧移至**編譯期 (Compile-Time)**，並在**運行期 (Runtime)** 透過 $\mathcal{O}(1)$ **確定性位元圖譜 (Bitmap)** 進行強制執行。

1.  **DROS Compiler**：撰寫 `Vajra.md` 策略即代碼，編譯為帶簽名的二進位策略檔 (`policy.bin`)。
2.  **DROS Engine (VajraClaw)**：嵌入式 C-FFI / JNI 核心，透過純位元運算（常數時間 $\mathcal{O}(1)$ 記憶體查表）驗證執行權限。**無 LLM 二次評估、無語意模糊、絕無繞過空間。**

### 鋼性預設拒絕保證 (Strict Fail-Closed)
DROS 遵循零信任架構運作。一旦 Agent 嘗試發起未經授權的操作、Ed25519 數位簽名不符或策略 Hash 受損，DROS 會直接在 OS 系統呼叫層切斷執行路徑（Panic）。寧可立即中止該進程，也絕不讓未授權的 Payload 接觸企業資料庫。

---

## 4. 企業級 AI 六大信任基石模型 (DROS-6P)

DROS-VajraClaw 於 C-ABI / FFI 帶內執行層實時強制執行六大基礎信任邊界：

1. **主體身分 (Principal)**：三層 PKI 簽發之 `DrosIdentityToken (DIT)`，實現不可偽造的 Agent 身分綁定。
2. **確定性授權 (Authorization)**：不可篡改之 $\mathcal{O}(1)$ 權限位元圖譜，精確將角色映射至執行向量。
3. **系統呼叫閘門 (Action Bound)**：亞微秒級 (<500ns) 二進位攔截，施加硬性物理邊界。
4. **動態策略控制 (Policy Gate)**：動態敏感資料遮蔽、人機協同 (HITL) 與 ZKP-Lite 零知識證明。
5. **不可篡改審計 (Audit Log)**：SHA-256 Merkle Hash 鏈結 ＋ Ed25519 簽名，完全符合歐盟 AI 法案 (EU AI Act) 第 12 條規範。
6. **微秒級撤銷 (Expiry/Revocation)**：常數時間 $\mathcal{O}(1)$ 動態位元圖譜更新，實現微秒級權限撤銷與即時 HTTP 403 阻斷。

---

## 📜 相關技術核心論文與實測驗證 (Technical Foundations & Benchmarks)

本專案之確定性執行治理、微秒級熔斷與密碼學存證機制，參考並延伸自以下核心技術論文與開源實測環境：

1. **核心架構與六大信任邊界 (Core Architecture)**:
   * **論文**: *DROS-6P: A Unified Deterministic Runtime Governance Architecture Closing the Six Fundamental Trust Boundaries of Enterprise AI Agents*
   * **Zenodo DOI**: [10.5281/zenodo.21833970](https://doi.org/10.5281/zenodo.21833970) | **記錄典藏**: [zenodo.org/records/21833970](https://zenodo.org/records/21833970)

2. **四層深度防禦架構 (Defense-in-Depth Model)**:
   * **論文**: *DROS 4-Layer Defense-in-Depth Architecture for Autonomous AI Workloads*
   * **Zenodo DOI**: [10.5281/zenodo.21903475](https://doi.org/10.5281/zenodo.21903475) | **記錄典藏**: [zenodo.org/records/21903475](https://zenodo.org/records/21903475)

3. **外掛 FFI 與不可否認存證模組 (Runtime Attribution Framework)**:
   * **論文**: *Runtime Attribution Framework: An External C-ABI and PKI-Based Zero-Trust Infrastructure for Non-Repudiable Execution Governance in Multi-Agent Systems*
   * **Zenodo DOI**: [10.5281/zenodo.21903687](https://doi.org/10.5281/zenodo.21903687) | **記錄典藏**: [zenodo.org/records/21903687](https://zenodo.org/records/21903687)

4. **開源技術標準與實測基準倉 (Open Standard & Verification Sandbox)**:
   * **RFC-010 規範**: 遵循開放 Agent 身分與存證規範（W3C DID did:key 與 Ed25519 簽章鏈）。
   * **實測基準環境**: [DROS-VEP Lite (可復現安全評測沙盒)](https://github.com/Top-Celestial-Company-Ltd/DROS-VEP-lite)
   * **實測報告**: 涵蓋 24 小時長效多場景測試數據（160,611 次請求驗證，決策延遲 26.1μs）。

## 🚀 多場景整合與部署指南 (Multi-Scenario Guide)

### 🌟 情境 A：DSH (DeepSeek Harness) 沙盒使用者
1. **一鍵啟動 DROS Docker 網關**：
   ```bash
   docker run -d -p 8080:8080 --name dros-gateway dros/hacker-gateway:v1.0.0
   ```
2. **在 DSH 中安裝社區外掛**：
   ```bash
   dsh plugin --profile web add dsh-plugin-dros
   ```
3. **享受即時微秒級防禦**：DSH 內的 Agent 工具調用將立即由 DROS 執行期微核心進行確定性防護與審計。

---

### 💻 情境 B：Antigravity 2.0 / Codex / Cursor 開發者 (MCP 協議)
在您的 `mcp_settings.json` 或 Claude Desktop 配置中加入 DROS 網關：
```json
{
  "mcpServers": {
    "dros-governance": {
      "url": "http://localhost:8080/mcp",
      "transport": "http"
    }
  }
}
```

---

### 🐍 情境 C：原生 Python / LangChain / AutoGen 開發者
```python
from integrations.vajraclaw.runtime import VajraClaw

vc = VajraClaw("demo_policy.yaml")
decision = vc.evaluate("execute_payment", {"amount": 500})
if not decision:
    raise PermissionError(f"Blocked by DROS: {decision.reason}")
```

---

## 📜 技術白皮書與 Zenodo 權威學術引用

### 📚 Zenodo 官方同行評審論文
* 🏛️ **DROS 4-Layer Defense-in-Depth Architecture for Autonomous AI Workloads**
  * **DOI**: [`10.5281/zenodo.21755654`](https://doi.org/10.5281/zenodo.21755654)
* 🏛️ **DROS-6P: A Unified Deterministic Runtime Governance Architecture**
  * **DOI**: [`10.5281/zenodo.21808499`](https://doi.org/10.5281/zenodo.21808499)

### 📖 完整技術白皮書
* 📖 **[完整技術白皮書 (繁體中文 v2.5)](docs/DROS_AgenticWeb_Defense_Whitepaper_CN.md)**：*自主型 AI 工作負載的零信任執行治理*
* 📖 **[Full Whitepaper (English v2.5)](docs/DROS_AgenticWeb_Defense_Whitepaper_EN.md)**
* 🏛️ **[先前技術防禦宣告 (Defensive Publication)](docs/DEFENSIVE_PUBLICATION.md)**

---

> **“Linux 定義了機器如何運行軟體，而 DROS 定義了 AI Agent 被允許如何行動。”**

---
*專利聲明：DROS 執行治理與安全技術已申請美國臨時專利保護（U.S. Provisional Patent Application No. 64/111,973，Patent Pending）。*