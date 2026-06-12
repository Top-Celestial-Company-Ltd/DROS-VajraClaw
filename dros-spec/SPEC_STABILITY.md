# DROS Specification Stability Promise

## Overview

The Vajra DSL (Domain Specific Language) serves as the open standard for Agentic AI Execution Governance.

To ensure enterprises can confidently build their security posture on top of DROS, we provide strict backward compatibility guarantees for the DSL.

## The v1.x Promise

> **All v1 DSL definitions guarantee backward compatibility within the v1.x series.**

When you define a policy using:

```yaml
vajra_version: 1
```

We guarantee that:

1. **Compiler Compatibility**: Any future `v1.x` version of the Vajra Compiler will successfully parse and compile your existing policy without requiring syntax modifications.
2. **Deterministic Enforcement**: The semantic meaning of your rules (e.g., capability propagation, rule precedence) will remain identical across compiler upgrades.
3. **Audit Log Consistency**: The cryptographic payload verification process ensures that a policy compiled today will enforce the same security boundaries tomorrow.

## Breaking Changes

Any breaking changes to the syntax, semantics, or capability resolution model will require a major version bump (e.g., `vajra_version: 2`).

If and when `vajra_version: 2` is released, the DROS compiler will continue to support `vajra_version: 1` indefinitely, ensuring your existing security infrastructure remains intact while providing a graceful migration path.
