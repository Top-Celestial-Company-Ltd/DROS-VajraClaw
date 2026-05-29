package mobile

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	_ "golang.org/x/mobile/bind"
)

// Operational Modes
const (
	ModeNormal           = 0
	ModeIsolated         = 1 // Network loss, relies on local
	ModeDegraded         = 2 // Safe Degraded Mode (Read-only Sandbox)
	ModeStrictFailClosed = 3 // Default Strict Fail-Closed
)

var currentOperationalMode int = ModeStrictFailClosed

// SetOperationalMode allows the host/SRE to configure the failure mode
func SetOperationalMode(mode int) {
	dynamicMutex.Lock()
	defer dynamicMutex.Unlock()
	currentOperationalMode = mode
}

// Condition represents a single field validation rule.
type Condition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // ">", "<", "==", "!=", "contains"
	Value    interface{} `json:"value"`
}

// ToolPolicy defines rules for a specific tool call.
type ToolPolicy struct {
	ToolName      string      `json:"tool_name"`
	Action        string      `json:"action"` // "BLOCK" or "ALLOW"
	IsWriteAction bool        `json:"is_write_action"` // For Mode C (Degraded Sandbox)
	Conditions    []Condition `json:"conditions"`
}

// DynamicPolicy is the root schema for JSON input.
type DynamicPolicy struct {
	ToolPolicies []ToolPolicy `json:"tool_policies"`
}

// Global Memory for dynamic policies and audit queue
var (
	dynamicPolicyMap map[string][]ToolPolicy // Maps tool_name -> list of policies
	dynamicMutex     sync.RWMutex

	auditQueue []AuditTrace
	queueMutex sync.Mutex
)

/**
 * InitDynamicPolicyFromJson loads and compiles a JSON AST into memory.
 * Returns 1 on success, 0 on JSON parse error.
 */
func InitDynamicPolicyFromJson(jsonStr string) int {
	dynamicMutex.Lock()
	defer dynamicMutex.Unlock()

	var policy DynamicPolicy
	err := json.Unmarshal([]byte(jsonStr), &policy)
	if err != nil {
		return 0
	}

	// Rebuild the map for O(1) tool lookups
	dynamicPolicyMap = make(map[string][]ToolPolicy)
	for _, tp := range policy.ToolPolicies {
		dynamicPolicyMap[tp.ToolName] = append(dynamicPolicyMap[tp.ToolName], tp)
	}

	return 1
}

/**
 * ClearDynamicPolicies completely purges the dynamic AST engine memory.
 */
func ClearDynamicPolicies() {
	dynamicMutex.Lock()
	defer dynamicMutex.Unlock()
	dynamicPolicyMap = nil
}

/**
 * SyncPolicyFromMesh is the Tier 3 (Corporate) API.
 * It receives a JSON payload and verifies it against the provided public key (or secret)
 * before compiling the AST. This ensures Zero-Trust OTA updates.
 * (Mocking actual RSA/Ed25519 with a simple SHA256 HMAC-like check for MVP)
 */
func SyncPolicyFromMesh(jsonStr string, signatureHex string, pubKeyHex string) int {
	pubKey, err := hex.DecodeString(pubKeyHex)
	if err != nil || len(pubKey) != ed25519.PublicKeySize {
		if currentOperationalMode == ModeStrictFailClosed {
			panic("[VajraClaw] FATAL: Invalid Public Key format in Strict Mode")
		}
		return 0
	}

	sig, err := hex.DecodeString(signatureHex)
	if err != nil || len(sig) != ed25519.SignatureSize {
		if currentOperationalMode == ModeStrictFailClosed {
			panic("[VajraClaw] FATAL: Invalid Signature format in Strict Mode")
		}
		return 0
	}

	if !ed25519.Verify(ed25519.PublicKey(pubKey), []byte(jsonStr), sig) {
		if currentOperationalMode == ModeStrictFailClosed {
			panic("[VajraClaw] FATAL: Cryptographic Signature Verification Failed! Tampered Policy detected. Halting execution.")
		}
		return 0 // Signature mismatch, reject the policy update
	}

	// 2. Signature valid, compile the AST
	return InitDynamicPolicyFromJson(jsonStr)
}

/**
 * FlushAuditTraces extracts all buffered T-Number audit tickets and clears the memory.
 * This is called by the Host App to send data to the SIEM.
 */
func FlushAuditTraces() string {
	queueMutex.Lock()
	defer queueMutex.Unlock()

	if len(auditQueue) == 0 {
		return "[]"
	}

	out, _ := json.Marshal(auditQueue)
	// Clear the queue to free memory (Stateless principle)
	auditQueue = make([]AuditTrace, 0)
	
	return string(out)
}

// evalCondition performs the actual Type Assertion and AST comparison.
func evalCondition(c Condition, args map[string]interface{}) (bool, string) {
	actualValueRaw, ok := args[c.Field]
	if !ok {
		return false, fmt.Sprintf("Missing required field: %s", c.Field)
	}

	// For MVP, handle float64 (JSON numbers) and strings
	switch expectedVal := c.Value.(type) {
	case float64:
		actualVal, isNum := actualValueRaw.(float64)
		if !isNum {
			return false, fmt.Sprintf("Field %s must be numeric", c.Field)
		}
		switch c.Operator {
		case ">":
			return actualVal > expectedVal, ""
		case "<":
			return actualVal < expectedVal, ""
		case "==":
			return actualVal == expectedVal, ""
		case "!=":
			return actualVal != expectedVal, ""
		default:
			return false, fmt.Sprintf("Unsupported operator for numeric: %s", c.Operator)
		}
	case string:
		actualVal, isStr := actualValueRaw.(string)
		if !isStr {
			return false, fmt.Sprintf("Field %s must be string", c.Field)
		}
		switch c.Operator {
		case "==":
			return actualVal == expectedVal, ""
		case "!=":
			return actualVal != expectedVal, ""
		default:
			return false, fmt.Sprintf("Unsupported operator for string: %s", c.Operator)
		}
	default:
		return false, fmt.Sprintf("Unsupported rule value type for field: %s", c.Field)
	}
}

/**
 * EvaluateDynamicToolCallWithAudit evaluates a tool call against the AST engine.
 */
func EvaluateDynamicToolCallWithAudit(toolName string, argsJson string, agentId string) string {
	start := time.Now()

	dynamicMutex.RLock()
	defer dynamicMutex.RUnlock()

	status := "PASS"
	matchedRule := ""

	// Check if tool has dynamic policies
	policies, exists := dynamicPolicyMap[toolName]

	if currentOperationalMode == ModeStrictFailClosed && len(dynamicPolicyMap) == 0 {
		// Mode D: Strict Fail-Closed -> no valid policy = no execution
		status = "BLOCK"
		matchedRule = "STRICT_MODE_NO_POLICY"
	} else if !exists {
		// Fail-Closed by Default: Unrecognized tool
		if currentOperationalMode == ModeStrictFailClosed || currentOperationalMode == ModeDegraded {
			status = "BLOCK"
			matchedRule = "UNAUTHORIZED_TOOL_FAIL_CLOSED"
		}
	} else {
		var args map[string]interface{}
		err := json.Unmarshal([]byte(argsJson), &args)
		
		if err != nil {
			// Fail-Closed
			status = "BLOCK"
			matchedRule = "MALFORMED_JSON_PAYLOAD"
		} else {
			// Evaluate AST
			// Pre-check Mode C: Safe Degraded
			isWrite := false
			for _, policy := range policies {
				if policy.IsWriteAction {
					isWrite = true
					break
				}
			}

			if currentOperationalMode == ModeDegraded && isWrite {
				// Mode C: Physically intercept all write operations
				status = "BLOCK"
				matchedRule = "DEGRADED_MODE_WRITE_BLOCKED"
			} else {
				for _, policy := range policies {
					if policy.Action == "BLOCK" {
					// If there are no conditions, it's an absolute block.
					if len(policy.Conditions) == 0 {
						status = "BLOCK"
						matchedRule = "FORBIDDEN_TOOL_NAME: " + policy.ToolName
						break
					}

					// If any condition is met, trigger the BLOCK.
					// Note: Multiple conditions in a single policy act as an AND gate for the block trigger.
					// But for simplicity in Wedge, we treat any matching condition as a violation.
					// Actually, typical policy engines treat array of conditions as AND. 
					// If ALL conditions are true -> action applies.
					allMet := true
					for _, cond := range policy.Conditions {
						condMet, reason := evalCondition(cond, args)
						if reason != "" {
							// If eval fails (e.g. wrong type/missing), we fail-closed for that condition
							status = "BLOCK"
							matchedRule = "SCHEMA_VIOLATION: " + reason
							allMet = false
							break
						}
						if !condMet {
							allMet = false
							break
						}
					}

					if status == "BLOCK" {
						break // Already blocked by schema violation
					}

					if allMet {
						status = "BLOCK"
						matchedRule = "DYNAMIC_CAPABILITY_VIOLATION"
						break
					}
				}
			}
		}
	}

	latency := float64(time.Since(start).Microseconds()) / 1000.0

	snippet := argsJson
	if len(snippet) > 512 {
		snippet = snippet[:512] + "...[TRUNCATED]"
	}

	trace := AuditTrace{
		TNumber:        generateUUID(),
		Timestamp:      time.Now().UTC().Format(time.RFC3339Nano),
		AgentID:        agentId,
		Action:         "EVALUATE_DYNAMIC_TOOL_CALL",
		Status:         status,
		MatchedRule:    matchedRule,
		PayloadSnippet: toolName + " => " + snippet,
		LatencyMs:      latency,
	}

	// Queue the trace for telemetry export (Mesh Protocol)
	queueMutex.Lock()
	if len(auditQueue) < 1000 { // Bound the memory size
		auditQueue = append(auditQueue, trace)
	}
	queueMutex.Unlock()

	out, _ := json.Marshal(trace)

	// Write to append-only JSONL audit log
	writeAuditEntry(AuditLogEntry{
		Timestamp:   trace.Timestamp,
		AgentID:     agentId,
		Surface:     "NETWORK",
		Tool:        toolName,
		Decision:    status,
		PolicyEpoch: func() string {
			if dynamicPolicyMap != nil {
				return "v1" // Will be replaced with real epoch in Sprint 1 Phase 2
			}
			return "NO_POLICY"
		}(),
		Mode:        modeLabel(currentOperationalMode),
		ReasonCode:  matchedRule,
		LatencyMs:   latency,
	})

	return string(out)
}
