package mobile

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMachineID(t *testing.T) {
	id1, err := GetOrCreateMachineID()
	if err != nil {
		t.Fatalf("Failed to generate Machine ID: %v", err)
	}

	if len(id1) == 0 {
		t.Fatalf("Machine ID is empty")
	}

	id2, err := GetOrCreateMachineID()
	if err != nil {
		t.Fatalf("Failed to retrieve Machine ID: %v", err)
	}

	if id1 != id2 {
		t.Fatalf("Machine ID should be persistent. Expected %s, got %s", id1, id2)
	}
}

func TestLicenseValidationSuccess(t *testing.T) {
	// Generate a temporary keypair for the test
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	OfficialPubKeyHex = hex.EncodeToString(pub)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := "startup|5|0"
		sig := ed25519.Sign(priv, []byte(payload))

		resp := LicenseVerifyResponse{
			LicenseTier: "startup",
			Concurrency: 5,
			ExpiresAt:   0,
			Signature:   hex.EncodeToString(sig),
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	os.Setenv("VAJRA_LICENSE_API", mockServer.URL)
	defer os.Unsetenv("VAJRA_LICENSE_API")

	// Reset state
	SetLicenseStatus(StatusExpired, time.Time{})

	isValid := ValidateLicense("TEST_KEY")
	if !isValid {
		t.Fatalf("Expected valid license, but got invalid")
	}

	if GetCurrentLicenseStatus() != StatusValid {
		t.Fatalf("Expected state to be StatusValid")
	}
}

func TestLicenseValidationSignatureFailure(t *testing.T) {
	// Generate a temporary keypair for the test
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	OfficialPubKeyHex = hex.EncodeToString(pub)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrong payload -> invalid signature
		payload := "hacker|2|0"
		sig := ed25519.Sign(priv, []byte(payload))

		resp := LicenseVerifyResponse{
			LicenseTier: "startup", // Tampered tier
			Concurrency: 5,
			ExpiresAt:   0,
			Signature:   hex.EncodeToString(sig),
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	os.Setenv("VAJRA_LICENSE_API", mockServer.URL)
	defer os.Unsetenv("VAJRA_LICENSE_API")

	// Reset state
	SetLicenseStatus(StatusExpired, time.Time{})

	isValid := ValidateLicense("TEST_KEY")
	if isValid {
		t.Fatalf("Expected invalid license due to signature mismatch")
	}

	if GetCurrentLicenseStatus() != StatusRevoked {
		t.Fatalf("Expected state to be StatusRevoked due to spoofing attempt")
	}
}

func TestLicenseValidationGracePeriod(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	os.Setenv("VAJRA_LICENSE_API", mockServer.URL)
	defer os.Unsetenv("VAJRA_LICENSE_API")

	// Scenario 1: Last success was 2 days ago (within grace period)
	SetLicenseStatus(StatusValid, time.Now().Add(-48*time.Hour))
	isValid := ValidateLicense("TEST_KEY")
	if !isValid {
		t.Fatalf("Expected valid license (GRACE), but got invalid")
	}
	if GetCurrentLicenseStatus() != StatusGrace {
		t.Fatalf("Expected state to be StatusGrace")
	}

	// Scenario 2: Last success was 8 days ago (expired grace period)
	SetLicenseStatus(StatusValid, time.Now().Add(-8*24*time.Hour))
	isValid = ValidateLicense("TEST_KEY")
	if isValid {
		t.Fatalf("Expected invalid license (EXPIRED), but got valid")
	}
	if GetCurrentLicenseStatus() != StatusExpired {
		t.Fatalf("Expected state to be StatusExpired")
	}
}
