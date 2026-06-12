# Vajra DSL Specification v1

## Overview

Vajra DSL is a declarative policy language designed for deterministic AI agent governance.

It defines **what an agent is allowed to do**, not how it executes.

---

## Versioning

```yaml
vajra_version: 1
```

* Required field
* Missing field → FATAL compile error
* Ensures backward compatibility within v1.x

---

## Core Concepts

### 1. Agents

Represents execution identity:

```yaml
agents:
  - id: customer_service
    role: support
```

---

### 2. Capabilities

Logical permissions grouped by intent:

```yaml
capabilities:
  customer_service:
    - READ_CRM
```

---

### 3. Deny Priority

Deny always overrides allow:

```text
DENY > ALLOW (absolute rule)
```

---

### 4. Wildcards

Supports prefix-based expansion in explicitly declared rules:

```yaml
rules:
  - match:
      agent: customer_service
      tool: crm.read.*
    effect: ALLOW
```

Expands at compile-time into matching tools.

---

### 5. Requires (Capability-Based Security)

```yaml
tools:
  - name: crm.read.profile
    requires:
      - READ_CRM
```

Capabilities are resolved during compilation into bitmap mappings. If an agent possesses the required capabilities, the compiler implicitly maps an ALLOW for that tool, unless explicitly overridden by a DENY rule.

---

## Evaluation Model

DROS does NOT evaluate policies at runtime.

Instead:

```text
Vajra.md → Compiler → Bitmap Artifact → O(1) lookup
```

---

## Deterministic Compilation Rules

* Alphabetical sorting of all inputs (Agents, Tools, Capabilities)
* No runtime parsing of YAML
* Stable hash output required
* Same input + same static key → identical binary output (SHA-256)

---

## Enforcement Semantics

* Missing policy → DENY
* Unknown agent → DENY
* Unknown tool → DENY
* Invalid signature → BLOCK
* Tampered payload → BLOCK

---

## Security Principle

> Deny-by-default, enforce-at-runtime, verify-at-compile-time
