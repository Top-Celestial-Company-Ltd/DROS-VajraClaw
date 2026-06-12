# Vajra.md (Physical Execution Contract)

這份合約定義了 Agent 在作業系統級別的絕對物理邊界。
請注意：這不是給 LLM 看的「提示詞 (Prompt)」，這是直接編譯給 VajraClaw 引擎 (GuardVM) 執行的二進位授權規則。

```yaml
vajra_contract:
  version: "1.0"
  priority: "highest"
  rules:
    # [預設絕對紅線] 禁止任何底層 Shell 系統呼叫與檔案刪除
    - action: "deny"
      condition: "contains 'os.system' or 'rm -rf' or 'del /s'"
      response: "physical_melt"
      audit_log: true
      
    # [企業資料保護] 嚴禁讀取薪資或財務資料庫
    - action: "deny"
      condition: "api_call == 'query_database' and parameter contains 'salary' or 'finance_records'"
      response: "physical_melt"
      alert_admin: true

    # [範例放行] 允許 Agent 查詢網域狀態或重置密碼表單
    - action: "allow"
      condition: "api_call == 'check_network_status' or api_call == 'send_password_reset_link'"
      require_tnumber: false
```
