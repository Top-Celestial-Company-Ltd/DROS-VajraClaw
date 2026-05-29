package mobile

import (
	"encoding/json"
	"testing"
)

func TestInitDynamicPolicyFromJson(t *testing.T) {
	policyJson := `
	{
		"tool_policies": [
			{
				"tool_name": "drop_database",
				"action": "BLOCK",
				"conditions": []
			},
			{
				"tool_name": "execute_payment",
				"action": "BLOCK",
				"conditions": [
					{"field": "amount", "operator": ">", "value": 5000}
				]
			}
		]
	}
	`
	res := InitDynamicPolicyFromJson(policyJson)
	if res != 1 {
		t.Fatalf("InitDynamicPolicyFromJson failed")
	}

	if len(dynamicPolicyMap["drop_database"]) != 1 {
		t.Errorf("Expected 1 policy for drop_database")
	}
	if len(dynamicPolicyMap["execute_payment"]) != 1 {
		t.Errorf("Expected 1 policy for execute_payment")
	}
}

func TestEvaluateDynamicToolCallWithAudit_Pass(t *testing.T) {
	policyJson := `{ "tool_policies": [ { "tool_name": "execute_payment", "action": "BLOCK", "conditions": [ {"field": "amount", "operator": ">", "value": 5000} ] } ] }`
	InitDynamicPolicyFromJson(policyJson)

	agentId := "agent-fin-01"
	toolName := "execute_payment"
	argsJson := `{"amount": 4999.0, "currency": "USD"}`

	jsonStr := EvaluateDynamicToolCallWithAudit(toolName, argsJson, agentId)
	var trace AuditTrace
	json.Unmarshal([]byte(jsonStr), &trace)

	if trace.Status != "PASS" {
		t.Errorf("Expected PASS, got %s. MatchedRule: %s", trace.Status, trace.MatchedRule)
	}
}

func TestEvaluateDynamicToolCallWithAudit_Block(t *testing.T) {
	policyJson := `{ "tool_policies": [ { "tool_name": "execute_payment", "action": "BLOCK", "conditions": [ {"field": "amount", "operator": ">", "value": 5000} ] } ] }`
	InitDynamicPolicyFromJson(policyJson)

	agentId := "agent-fin-01"
	toolName := "execute_payment"
	argsJson := `{"amount": 5001.0, "currency": "USD"}`

	jsonStr := EvaluateDynamicToolCallWithAudit(toolName, argsJson, agentId)
	var trace AuditTrace
	json.Unmarshal([]byte(jsonStr), &trace)

	if trace.Status != "BLOCK" {
		t.Errorf("Expected BLOCK, got %s", trace.Status)
	}
	if trace.MatchedRule != "DYNAMIC_CAPABILITY_VIOLATION" {
		t.Errorf("Unexpected matched rule: %s", trace.MatchedRule)
	}
}

func TestEvaluateDynamicToolCallWithAudit_ForbiddenTool(t *testing.T) {
	policyJson := `{ "tool_policies": [ { "tool_name": "drop_database", "action": "BLOCK", "conditions": [] } ] }`
	InitDynamicPolicyFromJson(policyJson)

	jsonStr := EvaluateDynamicToolCallWithAudit("drop_database", `{"db": "users"}`, "agent-db")
	var trace AuditTrace
	json.Unmarshal([]byte(jsonStr), &trace)

	if trace.Status != "BLOCK" {
		t.Errorf("Expected BLOCK for forbidden tool, got %s", trace.Status)
	}
}

func TestEvaluateDynamicToolCallWithAudit_SchemaViolation(t *testing.T) {
	policyJson := `{ "tool_policies": [ { "tool_name": "execute_payment", "action": "BLOCK", "conditions": [ {"field": "amount", "operator": ">", "value": 5000} ] } ] }`
	InitDynamicPolicyFromJson(policyJson)

	// Missing amount
	jsonStr := EvaluateDynamicToolCallWithAudit("execute_payment", `{"currency": "USD"}`, "agent-fin-01")
	var trace AuditTrace
	json.Unmarshal([]byte(jsonStr), &trace)

	if trace.Status != "BLOCK" {
		t.Errorf("Expected BLOCK for missing field, got %s", trace.Status)
	}
	if trace.MatchedRule != "SCHEMA_VIOLATION: Missing required field: amount" {
		t.Errorf("Unexpected matched rule: %s", trace.MatchedRule)
	}
}
