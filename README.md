# 🛡️ DROS - Execution Governance Standard for Agentic Systems

> **The missing execution layer for autonomous AI systems.**

[English](README.md) | [繁體中文](README_zh.md)

DROS (Deterministic Runtime Operating System) is an invisible, military-grade execution governance infrastructure for Agentic AI. 
It operates completely outside of the LLM reasoning space, sitting physically between the Agent's output and your enterprise operating system.

## 1. The Core Problem: Probabilistic Security Fails

Without execution control, AI Agents are a loaded gun pointed at your infrastructure. Current "security" relies on LLM-as-a-judge or JSON parsers, which fail because:
*   **Prompt Injections**: Zero-day semantic attacks bypass natural language guards.
*   **TOCTOU & Latency**: JSON parsing at runtime creates unpredictable delays and vulnerabilities.
*   **Lack of Auditability**: You cannot cryptographically prove *why* a prompt engineering wrapper allowed an action.

## 2. Developer Experience: Policy-as-Code via AI

Forget complex YAML or JSON configurations. With DROS, your CISO or Security Engineer simply writes the policy in natural language mixed with Markdown (`Vajra.md`). 
You can even use ChatGPT or Claude to generate your `Vajra.md` based on your company's security playbook. 

Once written, the **DROS Compiler** transforms this human-readable Markdown into a highly optimized, cryptographically signed `policy.bin` artifact ready for production.

## 3. The DROS Solution: Deterministic OS Layer

DROS shifts intelligence to **compile-time** and enforces rules via an **O(1) deterministic bitmap** at **runtime**.

1.  **DROS Compiler**: Write `Vajra.md` Policy-as-Code. Compile it into a cryptographically signed binary artifact (`policy.bin`).
2.  **DROS Engine**: A C-FFI / JNI embedded engine that performs pure bitwise AND operations (O(1) memory lookup) to validate execution. **No LLM evaluation. No semantic interpretation. No bypass.**

### Strict Fail-Closed Guarantee
DROS operates on a Zero-Trust basis. If an Agent attempts an unauthorized action, if the Ed25519 signature doesn't match, or if the Policy Hash is corrupted, DROS physically severs the execution path at the OS level (Panic). It will rather crash the application than let an unverified payload touch your database.

## 4. The 6-Pillars Enterprise AI Trust Model (DROS-6P)

DROS-VajraClaw enforces six fundamental trust boundaries in real time at the C-ABI / FFI in-band execution layer:

1. **Principal (Identity)**: 3-tier PKI-signed `DrosIdentityToken (DIT)` for unbypassable agent identity binding.
2. **Authorization (Deterministic)**: Immutable $O(1)$ capability bitmaps mapping roles to execution vectors.
3. **Action Bound (Syscall Gate)**: Sub-microsecond (<500ns) binary interception enforcing hard physical limits.
4. **Policy Gate (Dynamic Control)**: Dynamic data redaction, Human-In-The-Loop (HITL), and ZKP-Lite zero-knowledge proofs.
5. **Audit Log (Non-repudiability)**: SHA-256 Merkle Hash Chain + Ed25519 signatures, fully compliant with EU AI Act Art. 12.
6. **Expiry / Revocation (Microsecond)**: Constant-time $O(1)$ dynamic bitmap updates for microsecond-level revocation and instant HTTP 403 enforcement.

---


---

## 📝 How to Configure Security Policies (Vajra.md Guide)

DROS supports two straightforward formats: **Intuitive Markdown (`Vajra.md`)** and **Structured YAML (`demo_policy.yaml`)**.

### 1. 📄 Intuitive Markdown Example (`Vajra.md`)
Declare allowed capabilities and hard security boundaries in plain Markdown:

```markdown
# 🛡️ DROS Agent Security Policy (Vajra.md)

## 1. Allowed Capabilities
- Allow reading workspace files (`file_read`)
- Allow standard queries (`search_web`, `query_db`)
- Allow safe terminal commands (`git status`, `npm test`, `cargo check`)

## 2. Strict Fail-Closed Boundaries
- Block all recursive deletion or wiping commands (`rm -rf`, `rmdir /s`, `format`)
- Block access to credential paths (`.env`, `id_rsa`, `secrets.json`, `.aws/credentials`)
- Restrict transaction amounts exceeding $1,000 threshold (`amount <= 1000`)
```

---


> [!IMPORTANT]
> 🔒 **Crucial Security Best Practice: Lock `Vajra.md` to Read-Only After Configuration!**
> To prevent compromised or hallucinating AI Agents from attempting to rewrite their own security rules to escalate privileges, **always set your policy file to read-only once configured**:
> - **Linux / macOS**: `chmod 444 Vajra.md`
> - **Windows (PowerShell)**: `Set-ItemProperty -Path Vajra.md -Name IsReadOnly -Value $true`
> - **Docker Container Mount**: Mount with the read-only flag `-v $(pwd)/Vajra.md:/app/demo_policy.yaml:ro`
> 
> *(Note: DROS kernel enforces 4-Layer Invariant Defense to intercept unauthorized policy modifications in-band; combining this with OS file-level locks achieves 100% airtight physical defense!)*


### 2. 🤖 Let AI Generate Your Policy in 1 Second! (AI Prompt Template)

You don't need to write policies from scratch! Copy the following universal prompt to ChatGPT, Claude, or Cursor:

> 📋 **Copy this Prompt to any LLM / AI Assistant:**
> 
> ```text
> You are a DROS deterministic security architecture expert. Based on my Agent requirements, generate a standard DROS "Vajra.md" security policy in Markdown.
> 
> Agent Details:
> - Agent Role & Scenario: [e.g., Fullstack Developer / Customer Service / Financial Automation]
> - Allowed Tools & Operations: [e.g., Read/Write src/, Run tests, Query order database]
> - Strict Boundaries & Denials: [e.g., Block deletion of root/workspace, Block .env access, Payment limit $500]
> 
> Follow the DROS "Default Fail-Closed" whitelist principle and structure the output into:
> 1. Role & Capability Scope
> 2. Allowed Capabilities (Whitelist)
> 3. Security Boundary Constraints (Thresholds & Pattern Failsafes)
> ```

---

### 3. 🔄 Instant Hot Reloading
Simply mount your `Vajra.md` when launching the Docker gateway. Policy changes take effect in **<1 microsecond without container restarts**:
```bash
docker run -d -p 8080:8080 --name dros-gateway \
  -v $(pwd)/Vajra.md:/app/demo_policy.yaml \
  dros/hacker-gateway:v1.0.0
```


## 📜 Technical Foundations & Benchmark Publications


---

## 📝 How to Configure Security Policies (Vajra.md Guide)

DROS supports two straightforward formats: **Intuitive Markdown (`Vajra.md`)** and **Structured YAML (`demo_policy.yaml`)**.

### 1. 📄 Intuitive Markdown Example (`Vajra.md`)
Declare allowed capabilities and hard security boundaries in plain Markdown:

```markdown
# 🛡️ DROS Agent Security Policy (Vajra.md)

## 1. Allowed Capabilities
- Allow reading workspace files (`file_read`)
- Allow standard queries (`search_web`, `query_db`)
- Allow safe terminal commands (`git status`, `npm test`, `cargo check`)

## 2. Strict Fail-Closed Boundaries
- Block all recursive deletion or wiping commands (`rm -rf`, `rmdir /s`, `format`)
- Block access to credential paths (`.env`, `id_rsa`, `secrets.json`, `.aws/credentials`)
- Restrict transaction amounts exceeding $1,000 threshold (`amount <= 1000`)
```

---


> [!IMPORTANT]
> 🔒 **Crucial Security Best Practice: Lock `Vajra.md` to Read-Only After Configuration!**
> To prevent compromised or hallucinating AI Agents from attempting to rewrite their own security rules to escalate privileges, **always set your policy file to read-only once configured**:
> - **Linux / macOS**: `chmod 444 Vajra.md`
> - **Windows (PowerShell)**: `Set-ItemProperty -Path Vajra.md -Name IsReadOnly -Value $true`
> - **Docker Container Mount**: Mount with the read-only flag `-v $(pwd)/Vajra.md:/app/demo_policy.yaml:ro`
> 
> *(Note: DROS kernel enforces 4-Layer Invariant Defense to intercept unauthorized policy modifications in-band; combining this with OS file-level locks achieves 100% airtight physical defense!)*


### 2. 🤖 Let AI Generate Your Policy in 1 Second! (AI Prompt Template)

You don't need to write policies from scratch! Copy the following universal prompt to ChatGPT, Claude, or Cursor:

> 📋 **Copy this Prompt to any LLM / AI Assistant:**
> 
> ```text
> You are a DROS deterministic security architecture expert. Based on my Agent requirements, generate a standard DROS "Vajra.md" security policy in Markdown.
> 
> Agent Details:
> - Agent Role & Scenario: [e.g., Fullstack Developer / Customer Service / Financial Automation]
> - Allowed Tools & Operations: [e.g., Read/Write src/, Run tests, Query order database]
> - Strict Boundaries & Denials: [e.g., Block deletion of root/workspace, Block .env access, Payment limit $500]
> 
> Follow the DROS "Default Fail-Closed" whitelist principle and structure the output into:
> 1. Role & Capability Scope
> 2. Allowed Capabilities (Whitelist)
> 3. Security Boundary Constraints (Thresholds & Pattern Failsafes)
> ```

---

### 3. 🔄 Instant Hot Reloading
Simply mount your `Vajra.md` when launching the Docker gateway. Policy changes take effect in **<1 microsecond without container restarts**:
```bash
docker run -d -p 8080:8080 --name dros-gateway \
  -v $(pwd)/Vajra.md:/app/demo_policy.yaml \
  dros/hacker-gateway:v1.0.0
```


## 📜 Technical Foundations & Benchmark Publications

The deterministic execution governance, microsecond fusing, and cryptographic audit mechanisms in this project are referenced from and build upon the following core technical papers and verification environments:

1. **Core Architecture & Six Trust Boundaries (Core Architecture)**:
   * **Paper**: *DROS-6P: A Unified Deterministic Runtime Governance Architecture Closing the Six Fundamental Trust Boundaries of Enterprise AI Agents*
   * **Zenodo DOI**: [10.5281/zenodo.21833970](https://doi.org/10.5281/zenodo.21833970) | **Archived Record**: [zenodo.org/records/21833970](https://zenodo.org/records/21833970)

2. **Defense-in-Depth Model (4-Layer Security)**:
   * **Paper**: *DROS 4-Layer Defense-in-Depth Architecture for Autonomous AI Workloads*
   * **Zenodo DOI**: [10.5281/zenodo.21903475](https://doi.org/10.5281/zenodo.21903475) | **Archived Record**: [zenodo.org/records/21903475](https://zenodo.org/records/21903475)

3. **Runtime Attribution & C-ABI Module (Attribution Framework)**:
   * **Paper**: *Runtime Attribution Framework: An External C-ABI and PKI-Based Zero-Trust Infrastructure for Non-Repudiable Execution Governance in Multi-Agent Systems*
   * **Zenodo DOI**: [10.5281/zenodo.21903687](https://doi.org/10.5281/zenodo.21903687) | **Archived Record**: [zenodo.org/records/21903687](https://zenodo.org/records/21903687)

4. **Open Standards & Verification Sandbox**:
   * **RFC-010 Specification**: Adheres to open Agent Identity & Attestation standard (W3C DID did:key & Ed25519 signature chain).
   * **Verification Sandbox**: [DROS-VEP Lite (Reproducible Evaluation Sandbox)](https://github.com/Top-Celestial-Company-Ltd/DROS-VEP-lite)
   * **Evaluation Metrics**: 24-hour soak benchmark results (160,611 verified requests, 26.1μs decision latency).

## 🚀 Multi-Scenario Deployment & Setup Guide

### 🌟 Scenario A: DSH (DeepSeek Harness) Sandbox Users
1. **Start the DROS Docker Gateway**:
   ```bash
   docker run -d -p 8080:8080 --name dros-gateway dros/hacker-gateway:v1.0.0
   ```
2. **Install DROS Community Plugin in DSH**:
   ```bash
   dsh plugin --profile web add dsh-plugin-dros
   ```
3. **Enjoy Zero-Friction Protection**: DSH Agents are immediately bound to microsecond $O(1)$ tool interception.

---

### 💻 Scenario B: Antigravity 2.0 / Codex / Cursor Developers (MCP Protocol)
Add the DROS Gateway to your `mcp_settings.json` / Claude Config:
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

### 🐍 Scenario C: Native Python / LangChain / AutoGen Developers
```python
from integrations.vajraclaw.runtime import VajraClaw

vc = VajraClaw("demo_policy.yaml")
decision = vc.evaluate("execute_payment", {"amount": 500})
if not decision:
    raise PermissionError(f"Blocked by DROS: {decision.reason}")
```

---

> **"Linux defined how machines run software. DROS defines how agents are allowed to act."**

*Developed by DROS Labs / 康宸園有限公司 (Top-Celestial Company Ltd.)*  
*Protected under U.S. Provisional Patent Application No. 64/111,973 (Patent Pending).*
