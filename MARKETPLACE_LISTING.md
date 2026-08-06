# 🛡️ DROS™ (Deterministic Runtime Operating System)
**The Execution Governance Standard for Agentic AI**

[![Marketplace Status](https://img.shields.io/badge/GitHub_Marketplace-Available-success)](https://github.com/marketplace/dros-engine)
[![License](https://img.shields.io/badge/License-Commercial-blue.svg)](#)
[![Official Website](https://img.shields.io/badge/Website-dr--os.io-purple.svg)](https://dr-os.io)

## 🛑 "If runtime needs intelligence, the system is already broken."
Prompt Engineering is dead when it comes to enterprise security. No matter how complex your System Prompt is, Jailbreaks and Prompt Injections will find a way through. API Gateways fail because they assume human-driven execution constraints apply to agents.

**DROS** is NOT a prompt wrapper. It is the **Execution Governance Standard**. We move intelligence to compile-time (Vajra.md) and enforce rules via a **Dual-Layer Security Architecture**:
*   **L1 ATR Semantic Radar**: Intercepts prompt injections and indirect RAG/web context contamination (T001/T002) at the ingestion gateway.
*   **L2 Vajra Enforcer**: An **O(1) Bitmap Runtime Engine** (`.dll` / `.so` / `.aar` / C-ABI dynamic library) enforcing strict syscall boundaries.

If the Agent attempts an unauthorized action, DROS triggers a **Physical Fusing (Strict Fail-Closed)** in under 1 microsecond. 

Verified under the **Agent Security Benchmark (ASB v1.1.0)**:
*   **Zero-Trust Security**: 100.0% enforcement rate against Easy, Medium, Hard, and Adaptive attacks at L2 boundary (Difficulty Insensitivity).
*   **Production Efficacy**: F1-Score of **0.973** and False Positive Rate (FPR) of **1.8%** (tested with $n=500$ benign workflows).
*   **Extreme Performance**: 484.8 ns pure bitmap validation latency.

---

## 🔥 Enterprise Product Lines (V3 Ecosystem)

### 1. DROS Free-Trial (PLG Edition)
- **Zero-Friction Adoption**: Download the binary and experience O(1) interception in 3 minutes.
- **Offline Timebomb**: 30-day RSA signed license. No need to connect to our servers to validate.
- **Limit**: Binds 1 Machine UUID, Max 2 Concurrent Agents (Hardcoded at compile time).

### 2. VajraClaw Hacker Tier
- **Developer/Indie Focused**: Perfect for individual developers, researchers, and local single-machine experiments.
- **Full Local SDK**: Access to core Go/Zig-based dynamic libraries (`.dll` / `.so` / `.dylib`) and C-FFI bindings. Includes `dros-cli` for automated AST contract compilation.
- **Limit**: Binds 1 Machine UUID, Max 5 Concurrent Agents. No control plane mesh.

### 3. DROS Engine (Core Runtime)
- **Compile-Time Intelligence**: Write policies in `Vajra.md`, compile them into deterministic Bitmaps.
- **O(1) Execution**: Pure bitwise AND operations. Zero JSON parsing at runtime. Zero TOCTOU vulnerabilities.
- **Multi-Language SDKs**: Native Python `vajraclaw` module included. Wrap LangChain/AutoGen tools in 1 line of code.

### 4. DROS Engine+ (Enterprise Audit)
- **Cryptographic Audit**: `SHA-256 Policy Hash` and `Compiler Version` are bound to every audit trace. Prove to your CISO exactly *why* an agent made a decision.
- **Sealed Binary Policy**: Ed25519 signatures prevent tampering.
- **Policy Epoch Lock**: Prevent downgrade attacks in distributed microservices.
- **Scale Limit**: Max 15 installations, each supporting 30 concurrent agents (Total 450 concurrent agents). Scale beyond this requires Corporate / VajraAgent Mesh custom quote.

### 5. VajraAgent (Mesh / Corporate)
- **Centralized Visual CommandCenter**: Point-and-click IAM monitoring, real-time 0.37ms interception visualization, and global system health tracking.
- **fully On-Premise & Air-Gapped Support**: Deployable in fully disconnected, high-security private datacenters (FinTech/Healthcare/Government).
- **K8s Helm & Ansible Automated Deployment**: Native Helm charts and Ansible playbooks for zero-touch provisioning and auto-scaling across thousands of microservices.
- **Mesh OTA & 1ms Global Kill-Switch**: Push cryptographically signed `policy.bin` updates to all nodes globally in real-time, featuring a 1ms instantaneous emergency block broadcast.
- **MCP Reverse Proxy**: Zero-code integration. Point any Model Context Protocol (MCP) compatible agent to VajraAgent for instant execution governance.

### 6. DROS Mobile SDK (Zero-Trust Endpoint)
- **Gomobile Android AAR**: O(1) engine compiled natively for mobile edge devices.
- **Offline Safe Degraded (Mode C)**: If the phone is offline, automatically downgrade to a Read-Only sandbox. Block all write actions (e.g., transfers) while keeping read queries functional.

---

## 💻 Zero-Pollution Integration (Python Example)
```python
from vajraclaw import VajraClaw

# 1. Mount the Signed Binary Policy
vc = VajraClaw("policy.bin")

# 2. Evaluate the Tool Call in O(1) Time
decision = vc.evaluate("customer_service", "admin.delete_user")

# 3. Deterministic Enforcement
if decision == "BLOCK":
    raise PermissionError("Action blocked by DROS Execution Kernel.")
```

---

## 🔒 Zero-Trust Privilege Separation
DROS enforces strict privilege separation at the OS level. The AI Agent must **never** be granted write permissions to `policy.bin`. In enterprise Kubernetes/Docker deployments, the policy is mounted as a **Read-Only ConfigMap**. Even if the LLM becomes entirely hostile, OS-level file permissions will physically deny tampering.

---

## 💰 Pricing & Licensing (Commercial)

| Tier | Price | Best For | Included |
| :--- | :--- | :--- | :--- |
| **Free Trial** | $0 | Solo devs, PoC | Binds 1 Machine UUID, Max 2 Concurrent Agents |
| **VajraClaw Hacker** | $149 / yr | Personal, Single-Machine | Binds 1 Machine UUID, Max 5 Concurrent Agents, Local SDK |
| **Startup License** | $799 / yr | Small teams & SaaS MVP | Binds 3 Machine UUIDs, Max 10 Concurrent Agents per machine, SDK |
| **Enterprise Audit** | $7,990 / yr | High-Security Corporate (FinTech) | Binds 15 Machine UUIDs, Max 30 Concurrent Agents per machine, Ed25519, SHA-256 Audit |
| **VajraAgent Mesh**| Custom | Enterprise Deployments | Control Plane, MCP Proxy, Global OTA Updates |
| **Source Code Buyout**| Custom | Defense, Medical | Full Source Code (Go/Zig), Audit Reports, White-labeling rights |

---
**Developed by DROS Labs / 康宸園有限公司 / Jimmy Chen**
*Securing the autonomous frontier through Deterministic Execution.*
