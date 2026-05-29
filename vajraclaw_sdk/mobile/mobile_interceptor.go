package mobile

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
)

// VajraClawRoundTripper is an http.RoundTripper implementation that enforces
// the Execution Path Uniqueness rule. All outbound requests made by the Agent
// using this RoundTripper will be physically inspected by VajraClaw.
// If the payload contains an unauthorized tool call, the connection is melted.
type VajraClawRoundTripper struct {
	Proxied http.RoundTripper
	AgentID string
}

// Default creation
func NewVajraClawRoundTripper(proxied http.RoundTripper, agentID string) *VajraClawRoundTripper {
	if proxied == nil {
		proxied = http.DefaultTransport
	}
	return &VajraClawRoundTripper{
		Proxied: proxied,
		AgentID: agentID,
	}
}

func (vrt *VajraClawRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// 1. If it's a GET request, it's generally a read operation. We might still check URL params,
	// but for this example, we focus on POST/PUT where tool call payloads usually live.
	if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH" {
		if req.Body != nil {
			// Read the body for inspection
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			
			// Restore the body for the actual request if it passes
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			bodyStr := string(bodyBytes)
			
			// Simple heuristic to extract tool name from JSON payload
			// In a real environment, this would parse the specific LLM provider's format (e.g. OpenAI tool_calls)
			// Here we assume the body is the JSON payload of the tool call itself and try to extract "tool_name"
			// Or we just evaluate the raw body against all known dynamic tools
			
			// Physical Melt Check
			// To enforce Fail-Closed, if the current mode is StrictFailClosed and we can't parse the tool,
			// we should theoretically block it. For demonstration, we simulate evaluating the body.
			
			// A mock extraction of tool_name for demo:
			toolName := "unknown"
			if strings.Contains(bodyStr, `"tool_name":"transfer_funds"`) {
				toolName = "transfer_funds"
			} else if strings.Contains(bodyStr, `"tool_name":"query_faq"`) {
				toolName = "query_faq"
			}

			// Perform actual VajraClaw evaluation
			auditTraceJSON := EvaluateDynamicToolCallWithAudit(toolName, bodyStr, vrt.AgentID)
			
			// If the audit trace indicates a BLOCK, we physically sever the connection here.
			if strings.Contains(auditTraceJSON, `"status":"BLOCK"`) {
				// PHYSICAL MELT - Execution Path Uniqueness guaranteed.
				// The request never leaves the OS, returning a hard network error.
				return nil, errors.New("[VajraClaw] FATAL: Tool invocation intercepted and blocked by Physical Melt. Connection severed.")
			}
		}
	}

	// 2. Pass through to the real network if authorized
	return vrt.Proxied.RoundTrip(req)
}
