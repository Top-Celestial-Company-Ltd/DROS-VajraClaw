# VajraClaw Architecture Boundary: ADR-001

**Decision Date**: 2026-05-30  
**Status**: Accepted

---

## Decision

VajraClaw maintains a strict two-layer architecture. The boundary between layers is **permanent and intentional**.

---

## Two Layers

### Layer 1: VajraClaw Core (`core/vajra_claw.go` → `.dll` / `.so`)

**Philosophy**: Ultra-thin physical enforcement primitive. Does one thing. Does it at O(1).

**Capabilities**:
- Token stream matching against a crystallized rule matrix
- Static rule loading (`init_static_vajra_from_string`)
- Ephemeral rule injection and evaporation (`inject_ephemeral_rule`)
- Binary ALLOW / BLOCK decision (`match_token_stream` → 1 or 0)

**Intentionally excluded**:
- Operational Modes (C/D)
- Audit Logging
- Cryptographic Signature Verification
- Policy Epoch Locking
- HTTP Interceptor

**Target consumers**: OEM partners, advanced integrators who need a bare enforcement primitive and will build their own governance wrapper on top.

---

### Layer 2: VajraClaw SDK (`vajraclaw_sdk/mobile/` and future `server/`)

**Philosophy**: Complete enterprise governance layer built on top of Core. This is what enterprises deploy.

**Capabilities** (all SRE v1 guarantees live here):
- `SetOperationalMode()` — Mode C (Safe Degraded) / Mode D (Strict Fail-Closed)
- `SyncPolicyFromMesh()` — Ed25519 cryptographic policy verification. Invalid signature → Fatal Panic.
- `EvaluateDynamicToolCallWithAudit()` — Full capability-based AST evaluation with audit output
- `ConfigureAuditLog()` — Append-only JSONL audit trail (all 9 fields per Kernel Spec v1.0)
- `VajraClawRoundTripper` — HTTP network interceptor (Execution Path Uniqueness)
- Policy Epoch binding (coming Sprint 1 Phase 2)

**Target consumers**: All standard enterprise customers, AI SaaS startups, FinTech teams, Mobile/Edge AI deployments.

---

## Deployment Boundary Rule

> **Enterprise SRE guarantees (Audit, Mode C/D, Ed25519, Epoch Lock) are ONLY available through the SDK layer. The Core is a bare engine.**

Any customer integrating Core directly assumes full responsibility for building their own governance wrapper. This must be stated clearly in the EULA and technical documentation.

---

## Why This Is The Right Call

| Concern | Core | SDK |
|:---|:---|:---|
| Binary size | ~3MB | ~5MB |
| Startup time | Microseconds | Microseconds |
| External dependencies | Zero | Zero (Go stdlib only) |
| Audit capability | None (by design) | Full JSONL trail |
| SRE Mode support | None (by design) | Mode C / Mode D |
| Target | OEM / Advanced | Enterprise standard |

Keeping Core surgically thin means it can be embedded in the most constrained environments (microcontrollers, RTOS, custom hardware) without dragging in any governance overhead. The SDK handles governance for everyone else.

**This is the same pattern as SQLite (core) vs. any application built on top of it.**
