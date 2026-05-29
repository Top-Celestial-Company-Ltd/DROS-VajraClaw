package mobile

import (
	"crypto/ed25519"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestModeCDegradedWrite(t *testing.T) {
	ClearDynamicPolicies()
	SetOperationalMode(ModeDegraded)

	// Tool with IsWriteAction = true
	jsonPolicy := `{
		"tool_policies": [
			{
				"tool_name": "transfer_funds",
				"action": "ALLOW",
				"is_write_action": true,
				"conditions": []
			}
		]
	}`
	InitDynamicPolicyFromJson(jsonPolicy)

	// Since Mode = Degraded and tool is Write, it must block regardless of ALLOW action
	auditJSON := EvaluateDynamicToolCallWithAudit("transfer_funds", `{"amount":100}`, "agent-1")
	if !strings.Contains(auditJSON, `"status":"BLOCK"`) {
		t.Fatalf("Expected BLOCK for write action in Degraded mode, got: %s", auditJSON)
	}
	if !strings.Contains(auditJSON, `"matched_rule":"DEGRADED_MODE_WRITE_BLOCKED"`) {
		t.Fatalf("Expected matched_rule to be DEGRADED_MODE_WRITE_BLOCKED")
	}
}

func TestModeDFailClosedNoPolicy(t *testing.T) {
	ClearDynamicPolicies()
	SetOperationalMode(ModeStrictFailClosed)

	// No policies loaded. Evaluating ANY tool should block.
	auditJSON := EvaluateDynamicToolCallWithAudit("unknown_tool", `{"foo":"bar"}`, "agent-1")
	if !strings.Contains(auditJSON, `"status":"BLOCK"`) {
		t.Fatalf("Expected BLOCK for unknown tool in Strict mode")
	}
	if !strings.Contains(auditJSON, `"matched_rule":"STRICT_MODE_NO_POLICY"`) {
		t.Fatalf("Expected matched_rule to be STRICT_MODE_NO_POLICY")
	}
}

func TestSignatureVerificationFailClosed(t *testing.T) {
	SetOperationalMode(ModeStrictFailClosed)

	// Using a defer to catch the expected panic
	defer func() {
		if r := recover(); r != nil {
			errStr := r.(string)
			if !strings.Contains(errStr, "FATAL: Cryptographic Signature Verification Failed") && 
			   !strings.Contains(errStr, "Invalid Signature format") &&
			   !strings.Contains(errStr, "Invalid Public Key format") {
				t.Fatalf("Unexpected panic message: %v", r)
			}
		} else {
			t.Fatalf("Expected panic due to strict signature failure, but code continued executing")
		}
	}()

	pub, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	pubHex := hex.EncodeToString(pub)
	invalidSigHex := hex.EncodeToString(make([]byte, ed25519.SignatureSize))

	// This should panic in Strict mode
	SyncPolicyFromMesh(`{"tool_policies":[]}`, invalidSigHex, pubHex)
}

// Mock RoundTripper for interceptor testing
type mockTransport struct{}
func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200}, nil
}

func TestInterceptorBlocks(t *testing.T) {
	ClearDynamicPolicies()
	SetOperationalMode(ModeStrictFailClosed)

	interceptor := NewVajraClawRoundTripper(&mockTransport{}, "agent-1")

	// Create a mock request that attempts a tool call
	req, _ := http.NewRequest("POST", "http://api.llm.com/v1/chat", strings.NewReader(`{"tool_name":"transfer_funds"}`))

	_, err := interceptor.RoundTrip(req)
	if err == nil {
		t.Fatalf("Expected interceptor to block and return error, but got nil")
	}
	if !strings.Contains(err.Error(), "Physical Melt") {
		t.Fatalf("Expected error to mention Physical Melt, got: %v", err)
	}
}

func TestInterceptorAllows(t *testing.T) {
	ClearDynamicPolicies()
	SetOperationalMode(ModeNormal) // Normal mode allows unknown tools unless explicitly defined to block in our current MVP logic

	// Let's actually define it to allow
	jsonPolicy := `{
		"tool_policies": [
			{
				"tool_name": "query_faq",
				"action": "ALLOW",
				"is_write_action": false,
				"conditions": []
			}
		]
	}`
	InitDynamicPolicyFromJson(jsonPolicy)

	interceptor := NewVajraClawRoundTripper(&mockTransport{}, "agent-1")
	req, _ := http.NewRequest("POST", "http://api.llm.com/v1/chat", strings.NewReader(`{"tool_name":"query_faq"}`))

	resp, err := interceptor.RoundTrip(req)
	if err != nil {
		t.Fatalf("Expected interceptor to allow, got error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected response status 200")
	}
	
	// Read the body to make sure it was restored properly
	bodyBytes, _ := io.ReadAll(req.Body)
	if !strings.Contains(string(bodyBytes), `"tool_name":"query_faq"`) {
		t.Fatalf("Body was not restored correctly by the interceptor")
	}
}
