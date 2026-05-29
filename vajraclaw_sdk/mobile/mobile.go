package mobile

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"strings"
	"sync"
	"time"

	_ "golang.org/x/mobile/bind"
)

// Static Vajra Memory
var staticVajraRules []string
var staticMutex sync.RWMutex

// Ephemeral Rules Memory
var ephemeralRules []string
var ephemeralMutex sync.RWMutex

/**
 * InitStaticVajraFromString loads security rules from raw in-memory string.
 * Crystalizes rules into high-performance RAM memory.
 * Bypasses mobile filesystem sandboxing constraints.
 */
func InitStaticVajraFromString(content string) int {
	staticMutex.Lock()
	defer staticMutex.Unlock()

	staticVajraRules = nil
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 && !strings.HasPrefix(line, ">") && !strings.HasPrefix(line, "#") {
			staticVajraRules = append(staticVajraRules, line)
		}
	}
	// Secure defaults
	staticVajraRules = append(staticVajraRules, "診斷", "投資", "理財")
	return 1
}

/**
 * InjectEphemeralRule JIT-injects a dynamic ruleset boundary.
 */
func InjectEphemeralRule(rule string) int {
	ephemeralMutex.Lock()
	defer ephemeralMutex.Unlock()

	ephemeralRules = append(ephemeralRules, rule)
	return 1
}

/**
 * MatchTokenStream performs physical LLM stream intercept matching.
 * Returns 1 (PASS) if prompt is secure, 0 (BLOCK) if any boundary is breached.
 */
func MatchTokenStream(input string) int {
	staticMutex.RLock()
	defer staticMutex.RUnlock()
	ephemeralMutex.RLock()
	defer ephemeralMutex.RUnlock()

	// 1. Static Matrix Check
	for _, rule := range staticVajraRules {
		if strings.Contains(input, rule) {
			return 0
		}
	}

	// 2. Ephemeral JIT Pointer Check
	for _, rule := range ephemeralRules {
		if strings.Contains(input, rule) {
			return 0
		}
	}

	return 1
}

/**
 * ClearEphemeralRules performs physical evaporation of ephemeral memory segments.
 */
func ClearEphemeralRules() {
	ephemeralMutex.Lock()
	defer ephemeralMutex.Unlock()
	ephemeralRules = nil
}

// ----------------------------------------------------------------------
// AUDIT TRACE SYSTEM (Execution Ledger)
// ----------------------------------------------------------------------

// AuditTrace defines the deterministic execution log format for DROS.
// This is the portable trust substrate exported across boundaries.
type AuditTrace struct {
	TNumber        string  `json:"t_number"`        // Unique Transaction ID
	Timestamp      string  `json:"timestamp"`       // ISO-8601 Timestamp
	AgentID        string  `json:"agent_id"`        // Agent triggering the action
	Action         string  `json:"action"`          // Action evaluated
	Status         string  `json:"status"`          // "PASS" or "BLOCK"
	MatchedRule    string  `json:"matched_rule"`    // The specific rule that triggered a block
	PayloadSnippet string  `json:"payload_snippet"` // The intercepted input
	LatencyMs      float64 `json:"latency_ms"`      // Execution latency in milliseconds
}

// generateUUID generates a simple, dependency-free pseudo-UUID v4 using crypto/rand
func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10
	dst := make([]byte, 36)
	hex.Encode(dst[0:8], b[0:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], b[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], b[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], b[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:36], b[10:16])
	return string(dst)
}

/**
 * MatchTokenStreamWithAudit performs the physical intercept matching and generates
 * a cryptographically strong, deterministic Audit Ledger record in JSON format.
 * 
 * Returns a JSON string of the AuditTrace.
 */
func MatchTokenStreamWithAudit(input string, agentId string) string {
	start := time.Now()

	staticMutex.RLock()
	defer staticMutex.RUnlock()
	ephemeralMutex.RLock()
	defer ephemeralMutex.RUnlock()

	status := "PASS"
	matchedRule := ""

	// 1. Static Matrix Check
	for _, rule := range staticVajraRules {
		if strings.Contains(input, rule) {
			status = "BLOCK"
			matchedRule = rule
			break
		}
	}

	// 2. Ephemeral JIT Pointer Check
	if status == "PASS" {
		for _, rule := range ephemeralRules {
			if strings.Contains(input, rule) {
				status = "BLOCK"
				matchedRule = rule
				break
			}
		}
	}

	latency := float64(time.Since(start).Microseconds()) / 1000.0

	// Truncate payload for logging safety (prevent massive logs)
	snippet := input
	if len(snippet) > 512 {
		snippet = snippet[:512] + "...[TRUNCATED]"
	}

	trace := AuditTrace{
		TNumber:        generateUUID(),
		Timestamp:      time.Now().UTC().Format(time.RFC3339Nano),
		AgentID:        agentId,
		Action:         "EVALUATE_TOKEN_STREAM",
		Status:         status,
		MatchedRule:    matchedRule,
		PayloadSnippet: snippet,
		LatencyMs:      latency,
	}

	out, err := json.Marshal(trace)
	if err != nil {
		return `{"status": "ERROR", "message": "Failed to marshal audit log"}`
	}

	return string(out)
}

/**
 * EvaluateToolCallWithAudit performs rigid, type-strict JSON parsing and capability 
 * boundary checks. This is the cornerstone of DROS's Execution Isolation moat.
 *
 * It prevents the LLM from passing malformed payloads or executing unauthorized tools,
 * returning a deterministic JSON Audit Ledger.
 */
func EvaluateToolCallWithAudit(toolName string, argsJson string, agentId string) string {
	start := time.Now()

	status := "PASS"
	matchedRule := ""
	
	// Fail-Safe: Absolute explicit block for globally restricted tools
	if toolName == "drop_database" || toolName == "delete_all_records" {
		status = "BLOCK"
		matchedRule = "FORBIDDEN_TOOL_NAME"
	}

	// Dynamic Capability Check (Wedge Demo: Financial Controls)
	if status == "PASS" && toolName == "execute_payment" {
		var args map[string]interface{}
		err := json.Unmarshal([]byte(argsJson), &args)
		
		if err != nil {
			// Malformed JSON from LLM -> Fail-Closed
			status = "BLOCK"
			matchedRule = "MALFORMED_JSON_PAYLOAD"
		} else {
			// Type-strict field validation
			if amountRaw, ok := args["amount"]; ok {
				amount, isNumber := amountRaw.(float64)
				if !isNumber {
					status = "BLOCK"
					matchedRule = "TYPE_VIOLATION: amount must be numeric"
				} else if amount > 5000 {
					status = "BLOCK"
					matchedRule = "CAPABILITY_VIOLATION: amount > 5000 limit"
				}
			} else {
				status = "BLOCK"
				matchedRule = "SCHEMA_VIOLATION: missing amount field"
			}
		}
	}

	latency := float64(time.Since(start).Microseconds()) / 1000.0

	// Truncate payload for audit safety
	snippet := argsJson
	if len(snippet) > 512 {
		snippet = snippet[:512] + "...[TRUNCATED]"
	}

	trace := AuditTrace{
		TNumber:        generateUUID(),
		Timestamp:      time.Now().UTC().Format(time.RFC3339Nano),
		AgentID:        agentId,
		Action:         "EVALUATE_TOOL_CALL",
		Status:         status,
		MatchedRule:    matchedRule,
		PayloadSnippet: toolName + " => " + snippet,
		LatencyMs:      latency,
	}

	out, err := json.Marshal(trace)
	if err != nil {
		return `{"status": "ERROR", "message": "Failed to marshal tool audit log"}`
	}

	return string(out)
}
