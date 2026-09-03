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


---

## 📝 如何設定安全策略？(How to Configure Vajra.md)

DROS 支援兩種極簡設定方式：**人類直覺 Markdown 格式 (`Vajra.md`)** 與 **結構化 YAML 格式 (`demo_policy.yaml`)**。

### 1. 📄 人類直覺寫法範例 (`Vajra.md`)
只需以白話 Markdown 宣告允許執行的白名單與防禦邊界：

```markdown
# 🛡️ DROS Agent 安全策略規範 (Vajra.md)

## 1. 允許執行的工具 (Allowed Capabilities)
- 允許讀取當前工作區檔案 (`file_read`)
- 允許執行一般查詢 (`search_web`, `query_db`)
- 允許終端執行唯讀指令 (`git status`, `npm test`, `cargo check`)

## 2. 嚴格禁止的邊界 (Strict Fail-Closed Boundaries)
- 禁止執行任何遞迴刪除或清空指令 (`rm -rf`, `rmdir /s`, `format`)
- 禁止存取敏感憑證檔案 (`.env`, `id_rsa`, `secrets.json`, `.aws/credentials`)
- 禁止單筆交易金額超過 1,000 元 (`amount <= 1000`)
```

---


> [!CAUTION]
> 🔒 **極核心元權限警語：人類主權剛性封印與作業系統唯讀硬鎖定！**
> **絕對禁止讓 AI Agent 自行生成或寫入「禁止修改 Vajra.md」這條指令！** 若 AI 擁有編輯權限，一旦遭到提示詞注入劫持，隨時可以在攻擊前先調用工具將該規則刪除或註釋掉（自指迴圈漏洞）。
> **鐵律：** 保護規則檔本身的剛性邊界，在邏輯上**必須由人類開發者親手貼入**，並在終端執行作業系統唯讀鎖死：
> - **人類封印**：在檔案頂部貼入 `<!-- [HUMAN SEALED] -->` 與 `DENY WRITE/DELETE path: "**/vajra.md"`
> - **Linux / macOS**: `chmod 444 Vajra.md`
> - **Windows (PowerShell)**: `attrib +r Vajra.md`（或 `icacls Vajra.md /deny "Users:(W)"`）
> - **Docker 掛載時**: 使用唯讀掛載模式 `-v $(pwd)/Vajra.md:/app/demo_policy.yaml:ro`
> 
> 📖 **完整教學**：請參閱 [Hacker Edition (個人社群免費版) 長任務開發者手冊](https://github.com/Top-Celestial-Company-Ltd/DROS-VEP-lite/blob/main/docs/guides/DROS_HACKER_EDITION_MANUAL.md)，內含 4 步曲手動設定、越權三重感官警報與暫時解封 SOP。


### 2. 🤖 讓 AI 幫你一秒生成策略！(AI Prompt Template)

您不需要從零手寫！直接將以下**「萬用提示詞 (Prompt)」**複製給 ChatGPT、Claude 或 Cursor，AI 就會自動產出標準合規的 `Vajra.md`：

> 📋 **複製這段 Prompt 給任何 LLM / Agent：**
> 
> ```text
> 你現在是 DROS 確定性安全架構專家。請根據我的 Agent 角色，為我生成一份標準的 DROS「Vajra.md」安全策略 Markdown 檔案。
> 
> 我的 Agent 需求如下：
> - Agent 角色與場景：【例如：全端工程師 / 客服機器人 / 自動化財務助理】
> - 允許的工具與操作：【例如：讀寫代碼、執行 npm test、查詢訂單資料庫】
> - 嚴格禁止的邊界：【例如：禁止刪除根目錄、禁止讀取 .env、單次轉帳上限 500】
> 
> 請遵循 DROS「預設拒絕 (Default Fail-Closed)」白名單原則，生成清晰的 Markdown 規則區塊，包含：
> 1. 角色定義與授權範疇 (Role & Scope)
> 2. 白名單工具 (Allowed Capabilities)
> 3. 邊界條件約束 (Thresholds & Security Patterns)
> ```

---

### 3. 🔄 策略即時熱更新 (Hot Reloading)
啟動 Docker 網關時，只需將您的 `Vajra.md` 掛載進去，修改存檔後 **1 微秒內即時生效，無需重啟容器**：
```bash
docker run -d -p 8080:8080 --name dros-gateway \
  -v $(pwd)/Vajra.md:/app/demo_policy.yaml \
  dros/hacker-gateway:v1.0.0
```


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

## 🛡️ 治理與防禦能力對照矩陣 (Defense Capability Matrix)

| 威脅防禦維度 / 核心能力 | 傳統 LLM 防護 (NeMo / 提示詞審查) | 📦 DSH 獨立 TypeScript 外掛 | 🛡️ DROS Hacker Docker 網關 | 🏢 企業版 / K8s 集群 |
| :--- | :---: | :---: | :---: | :---: |
| **運行載體 (Vehicle)** | 雲端 API / 外部大模型 | 進程內原生 JS (零外部依賴) | **本地 Docker 容器 (`:8080`)** | 企業集群 / K8s / C-ABI 微核心 |
| **保護範圍 (Scope)** | 單次對話 Session | DSH 單一本機進程 | **全生態 (Claude+Codex+Cursor+DSH+AGY)** | 跨主機節點集群 / 私有雲 |
| **執行意圖治理 (Governance)** | 🔴 僅限文字模糊比對 | 🟢 **正則表達式硬防線 (Regex Failsafe)** | 🟢 **100% 確定性 AST 語法樹熔斷 (<1µs)** | 🟢 **AST 點陣圖 ＋ eBPF 內核級攔截** |
| **破壞性指令攔截 (Destructive)** | 🔴 易遭提示注入與編碼繞過 | 🟢 **敏感路徑物理阻斷** | 🟢 **底層 Syscall 物理硬熔斷** | 🟢 **硬體 HSM 隔離 ＋ 內核檔案鎖** |
| **機密與金鑰防洩漏 (Secrets)** | 🔴 無物理安全防線 | 🟢 **敏感路徑讀取阻斷** | 🟢 **動態遮蔽 ＋ 虛擬沙箱隔離** | 🟢 **硬體 HSM ＋ 零知識微證明 (ZKP)** |
| **Agent 主體身分綁定 (Identity)** | 🔴 無身分認證 | 🟢 Session 級別識別碼 | 🟢 **原生 W3C `did:key` (Ed25519 簽章)** | 🟢 **三層 PKI `DrosIdentityToken (DIT)`** |
| **不可篡改審計存證 (Audit)** | 🔴 普通可竄改文字 Log | 🟢 **本地 SHA-256 雜湊鏈** | 🟢 **Ed25519 簽章 Merkle 雜湊鏈** | 🟢 **歐盟 AI 法案第 12 條法院級存證** |
| **RFC-010 代理通行證 (Passport)** | 🔴 不支援 | 🟢 格式解析器 | 🟢 **本地發行 ＋ 跨 Agent 密碼學驗證** | 🟢 **跨組織漫遊通行證與權限繼承** |
| **決策延遲 (Decision Latency)** | 🔴 1,000 ~ 3,000 ms (二次模型極慢) | 🟢 **<1 ms (記憶體直接攔截)** | 🟢 **<1 µs (C-ABI) / <1 ms (REST 網關)** | 🟢 **<500 ns (零拷貝常數時間查表)** |
| **授權條款 (License)** | 按 Token 計費 | **100% 免費開源 (Apache-2.0)** | **個人永久免費授權 (Free for Individuals)** | 新創版 $2,990 / 企業版 $29,990 |
