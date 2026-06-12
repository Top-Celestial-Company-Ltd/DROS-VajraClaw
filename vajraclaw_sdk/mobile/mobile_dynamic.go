package mobile

import (
	"bytes"
	"crypto/ed25519"
	"encoding/binary"
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
	Epoch        string       `json:"epoch"`
	ToolPolicies []ToolPolicy `json:"tool_policies"`
}

// Global Memory for dynamic policies and audit queue
var (
	dynamicPolicyMap map[string][]ToolPolicy // Legacy V1
	
	// V2 O(1) Capability Bitmap Memory
	v2AgentIndexMap map[string]int
	v2ToolIndexMap  map[string]int
	v2CapabilityBitmap []byte
	isV2Binary bool

	currentPolicyEpoch string = "NO_POLICY"
	currentPolicyHash  string = ""
	currentCompilerVersion string = ""
	currentDslVersion  int = 0
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

	if policy.Epoch != "" {
		currentPolicyEpoch = policy.Epoch
	} else {
		currentPolicyEpoch = "UNVERSIONED"
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
	v2AgentIndexMap = nil
	v2ToolIndexMap = nil
	v2CapabilityBitmap = nil
	isV2Binary = false
	currentPolicyEpoch = "NO_POLICY"
	currentPolicyHash = ""
	currentCompilerVersion = ""
	currentDslVersion = 0
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
 * InitDynamicPolicyFromBinary loads a cryptographically signed binary matrix (.bin)
 * Format:
 *  - 6 bytes Header: "VAJRAC"
 *  - 1 byte Version: 0x01
 *  - 32 bytes Epoch ID (null-padded)
 *  - 64 bytes Ed25519 Signature
 *  - 4 bytes Payload Length (BigEndian uint32)
 *  - Variable length JSON Payload bytes
 */
func InitDynamicPolicyFromBinary(binBytes []byte, pubKeyHex string) int {
	pubKey, err := hex.DecodeString(pubKeyHex)
	if err != nil || len(pubKey) != ed25519.PublicKeySize {
		if currentOperationalMode == ModeStrictFailClosed {
			panic("[VajraClaw] FATAL: Invalid Public Key format in Strict Mode")
		}
		return 0
	}

	if len(binBytes) < (6 + 1 + 32 + 64 + 4) {
		return 0 // Too short
	}

	header := binBytes[0:6]
	if string(header) != "VAJRAC" {
		return 0 // Invalid magic header
	}

	version := binBytes[6]
	if version != 0x01 && version != 0x02 && version != 0x03 {
		return 0 // Unsupported version
	}

	sig := binBytes[39 : 39+64]
	
	var payload []byte
	if version == 0x01 || version == 0x02 {
		payloadLen := binary.BigEndian.Uint32(binBytes[103 : 107])
		if len(binBytes) < 107+int(payloadLen) {
			return 0 // Truncated payload
		}
		payload = binBytes[107 : 107+payloadLen]
	} else if version == 0x03 {
		payloadLen := binary.BigEndian.Uint32(binBytes[152 : 156])
		if len(binBytes) < 156+int(payloadLen) {
			return 0 // Truncated payload
		}
		payload = binBytes[156 : 156+payloadLen]
	}

	if !ed25519.Verify(ed25519.PublicKey(pubKey), payload, sig) {
		if currentOperationalMode == ModeStrictFailClosed {
			panic("[VajraClaw] FATAL: Cryptographic Signature Verification Failed! Tampered Policy detected.")
		}
		return 0
	}

	if version == 0x01 {
		return InitDynamicPolicyFromJson(string(payload))
	} else if version == 0x02 {
		return initDynamicPolicyFromV2Payload(payload, binBytes[7:39], "", "", 0)
	} else if version == 0x03 {
		hashHex := fmt.Sprintf("%x", binBytes[103:135])
		dslVer := int(binBytes[135])
		cvStr := string(bytes.TrimRight(binBytes[136:152], "\x00"))
		return initDynamicPolicyFromV2Payload(payload, binBytes[7:39], hashHex, cvStr, dslVer)
	}
	return 0
}

func initDynamicPolicyFromV2Payload(payload []byte, epochBytes []byte, pHash string, cVer string, dVer int) int {
	dynamicMutex.Lock()
	defer dynamicMutex.Unlock()

	if len(payload) < 4 {
		return 0
	}

	numAgents := int(binary.BigEndian.Uint16(payload[0:2]))
	numTools := int(binary.BigEndian.Uint16(payload[2:4]))

	v2AgentIndexMap = make(map[string]int)
	v2ToolIndexMap = make(map[string]int)

	offset := 4
	
	// Read Agents
	for i := 0; i < numAgents; i++ {
		start := offset
		for offset < len(payload) && payload[offset] != 0 {
			offset++
		}
		if offset >= len(payload) {
			return 0
		}
		v2AgentIndexMap[string(payload[start:offset])] = i
		offset++ // Skip null byte
	}

	// Read Tools
	for i := 0; i < numTools; i++ {
		start := offset
		for offset < len(payload) && payload[offset] != 0 {
			offset++
		}
		if offset >= len(payload) {
			return 0
		}
		v2ToolIndexMap[string(payload[start:offset])] = i
		offset++ // Skip null byte
	}

	// Read Bitmap
	expectedBitmapSize := (numAgents * numTools + 7) / 8
	if len(payload)-offset < expectedBitmapSize {
		return 0
	}
	
	v2CapabilityBitmap = make([]byte, expectedBitmapSize)
	copy(v2CapabilityBitmap, payload[offset:offset+expectedBitmapSize])

	isV2Binary = true
	
	// Trim null bytes from epoch
	epochStr := string(epochBytes)
	for i := 0; i < len(epochStr); i++ {
		if epochStr[i] == 0 {
			epochStr = epochStr[:i]
			break
		}
	}
	currentPolicyEpoch = epochStr
	currentPolicyHash = pHash
	currentCompilerVersion = cVer
	currentDslVersion = dVer

	return 1
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
 * Requires an expectedEpoch to enforce Policy Epoch Lock (Anti-drift).
 */
func EvaluateDynamicToolCallWithAudit(toolName string, argsJson string, agentId string, expectedEpoch string) string {
	// 1. License Check
	licStatus := GetCurrentLicenseStatus()
	if licStatus == StatusExpired || licStatus == StatusRevoked {
		if currentOperationalMode == ModeStrictFailClosed {
			panic(fmt.Sprintf("[VajraClaw] FATAL: License %v. Halting execution.", licStatus))
		}
		trace := AuditTrace{
			TNumber:        generateUUID(),
			Status:         "BLOCK",
			Timestamp:      time.Now().Format(time.RFC3339),
			AgentID:        agentId,
			Action:         toolName,
			PayloadSnippet: argsJson,
			MatchedRule:    "LICENSE_REQUIRED",
			LatencyMs:      0.0,
		}
		out, _ := json.Marshal(trace)
		return string(out)
	}

	// 2. Epoch Anti-Drift Check

	start := time.Now()

	dynamicMutex.RLock()
	defer dynamicMutex.RUnlock()

	status := "PASS"
	matchedRule := ""

	// Check if tool has dynamic policies
	policies, exists := dynamicPolicyMap[toolName]

	if expectedEpoch != "" && expectedEpoch != currentPolicyEpoch && currentPolicyEpoch != "NO_POLICY" {
		// Epoch drift detected
		status = "BLOCK"
		matchedRule = "EPOCH_MISMATCH"
	} else if currentOperationalMode == ModeStrictFailClosed && !isV2Binary && len(dynamicPolicyMap) == 0 {
		status = "BLOCK"
		matchedRule = "STRICT_MODE_NO_POLICY"
	} else if isV2Binary && matchedRule == "" {
		// ==========================================
		// V2 O(1) BITMAP LOOKUP (Zero JSON Parsing)
		// ==========================================
		agentIdx, aOk := v2AgentIndexMap[agentId]
		toolIdx, tOk := v2ToolIndexMap[toolName]

		if !aOk {
			status = "BLOCK"
			matchedRule = "V2_UNKNOWN_AGENT"
		} else if !tOk {
			status = "BLOCK"
			matchedRule = "V2_UNKNOWN_TOOL"
		} else {
			numTools := len(v2ToolIndexMap)
			bitOffset := agentIdx*numTools + toolIdx
			byteIdx := bitOffset / 8
			bitInByte := bitOffset % 8

			if (v2CapabilityBitmap[byteIdx] & (1 << bitInByte)) == 0 {
				status = "BLOCK"
				matchedRule = "V2_CAPABILITY_DENIED"
			}
		}
	} else if !exists && matchedRule == "" {
		// ==========================================
		// LEGACY V1 (JSON Parsing Fallback)
		// ==========================================
		// Fail-Closed by Default: Unrecognized tool
		if currentOperationalMode == ModeStrictFailClosed || currentOperationalMode == ModeDegraded {
			status = "BLOCK"
			matchedRule = "UNAUTHORIZED_TOOL_FAIL_CLOSED"
		}
	} else if matchedRule == "" {
		var args map[string]interface{}
		err := json.Unmarshal([]byte(argsJson), &args)
		
		if err != nil {
			status = "BLOCK"
			matchedRule = "MALFORMED_JSON_PAYLOAD"
		} else {
			isWrite := false
			for _, policy := range policies {
				if policy.IsWriteAction {
					isWrite = true
					break
				}
			}

			if currentOperationalMode == ModeDegraded && isWrite {
				status = "BLOCK"
				matchedRule = "DEGRADED_MODE_WRITE_BLOCKED"
			} else {
				for _, policy := range policies {
					if policy.Action == "BLOCK" {
						if len(policy.Conditions) == 0 {
							status = "BLOCK"
							matchedRule = "FORBIDDEN_TOOL_NAME: " + policy.ToolName
							break
						}

						allMet := true
						for _, cond := range policy.Conditions {
							condMet, reason := evalCondition(cond, args)
							if reason != "" {
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
							break 
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
		PolicyHash:     currentPolicyHash,
		CompilerVersion: currentCompilerVersion,
		DslVersion:     currentDslVersion,
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
				return currentPolicyEpoch 
			}
			return "NO_POLICY"
		}(),
		Mode:        modeLabel(currentOperationalMode),
		ReasonCode:  matchedRule,
		LatencyMs:   latency,
	})

	return string(out)
}
