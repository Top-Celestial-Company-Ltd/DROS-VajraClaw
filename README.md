# 🛡️ VajraClaw
**Deterministic execution enforcement layer for agentic AI systems.**

[![License](https://img.shields.io/badge/License-Commercial-blue.svg)](#-pricing--licensing)
[![Official Website](https://img.shields.io/badge/Website-dr--os.io-purple.svg)](https://dr-os.io)
[![Gumroad](https://img.shields.io/badge/Buy-Gumroad-orange.svg)](https://drosvajra.gumroad.com)

VajraClaw ensures AI agents can **only execute actions within explicitly defined capabilities** — enforced at runtime with O(1) deterministic checks, no LLM involved.

---

## ⚠️ Why This Exists

Modern AI agents are powerful — but unsafe by default:

- **Prompt injection**: Malicious instructions hijack agent behavior
- **Tool misuse**: Agents call APIs or write files they should never touch
- **Silent privilege escalation**: No log, no trace, no way to know it happened
- **Semantic guardrails fail**: LLM-as-a-judge inherits the same flaws as the model it guards

**The industry's answer — "use another LLM to check the first LLM" — is broken by design.**

---

## 🧠 Core Concept

Instead of filtering prompts or reasoning about intent, VajraClaw enforces a single rule:

```
Agent → Capability Check → Tool Execution
```

If the agent is not **explicitly authorized** for the action:

```
❌ Execution is blocked (fail-closed)
```

No inference. No heuristics. No bypass.

---

## ⚙️ Architecture

```
[ AI Agent Output ]
        ↓
[ Tool Invocation ]
        ↓
[ VajraClaw Enforcement Layer ]    ← O(1) capability check
        ↓
  ┌─────────────┐
  │ ALLOW       │ → Tool executes
  │ BLOCK       │ → Audit log + hard stop
  └─────────────┘
```

**VajraClaw is a C-shared binary microkernel (`.dll` / `.so`)** — not a Python wrapper, not a prompt template. It sits at the execution boundary and enforces capability policy before any system call is made.

---

## 🔑 Key Properties

| Property | Detail |
|:---|:---|
| **Enforcement model** | Capability-based (explicit allowlist, deny everything else) |
| **Decision speed** | O(1) — memory-mapped bitmap lookup, sub-millisecond |
| **LLM dependency** | Zero — security layer never calls any model |
| **Failure mode** | Fail-closed — missing policy = block, not pass |
| **Bypass surface** | None by design — interceptor sits at network/SDK layer |
| **Audit trail** | Append-only JSONL, every decision recorded |

---

## 🚀 Quick Start

### 1. Clone & Mount the Core Binary

```bash
git clone https://github.com/Top-Celestial-Company-Ltd/VajraClaw
cd VajraClaw
```

Drop `core/vajra_claw.dll` (Windows) or `core/vajra_claw.so` (Linux/macOS) into your project.

### 2. Define a Capability Policy (`Vajra.md`)

```yaml
agents:
  - id: finance-agent
    capabilities:
      - READ_DB        # allowed
      # WRITE_DB      # not listed = denied by default
```

### 3. Wrap Your Agent's Tool Call

```python
from adapters.claw_adapter import VajraClawAdapter

claw = VajraClawAdapter(dll_path="./core/vajra_claw.dll", rule_path="Vajra.md")

# Before executing any tool:
if not claw.evaluate(agent_id="finance-agent", tool="db.write"):
    raise PermissionError("VajraClaw: execution blocked — WRITE_DB not authorized")
```

### 4. What You See When It Blocks

```
❌ DENIED
Agent:       finance-agent
Tool:        db.write
Reason:      Capability WRITE_DB not in policy
Mode:        STRICT (Mode D)
Policy Hash: a3f9c12e...
Timestamp:   2026-05-30T07:14:22Z
```

---

## 🧪 Demo: Simulate an Attack

> **LangChain / MCP integration demos**: See [`integrations/`](./integrations/) — coming in next release.

The simplest way to trigger a DENY today — run the Go SDK test directly:

```bash
cd vajraclaw_sdk/mobile
go test -run TestModeCDegradedWrite -v
```

Expected output:
```
--- PASS: TestModeCDegradedWrite (0.00s)
    [VajraClaw] BLOCK | reason: DEGRADED_MODE_WRITE_BLOCKED
```

---

## 🔐 Security Model

- **No runtime LLM decision-making** — policy is compiled, signed, and locked
- **No prompt-based filtering** — string matching is explicitly not the primary mechanism
- **Cryptographic policy binding** — Ed25519 signature verification; tampered policy → fatal halt
- **Operational Modes**:
  - **Mode C (Safe Degraded)**: Network isolated? Write ops blocked, read ops continue
  - **Mode D (Strict Fail-Closed)**: No valid policy = no execution, period

---

## ⚠️ What VajraClaw Does NOT Protect Against

We are 100% transparent about the physical limits of our O(1) architecture:

1. **Semantic synonyms**: Blacklisting `"delete database"` does not catch `"drop the tables"` — VajraClaw enforces capability boundaries, not semantic intent
2. **Obfuscated payloads**: Base64 or encoded commands are not decoded and scanned unless explicitly defined
3. **Host environment vulnerabilities**: We enforce at the execution boundary — not a substitute for OS-level sandboxing

**The golden rule**: VajraClaw is the inner execution firewall. Keep Docker/K8s as your outer armor.

---

## 💰 Pricing & Licensing

| Tier | Price | Agents | Best For |
|:---|:---|:---|:---|
| **Hacker Edition** | Free | 2 concurrent | Developers, open source, non-commercial |
| **Startup License** | $499 / yr | 10 / machine | Small teams & SaaS MVPs |
| **Enterprise License** | $4,990 / yr | 30 / machine | High-security corporate & regulated industries |
| **Source Code Buyout** | Custom | Unlimited | Defense, Medical, Banking |

👉 [**Buy on Gumroad**](https://drosvajra.gumroad.com) · [**Official Portal**](https://dr-os.io)

---

## 📁 Repository Structure

```
VajraClaw/
├── core/                    # C-shared binary microkernel (.dll / .so)
├── vajraclaw_sdk/
│   └── mobile/              # Go SDK — full enterprise governance layer
│       ├── mobile.go            # Static rule enforcement
│       ├── mobile_dynamic.go    # Dynamic policy AST engine + Modes C/D
│       ├── mobile_audit.go      # Append-only JSONL audit log
│       └── mobile_interceptor.go # HTTP execution path interceptor
├── adapters/                # Language bindings (Python ctypes, Node.js)
├── integrations/            # LangChain, MCP — coming soon
├── docs/
│   ├── ADR-001-Architecture-Boundary.md
│   └── ...
└── rules/                   # Example Vajra.md policy files
```

---

**Developed by Top-Celestial Company Ltd. (康宸園有限公司) / Jimmy Chen**  
*"We control what AI is allowed to execute."*
