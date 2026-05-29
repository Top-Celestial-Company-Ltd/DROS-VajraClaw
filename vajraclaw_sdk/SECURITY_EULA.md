# ⚠️ DROS™ Vajra Claw™ - Security Boundary & EULA

> **CRITICAL WARNING FOR ALL DEPLOYMENTS**
> Please read this document before mounting `vajra_claw.dll` / `vajra_claw.so` into your production environment.

## 1. The "Fail-Closed" vs "Fail-Open" Concurrency Threshold
Your VajraClaw license is hard-coded with a **Max Concurrent Agents** limit (e.g., 5, 10, or 30 agents per machine) based on your pricing tier.

**What happens if you exceed the limit?**
If you attempt to spawn an Agent instance that exceeds your licensed concurrency limit, the VajraClaw C-FFI kernel will trigger a **Fail-Closed** exception. The initialization of the SDK adapter will be blocked. 
However, if your developers bypass the SDK error handling and route the LLM output directly to the system bypassing VajraClaw, **your system will be completely exposed to Prompt Injection and Data Exfiltration**. 

*Legal Disclaimer*: DROS™ (Top-Celestial Company Ltd.) is not liable for any data breaches, financial losses, or system destruction caused by LLM hallucinations or Prompt Injections on un-monitored agent threads that exceeded your licensing quota.

## 2. Heartbeat Severance Protocol
For Hacker and Startup Tiers, VajraClaw requires an active connection to the DROS Heartbeat Server (`https://api.dr-os.io`) to validate your RSA Token.
If your machine loses internet connection for more than 24 hours, the Ephemeral Token will expire. **The physical circuit breaker will enter absolute lockdown.** All LLM outputs routed through VajraClaw will return empty strings until the connection is restored. 

*If your infrastructure requires offline operation, you MUST upgrade to the Air-Gapped Enterprise Tier.*

## 3. The Static Matrix (`Vajra.md`) Configuration Risk
VajraClaw strictly enforces the boundaries you write in the Static Matrix file. 
- If you misconfigure the Regex patterns in `Vajra.md` to be too broad (e.g., allowing `*.*.*.*` IP connections), VajraClaw will permit it. 
- **VajraClaw is a physical fuse, not an AI.** It does not guess your intent. It executes your physical boundaries with O(1) mathematical precision. You are solely responsible for the precision of your Matrix logic.

## 4. Reverse Engineering & Tempering
Any attempt to decompile, reverse-engineer, or tamper with the binary structure of `vajra_claw.dll` or `.so` to bypass licensing constraints will trigger the self-corrupting defense mechanism, permanently burning your License Key on the central server. 

---
**By mounting this binary, you accept the physical laws of the Vajra Matrix.**
*Support: service@dr-os.io*
