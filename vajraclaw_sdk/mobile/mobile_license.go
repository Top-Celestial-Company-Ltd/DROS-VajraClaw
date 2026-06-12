package mobile

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LicenseStatus int

const (
	StatusValid LicenseStatus = iota
	StatusGrace
	StatusExpired
	StatusRevoked
)

var (
	currentLicenseStatus LicenseStatus = StatusExpired
	lastSuccessTime      time.Time
	licenseMutex         sync.RWMutex
	// In production, this would be the actual public key from DROS
	OfficialPubKeyHex = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" 
)

// GracePeriod is 7 days
const GracePeriod = 7 * 24 * time.Hour

func GetOrCreateMachineID() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to temp dir
		homeDir = os.TempDir()
	}
	
	vajraDir := filepath.Join(homeDir, ".vajra")
	if err := os.MkdirAll(vajraDir, 0755); err != nil {
		return "", err
	}
	
	idPath := filepath.Join(vajraDir, "machine.id")
	if data, err := os.ReadFile(idPath); err == nil {
		return string(bytes.TrimSpace(data)), nil
	}
	
	// Generate random ID
	b := make([]byte, 16)
	rand.Read(b)
	newID := hex.EncodeToString(b)
	
	if err := os.WriteFile(idPath, []byte(newID), 0600); err != nil {
		return "", err
	}
	
	return newID, nil
}

type LicenseVerifyRequest struct {
	LicenseKey string `json:"license_key"`
	MachineID  string `json:"machine_id"`
	Version    string `json:"version"`
}

type LicenseVerifyResponse struct {
	LicenseTier string `json:"license_tier"`
	Concurrency int    `json:"concurrency"`
	ExpiresAt   int64  `json:"expires_at"`
	Signature   string `json:"signature"`
}

func ValidateLicense(licenseKey string) bool {
	if licenseKey == "" {
		licenseKey = "TRIAL"
		fmt.Println("[VajraClaw] No License Key provided. Activating 30-Day Auto-Trial via Heartbeat Server.")
	}

	machineID, err := GetOrCreateMachineID()
	if err != nil {
		fmt.Printf("[VajraClaw] Failed to get machine ID: %v\n", err)
		return false
	}
	
	apiURL := os.Getenv("VAJRA_LICENSE_API")
	if apiURL == "" {
		apiURL = "https://api.dr-os.io/license/verify"
	}
	
	reqBody := LicenseVerifyRequest{
		LicenseKey: licenseKey,
		MachineID:  machineID,
		Version:    "v1.0",
	}
	
	jsonBody, _ := json.Marshal(reqBody)
	
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(jsonBody))
	
	licenseMutex.Lock()
	defer licenseMutex.Unlock()
	
	if err != nil || resp.StatusCode != 200 {
		// Network error or 5xx, check Grace Period
		if !lastSuccessTime.IsZero() && time.Since(lastSuccessTime) < GracePeriod {
			currentLicenseStatus = StatusGrace
			fmt.Printf("[VajraClaw] License Server unreachable. Entering GRACE period. Valid for %v more.\n", GracePeriod - time.Since(lastSuccessTime))
			return true
		}
		
		if resp != nil && resp.StatusCode == 403 {
			// Explicitly rejected
			currentLicenseStatus = StatusRevoked
			fmt.Println("[VajraClaw] License Revoked by server!")
			return false
		}
		
		currentLicenseStatus = StatusExpired
		fmt.Println("[VajraClaw] License check failed and no valid grace period.")
		return false
	}
	
	defer resp.Body.Close()
	var apiResp LicenseVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Println("[VajraClaw] Failed to decode API response.")
		return false
	}
	
	// Verify Signature
	if apiResp.LicenseTier == "Trial" {
		// Trial tier does not require strict public key signature validation for the demo
		fmt.Println("[VajraClaw] Trial tier detected. Bypassing strict signature validation.")
	} else {
		pubBytes, err := hex.DecodeString(OfficialPubKeyHex)
		if err == nil && len(pubBytes) == ed25519.PublicKeySize {
			sigBytes, err := hex.DecodeString(apiResp.Signature)
			if err != nil || len(sigBytes) != ed25519.SignatureSize {
				currentLicenseStatus = StatusRevoked
				fmt.Println("[VajraClaw] Invalid signature format from API.")
				return false
			}
			
			// The payload we expect the server to sign is the tier, concurrency, and expires_at combined
			payload := fmt.Sprintf("%s|%d|%d", apiResp.LicenseTier, apiResp.Concurrency, apiResp.ExpiresAt)
			if !ed25519.Verify(ed25519.PublicKey(pubBytes), []byte(payload), sigBytes) {
				currentLicenseStatus = StatusRevoked
				fmt.Println("[VajraClaw] FATAL: License API signature verification failed! Possible spoofing attack.")
				return false
			}
		} else {
			// If we can't parse our own pubkey, something is very wrong with the build
			currentLicenseStatus = StatusRevoked
			return false
		}
	}
	
	if time.Now().Unix() > apiResp.ExpiresAt && apiResp.ExpiresAt != 0 {
		currentLicenseStatus = StatusExpired
		fmt.Println("[VajraClaw] License has expired.")
		return false
	}
	
	currentLicenseStatus = StatusValid
	lastSuccessTime = time.Now()
	fmt.Printf("[VajraClaw] License Verified. Tier: %s\n", apiResp.LicenseTier)
	
	return true
}

func GetCurrentLicenseStatus() LicenseStatus {
	licenseMutex.RLock()
	defer licenseMutex.RUnlock()
	
	// Re-evaluate grace period expiration lazily
	if currentLicenseStatus == StatusGrace {
		if time.Since(lastSuccessTime) > GracePeriod {
			return StatusExpired
		}
	}
	return currentLicenseStatus
}

// For testing purposes
func SetLicenseStatus(status LicenseStatus, lastSuccess time.Time) {
	licenseMutex.Lock()
	defer licenseMutex.Unlock()
	currentLicenseStatus = status
	lastSuccessTime = lastSuccess
}
