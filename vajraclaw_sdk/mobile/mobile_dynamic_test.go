package mobile

import (
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
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

	jsonStr := EvaluateDynamicToolCallWithAudit(toolName, argsJson, agentId, "")
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

	jsonStr := EvaluateDynamicToolCallWithAudit(toolName, argsJson, agentId, "")
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

	jsonStr := EvaluateDynamicToolCallWithAudit("drop_database", `{"db": "users"}`, "agent-db", "")
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
	jsonStr := EvaluateDynamicToolCallWithAudit("execute_payment", `{"currency": "USD"}`, "agent-fin-01", "")
	var trace AuditTrace
	json.Unmarshal([]byte(jsonStr), &trace)

	if trace.Status != "BLOCK" {
		t.Errorf("Expected BLOCK for missing field, got %s", trace.Status)
	}
	if trace.MatchedRule != "SCHEMA_VIOLATION: Missing required field: amount" {
		t.Errorf("Unexpected matched rule: %s", trace.MatchedRule)
	}
}

func TestEvaluateDynamicToolCallWithAudit_EpochDrift(t *testing.T) {
	policyJson := `{ "epoch": "v2", "tool_policies": [ { "tool_name": "execute_payment", "action": "ALLOW", "conditions": [] } ] }`
	InitDynamicPolicyFromJson(policyJson)

	// Matching epoch should PASS
	jsonStr := EvaluateDynamicToolCallWithAudit("execute_payment", `{}`, "agent-1", "v2")
	var trace AuditTrace
	json.Unmarshal([]byte(jsonStr), &trace)
	if trace.Status != "PASS" {
		t.Errorf("Expected PASS for matching epoch, got %s", trace.Status)
	}

	// Mismatched epoch should BLOCK
	jsonStrDrift := EvaluateDynamicToolCallWithAudit("execute_payment", `{}`, "agent-1", "v1-outdated")
	var traceDrift AuditTrace
	json.Unmarshal([]byte(jsonStrDrift), &traceDrift)
	if traceDrift.Status != "BLOCK" || traceDrift.MatchedRule != "EPOCH_MISMATCH" {
		t.Errorf("Expected BLOCK with EPOCH_MISMATCH, got %s / %s", traceDrift.Status, traceDrift.MatchedRule)
	}
}

func TestInitDynamicPolicyFromBinary(t *testing.T) {
	SetOperationalMode(ModeDegraded)
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	astJson := `{ "epoch": "v3-signed", "tool_policies": [] }`
	payload := []byte(astJson)
	sig := ed25519.Sign(priv, payload)

	binData := make([]byte, 0)
	binData = append(binData, []byte("VAJRAC")...)
	binData = append(binData, 0x01)
	
	epochBytes := make([]byte, 32)
	copy(epochBytes, "v3-signed")
	binData = append(binData, epochBytes...)
	
	binData = append(binData, sig...)
	
	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(payload)))
	binData = append(binData, lenBytes...)
	
	binData = append(binData, payload...)

	pubHex := hex.EncodeToString(pub)

	res := InitDynamicPolicyFromBinary(binData, pubHex)
	if res != 1 {
		t.Fatalf("Failed to initialize from valid binary")
	}

	// Verify epoch was loaded
	jsonStr := EvaluateDynamicToolCallWithAudit("some_tool", `{}`, "agent-1", "v3-signed")
	var trace AuditTrace
	json.Unmarshal([]byte(jsonStr), &trace)
	if trace.MatchedRule == "EPOCH_MISMATCH" {
		t.Fatalf("Epoch lock failed, got EPOCH_MISMATCH for correct epoch v3-signed")
	}
	
	// Test Tampering
	binDataTampered := make([]byte, len(binData))
	copy(binDataTampered, binData)
	binDataTampered[len(binDataTampered)-1] = 'X' // Alter payload
	// Since payload is altered but signature is same, it should fail
	
	resTampered := InitDynamicPolicyFromBinary(binDataTampered, pubHex)
	if resTampered != 0 {
		t.Fatalf("Tampered binary should fail signature check")
	}
}

