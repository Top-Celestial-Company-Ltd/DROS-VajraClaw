package mobile

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestInitStaticVajraFromString(t *testing.T) {
	contract := `
# DROS Security Contract
[BLOCK]
DROP TABLE
DELETE FROM
EXECUTE IMMEDIATE
`
	res := InitStaticVajraFromString(contract)
	if res != 1 {
		t.Fatalf("Expected 1 from InitStaticVajraFromString, got %d", res)
	}

	if len(staticVajraRules) != 7 { // 4 custom + 3 default ("診斷", "投資", "理財")
		t.Fatalf("Expected 7 static rules, got %d", len(staticVajraRules))
	}
}

func TestMatchTokenStreamWithAudit_Pass(t *testing.T) {
	InitStaticVajraFromString("RESTRICTED_WORD")
	ClearEphemeralRules()

	input := "This is a completely benign sentence."
	agentId := "agent-fin-01"

	jsonStr := MatchTokenStreamWithAudit(input, agentId)

	var trace AuditTrace
	err := json.Unmarshal([]byte(jsonStr), &trace)
	if err != nil {
		t.Fatalf("Failed to unmarshal audit JSON: %v, raw: %s", err, jsonStr)
	}

	if trace.Status != "PASS" {
		t.Errorf("Expected PASS, got %s", trace.Status)
	}
	if trace.AgentID != agentId {
		t.Errorf("Expected agentID %s, got %s", agentId, trace.AgentID)
	}
	if trace.MatchedRule != "" {
		t.Errorf("Expected empty matched_rule on PASS, got %s", trace.MatchedRule)
	}
	if trace.PayloadSnippet != input {
		t.Errorf("Expected payload snippet to match input")
	}
	if trace.TNumber == "" {
		t.Errorf("TNumber is missing")
	}
	if trace.LatencyMs < 0 {
		t.Errorf("Latency is negative")
	}
}

func TestMatchTokenStreamWithAudit_Block(t *testing.T) {
	InitStaticVajraFromString("RESTRICTED_WORD")
	ClearEphemeralRules()

	input := "I will output RESTRICTED_WORD now."
	agentId := "agent-hr-02"

	jsonStr := MatchTokenStreamWithAudit(input, agentId)

	var trace AuditTrace
	err := json.Unmarshal([]byte(jsonStr), &trace)
	if err != nil {
		t.Fatalf("Failed to unmarshal audit JSON: %v, raw: %s", err, jsonStr)
	}

	if trace.Status != "BLOCK" {
		t.Errorf("Expected BLOCK, got %s", trace.Status)
	}
	if trace.MatchedRule != "RESTRICTED_WORD" {
		t.Errorf("Expected matched_rule to be RESTRICTED_WORD, got %s", trace.MatchedRule)
	}
}

func TestMatchTokenStreamWithAudit_Truncation(t *testing.T) {
	InitStaticVajraFromString("NO_MATCH")
	ClearEphemeralRules()

	// Generate a huge input (1000 chars)
	input := strings.Repeat("A", 1000)
	
	jsonStr := MatchTokenStreamWithAudit(input, "agent-test")
	var trace AuditTrace
	json.Unmarshal([]byte(jsonStr), &trace)

	if len(trace.PayloadSnippet) >= 1000 {
		t.Errorf("Payload should be truncated, but got length %d", len(trace.PayloadSnippet))
	}
	if !strings.HasSuffix(trace.PayloadSnippet, "[TRUNCATED]") {
		t.Errorf("Payload should end with [TRUNCATED]")
	}
}

func TestEvaluateToolCallWithAudit_Payment_Pass(t *testing.T) {
	agentId := "agent-fin-01"
	toolName := "execute_payment"
	argsJson := `{"amount": 4999.0, "currency": "USD"}`

	jsonStr := EvaluateToolCallWithAudit(toolName, argsJson, agentId)
	var trace AuditTrace
	json.Unmarshal([]byte(jsonStr), &trace)

	if trace.Status != "PASS" {
		t.Errorf("Expected PASS, got %s. MatchedRule: %s", trace.Status, trace.MatchedRule)
	}
}

func TestEvaluateToolCallWithAudit_Payment_Block(t *testing.T) {
	agentId := "agent-fin-01"
	toolName := "execute_payment"
	argsJson := `{"amount": 5001.0, "currency": "USD"}`

	jsonStr := EvaluateToolCallWithAudit(toolName, argsJson, agentId)
	var trace AuditTrace
	json.Unmarshal([]byte(jsonStr), &trace)

	if trace.Status != "BLOCK" {
		t.Errorf("Expected BLOCK, got %s", trace.Status)
	}
	if trace.MatchedRule != "CAPABILITY_VIOLATION: amount > 5000 limit" {
		t.Errorf("Unexpected matched rule: %s", trace.MatchedRule)
	}
}

func TestEvaluateToolCallWithAudit_MalformedJSON(t *testing.T) {
	agentId := "agent-fin-01"
	toolName := "execute_payment"
	argsJson := `{ amount: 10000, broken: json }`

	jsonStr := EvaluateToolCallWithAudit(toolName, argsJson, agentId)
	var trace AuditTrace
	json.Unmarshal([]byte(jsonStr), &trace)

	if trace.Status != "BLOCK" {
		t.Errorf("Expected BLOCK for malformed JSON, got %s", trace.Status)
	}
	if trace.MatchedRule != "MALFORMED_JSON_PAYLOAD" {
		t.Errorf("Unexpected matched rule: %s", trace.MatchedRule)
	}
}

func TestEvaluateToolCallWithAudit_ForbiddenTool(t *testing.T) {
	agentId := "agent-db-01"
	toolName := "drop_database"
	argsJson := `{"db_name": "production"}`

	jsonStr := EvaluateToolCallWithAudit(toolName, argsJson, agentId)
	var trace AuditTrace
	json.Unmarshal([]byte(jsonStr), &trace)

	if trace.Status != "BLOCK" {
		t.Errorf("Expected BLOCK for forbidden tool, got %s", trace.Status)
	}
}
