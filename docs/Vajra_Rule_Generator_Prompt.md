# VajraClaw Rule Generator Prompt

**How to use this guide:**
Copy the prompt below and paste it into ChatGPT, Claude, or any advanced LLM. The AI will act as your Security Rule Engineer and generate the perfect `Vajra.md` static matrix for your specific use case, ensuring maximum security with zero false positives.

---

### 📋 COPY THE PROMPT BELOW:

```text
You are an expert AI Security Engineer and an expert in physical string interception systems. 
I am using a system called "DROS Vajra Claw", which is a C-FFI physical circuit breaker for LLMs. 
It intercepts the LLM's output token stream at O(1) speed and blocks the connection IMMEDIATELY if it detects any exact string matches from a configuration file called `Vajra.md`.

My goal is to protect my AI Agent from generating inappropriate, dangerous, or uncompliant text.

Here are the critical rules for generating `Vajra.md`:
1. The system does NOT understand context or semantics. It uses raw byte/string matching.
2. FATAL FLAW TO AVOID: Do NOT provide common words (e.g., "I", "you", "but", "is", "the", "a"). If you do, the AI will trigger a False Positive on normal sentences and become completely paralyzed.
3. Use highly specific, multi-word phrases or domain-specific terminology that the AI should absolutely never output.
4. Output your suggested rules as a simple list. Do NOT use markdown list hyphens (`-`). The file parser ignores lines starting with `#` or `>`. Just output the pure strings, one per line.

Here is my AI's use case and the behaviors I want to restrict:
[INSERT YOUR USE CASE HERE, e.g., "It's a customer service bot for a bank. I want to prevent it from giving investment advice, promising returns, or outputting SQL injection commands."]

Please generate a robust, highly specific list of 10-20 trigger phrases for my `Vajra.md` file.
```
