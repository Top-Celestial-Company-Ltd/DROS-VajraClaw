package mobile

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// AuditLogV1 fields as per Execution Kernel Spec v1.0
// Every execution decision MUST be recorded here.
// This file is append-only and must never be truncated during runtime.
type AuditLogEntry struct {
	Timestamp    string  `json:"timestamp"`
	AgentID      string  `json:"agent_id"`
	Surface      string  `json:"surface"`       // NETWORK / SDK / ASYNC / OS_INTENT / DYNAMIC
	Tool         string  `json:"tool"`
	Decision     string  `json:"decision"`      // ALLOW / BLOCK
	PolicyEpoch  string  `json:"policy_epoch"`
	Mode         string  `json:"mode"`          // NORMAL / DEGRADED / STRICT
	ReasonCode   string  `json:"reason_code"`
	LatencyMs    float64 `json:"latency_ms"`
}

var (
	auditLogPath  string
	auditLogMutex sync.Mutex
	auditLogReady bool
)

// modeLabel converts integer mode to a readable string label
func modeLabel(mode int) string {
	switch mode {
	case ModeNormal:
		return "NORMAL"
	case ModeIsolated:
		return "ISOLATED"
	case ModeDegraded:
		return "DEGRADED"
	case ModeStrictFailClosed:
		return "STRICT"
	default:
		return "UNKNOWN"
	}
}

// ConfigureAuditLog sets the path for the append-only JSONL audit log.
// Call this once at startup. If not called, audit log is disabled.
func ConfigureAuditLog(path string) int {
	auditLogMutex.Lock()
	defer auditLogMutex.Unlock()

	// Verify we can open/create the file before committing
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[VajraClaw-Audit] ERROR: Cannot open audit log at %s: %v\n", path, err)
		return 0
	}
	f.Close()

	auditLogPath = path
	auditLogReady = true
	fmt.Printf("[VajraClaw-Audit] ✅ Audit log initialized at: %s\n", path)
	return 1
}

// writeAuditEntry appends a single JSON line to the audit log file.
// This is the ONLY function that writes to the log — guaranteeing append-only semantics.
func writeAuditEntry(entry AuditLogEntry) {
	auditLogMutex.Lock()
	defer auditLogMutex.Unlock()

	if !auditLogReady {
		return
	}

	line, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("[VajraClaw-Audit] ERROR: Failed to marshal audit entry: %v\n", err)
		return
	}

	f, err := os.OpenFile(auditLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[VajraClaw-Audit] ERROR: Cannot write audit log: %v\n", err)
		return
	}
	defer f.Close()

	// One JSON object per line (JSONL format)
	if _, err := fmt.Fprintf(f, "%s\n", line); err != nil {
		fmt.Printf("[VajraClaw-Audit] ERROR: Write failed: %v\n", err)
	}
}
