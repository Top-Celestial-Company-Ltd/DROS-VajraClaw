# DROS Golden Path Examples

This document demonstrates how to model real-world Agentic AI governance scenarios using the Vajra DSL.

## Scenario 1: The Principle of Least Privilege

You have a Data Analytics Agent that needs to query a database and write reports to a local file system, but under no circumstances should it be able to delete data or execute system commands.

### The Policy (`analytics_policy.yaml`)
```yaml
vajra_version: 1
epoch: 2026-Q3-PROD

agents:
  - id: data_analyst_agent
    role: analytics

tools:
  - name: db.query
    requires: [READ_DB]
  - name: db.delete
    requires: [ADMIN_DB]
  - name: file.write
    requires: [WRITE_REPORT]
  - name: sys.exec
    requires: [ROOT]

capabilities:
  data_analyst_agent:
    - READ_DB
    - WRITE_REPORT

# Implicit ALLOWs are automatically generated for db.query and file.write.
# All other tools remain implicitly DENIED.
```

## Scenario 2: Fail-Safe Override (The Kill Switch)

A Customer Service Agent has been granted wildcard access to read CRM data (`crm.read.*`). However, a new API `crm.read.billing_secrets` is introduced, and you want to explicitly block it without altering the capability structure.

### The Policy (`support_policy.yaml`)
```yaml
vajra_version: 1

agents:
  - id: support_agent

tools:
  - name: crm.read.profile
  - name: crm.read.history
  - name: crm.read.billing_secrets

rules:
  # 1. Grant broad access
  - match:
      agent: support_agent
      tool: crm.read.*
    effect: ALLOW

  # 2. Explicit Deny Override (DENY always wins)
  - match:
      agent: support_agent
      tool: crm.read.billing_secrets
    effect: DENY
```

## Scenario 3: CI/CD Pipeline Enforcement

DROS is built to integrate into your deployment pipelines. 

By running `vajra lint` on your pull requests, you can prevent dangerous configurations from ever reaching production.

```bash
# In your GitHub Actions or GitLab CI:

$ vajra lint my_policy.yaml
[Vajra Linter] Found 1 potential issues:
  [CRITICAL] Dangerous Grant: Non-admin agent 'support_agent' granted access to 'sys.reboot'.

$ echo $?
1 # Fails the build automatically!
```
