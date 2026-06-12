# Vajra Linter & Doctor Governance Guide

The Vajra Compiler includes a powerful static analysis engine designed to shift security left. By analyzing your `Vajra.md` policies at design-time, we catch logical flaws before they reach the enforcement engine.

## 5-Level Severity Model

`vajra lint` analyzes the topology of your policy and emits warnings across 5 severity levels:

### 1. `[INFO]` (Informational)
Identifies inefficiencies or orphaned configurations that do not pose a direct security threat.
- **Example**: An agent is assigned a capability (e.g., `READ_LOGS`) that is never required by any tool in the ecosystem.
- **Action**: Clean up to maintain Least Privilege.

### 2. `[WARN]` (Warning)
Identifies potential misconfigurations or dead-code paths.
- **Example**: `Unreachable Tool` - A tool requires a capability that no agent currently possesses.
- **Example**: `Wildcard Mismatch` - A wildcard rule (e.g., `invalid.tool.*`) matches zero tools.
- **Example**: `Capability Duplicate Meaning` - Two distinct capabilities resolve to the exact same set of tools, indicating an IAM sprawl.

### 3. `[ERROR]` (Error)
Identifies explicit logical conflicts. **Fails the build.**
- **Example**: `Rule Conflict` - An implicit ALLOW from a capability is explicitly overridden by a DENY rule. (Note: DENY always wins, but the linter flags this to ensure the developer is aware of the override).
- **Example**: Multiple explicit rules contradict each other for the same Agent/Tool pair.

### 4. `[CRITICAL]` (Critical Security Risk)
Identifies syntactically valid but highly dangerous policy decisions. **Fails the build.**
- **Example**: `Dangerous Grant` - A non-administrative agent is granted access to high-risk wildcards (e.g., `admin.*` or `sys.*`).

### 5. `[FATAL]` (Fatal Parsing Failure)
Identifies structural failures that prevent compilation. **Fails the build.**
- **Example**: Missing `vajra_version`, malformed YAML syntax, or referencing an undefined agent in the capabilities block.

---

## Vajra Doctor (Health Assessment)

Run `vajra doctor policy.yaml` to receive a high-level health report tailored for CISOs and Security Architects.

### Policy Matrix Density
The Doctor calculates the "Sparse Matrix Density" (Total ALLOWs / (Total Agents * Total Tools)). 
- If the density is extremely low (e.g., < 1%), it indicates a highly fragmented capability structure (IAM Explosion).
- **Recommendation**: Consider refactoring and consolidating capabilities to simplify the governance model.

### Complexity Score
Based on unused capabilities, conflict risks, and matrix density, the policy is graded (A to F) to provide a quick snapshot of your Agentic AI's security posture.
