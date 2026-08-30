# 🛡️ DROS™ (Deterministic Runtime Operating System)
**The Execution Governance Standard for Agentic AI**

[![Marketplace Status](https://img.shields.io/badge/GitHub_Marketplace-Available-success)](https://github.com/marketplace/dros-engine)
[![License: 3-Tier Model](https://img.shields.io/badge/License-Standard_3--Tier_Model-blue.svg)](#)
[![Official Website](https://img.shields.io/badge/Website-dr--os.io-purple.svg?style=for-the-badge)](https://dr-os.io)
[![Patent Status](https://img.shields.io/badge/U.S._Patent-64%2F111%2C973-blue.svg)](#)
[![RFC-010 Standard](https://img.shields.io/badge/Standard-RFC--010_Draft-orange.svg)](#)

---

## 🛑 "If runtime needs intelligence, the system is already broken."
Prompt Engineering is dead when it comes to enterprise execution security. No matter how complex your System Prompt is, Jailbreaks and Prompt Injections will bypass natural language filters. Traditional API Gateways fail because they assume human-driven execution constraints apply to autonomous agents.

**DROS** is NOT a prompt wrapper. It is the **Execution Governance Standard**. We move intelligence to compile-time (`Vajra.md` / `demo_policy.yaml`) and enforce rules via a **Dual-Layer Deterministic Runtime Architecture**:
* **L1 ATR Semantic Radar**: Intercepts prompt injections and indirect RAG/web context contamination at the ingestion gateway.
* **L2 Vajra Enforcer**: An **O(1) Bitmap Runtime Engine** (`vajra_claw.dll` / `.so` / `.aar` / C-ABI dynamic library / Docker Gateway) enforcing strict syscall and tool execution boundaries.

If the Agent attempts an unauthorized action, DROS triggers a **Physical Fusing (Strict Fail-Closed)** in under 1 microsecond.

---

## 🏛️ The 6-Pillars Enterprise AI Trust Model (DROS-6P)

DROS-VajraClaw enforces six fundamental trust boundaries in real time at the C-ABI / FFI in-band execution layer:

1. **Principal (Identity)**: 3-tier PKI-signed `DrosIdentityToken (DIT)` and native W3C `did:key` (Ed25519) for unbypassable agent identity binding.
2. **Authorization (Deterministic)**: Immutable $\mathcal{O}(1)$ capability bitmaps mapping agent roles to execution vectors.
3. **Action Bound (Syscall Gate)**: Sub-microsecond (<500ns) binary interception enforcing hard physical OS limits.
4. **Policy Gate (Dynamic Control)**: Dynamic data redaction, Human-In-The-Loop (HITL), and ZKP-Lite zero-knowledge proofs.
5. **Audit Log (Non-repudiability)**: SHA-256 Merkle Hash Chain + Ed25519 signatures, fully compliant with EU AI Act Art. 12.
6. **Expiry / Revocation (Microsecond)**: Constant-time $\mathcal{O}(1)$ dynamic bitmap updates for microsecond-level revocation and instant HTTP 403 enforcement.

---

## 🔥 Product Lines & Deployment Editions

### 1. DROS VajraClaw Hacker Edition (Free Docker Gateway)
* **Target**: Individual developers, AI researchers, and local single-machine multi-agent workstations.
* **Included**: Standalone Docker Gateway (`localhost:8080`), W3C `did:key` identity binding, O(1) AST policy fusing, and persistent SHA-256 Merkle audit chains.
* **Capacity**: Binds 1 Host, Max 5 Concurrent Agents.
* **Repository**: [`Top-Celestial-Company-Ltd/DROS-VajraClaw-Hacker`](https://github.com/Top-Celestial-Company-Ltd/DROS-VajraClaw-Hacker)

### 2. DSH Community Plugin (`dsh-plugin-vajraclaw`)
* **Target**: DeepSeek Harness (DSH) users.
* **Included**: 100% Zero-Dependency pure ESM TypeScript plugin. Instant regex tool-call failsafe and local persistent JSONL hash audit chain.
* **License**: Apache 2.0 Open Source (Free for Community).
* **Repository**: [`Top-Celestial-Company-Ltd/dsh-dros-vajraclaw`](https://github.com/Top-Celestial-Company-Ltd/dsh-dros-vajraclaw)

### 3. DROS Startup Edition
* **Target**: Fast-growing AI startups, agentic SaaS MVP platforms, and small engineering teams.
* **Included**: C-ABI Native Libraries (`.dll` / `.so`), Multi-Language SDKs (Python, TypeScript, Go, Java, Rust), Ed25519 Signed Policies, and REST/MCP Gateway.
* **Capacity**: Binds 3 Machine Hosts, Max 10 Concurrent Agents per host (Total 30 Concurrent Agents).

### 4. DROS Enterprise Audit Edition
* **Target**: High-security corporate environments (FinTech, Healthcare, Semiconductor, Government).
* **Included**: C-ABI Zero-Copy Enforcement Kernel (<500ns), Hardware HSM / TPM Key Binding, Dynamic PII Redaction, EU AI Act Article 12 Court-Grade Merkle Audit, and Air-Gapped / Private Cloud Deployment.
* **Capacity**: Binds 15 Machine Hosts, Max 30 Concurrent Agents per host (Total 450 Concurrent Agents).

### 5. VajraAgent Sovereign Mesh & Source Buyout
* **Target**: National defense, sovereign AI infrastructure, tier-1 financial institutions.
* **Included**: Full C-ABI / Go / Rust source code escrow, centralized IAM Control Plane, K8s Helm charts, sub-millisecond global kill-switch OTA, and customized hardware accelerators.

---

## 💻 Zero-Pollution Integration (Python Example)

```python
from integrations.vajraclaw.runtime import VajraClaw

# 1. Mount the Signed Binary Policy
vc = VajraClaw("demo_policy.yaml")

# 2. Evaluate the Tool Call in O(1) Time
decision = vc.evaluate("execute_payment", {"amount": 500})

# 3. Deterministic Enforcement
if not decision:
    raise PermissionError(f"Action blocked by DROS Execution Kernel: {decision.reason}")
```

---

## 💰 Commercial Pricing & Licensing Matrix

| Tier | Pricing | Best For | Licensing & Governance Scope |
| :--- | :--- | :--- | :--- |
| **Hacker (Docker Gateway)** | **$0 (Permanent Free)** | Solo developers, researchers | Free License for Individuals (1 Host, 5 Concurrent Agents, W3C DID, O(1) AST Gate) |
| **DSH Plugin (TS Failsafe)** | **$0 (Apache-2.0)** | DSH IDE users | Zero-Dependency in-process regex failsafe & local JSONL hash chain |
| **Startup License** | **$2,990 / yr** | AI Startups, SaaS MVPs | 3 Hosts, 10 Agents/host, Native C-ABI SDKs, Ed25519 Binary Policies |
| **Enterprise Audit** | **$29,990 / yr** | FinTech, Healthcare, Corp | 15 Hosts, 30 Agents/host, EU AI Act Art. 12 Merkle Audit, HSM Binding, Air-Gapped |
| **Sovereign Mesh / Source** | **Custom Quote** | Defense, Critical Infra | Full Source Code Buyout, K8s Helm Mesh, Global Instant Kill-Switch OTA |

---

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

## ⚖️ Standard 3-Tier License & Patent Constitution

* **Core Enforcement Substrate**: Protected under U.S. Provisional Patent Application (**U.S. PPA No. 64/111,973, Patent Pending**). All commercial and enterprise rights reserved by Top-Celestial Company Ltd.
* **Community Client & Docker Gateway**: Free License for Individuals (Permanent free non-commercial license for solo developers on single workstations).
* **Benchmark Testbed**: Open Source under Apache 2.0 ([DROS-VEP Lite](https://github.com/Top-Celestial-Company-Ltd/DROS-VEP-lite)).

---
**Developed by DROS Labs / 康宸園有限公司 (Top-Celestial Company Ltd.)**  
*Official Portal: [https://dr-os.io](https://dr-os.io) | Contact: [service@dr-os.io](mailto:service@dr-os.io)*