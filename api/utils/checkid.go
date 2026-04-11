package utils

import (
	"api/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	checkIDBaseURL  = "https://api.checkid.co.za"
	checkIDClient   = &http.Client{Timeout: 10 * time.Second}
	checkIDWarnOnce sync.Once
)

type checkIDResult struct {
	IDNumber string `json:"idNumber"`
	IsValid  bool   `json:"isValid"`
}

// CheckIDValidate validates a single SA ID number against the CheckID API.
// Returns (isValid, error). If the API key is not configured, returns (true, nil)
// so that local Luhn validation remains the only check.
func CheckIDValidate(idNumber string) (bool, error) {
	if config.CheckIDApiKey == "" {
		checkIDWarnOnce.Do(func() {
			log.Println("WARNING: CHECKID_API_KEY not set — skipping CheckID API validation")
		})
		return true, nil
	}

	url := fmt.Sprintf("%s/api/v1/validate/%s", checkIDBaseURL, idNumber)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("checkid request error: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.CheckIDApiKey)

	resp, err := checkIDClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("checkid API unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return false, fmt.Errorf("checkid API: invalid API key")
	}
	if resp.StatusCode == http.StatusBadRequest {
		return false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("checkid API returned status %d", resp.StatusCode)
	}

	var result checkIDResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("checkid response decode error: %w", err)
	}
	return result.IsValid, nil
}

// CheckIDBulkValidate validates up to 100 SA ID numbers at a time against the
// CheckID bulk API. Returns a map of idNumber -> isValid.
func CheckIDBulkValidate(idNumbers []string) (map[string]bool, error) {
	if config.CheckIDApiKey == "" {
		checkIDWarnOnce.Do(func() {
			log.Println("WARNING: CHECKID_API_KEY not set — skipping CheckID API validation")
		})
		results := make(map[string]bool, len(idNumbers))
		for _, id := range idNumbers {
			results[id] = true
		}
		return results, nil
	}

	results := make(map[string]bool, len(idNumbers))

	// Process in chunks of 100 (API limit)
	for start := 0; start < len(idNumbers); start += 100 {
		end := start + 100
		if end > len(idNumbers) {
			end = len(idNumbers)
		}
		chunk := idNumbers[start:end]

		body, err := json.Marshal(chunk)
		if err != nil {
			return nil, fmt.Errorf("checkid bulk marshal error: %w", err)
		}

		url := fmt.Sprintf("%s/api/v1/validate/bulk", checkIDBaseURL)
		req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
		if err != nil {
			return nil, fmt.Errorf("checkid bulk request error: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+config.CheckIDApiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := checkIDClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("checkid API unreachable: %w", err)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()
			return nil, fmt.Errorf("checkid API: invalid API key")
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("checkid API returned status %d", resp.StatusCode)
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("checkid bulk read error: %w", err)
		}

		var chunkResults []checkIDResult
		if err := json.Unmarshal(respBody, &chunkResults); err != nil {
			return nil, fmt.Errorf("checkid bulk decode error: %w", err)
		}

		for _, r := range chunkResults {
			results[r.IDNumber] = r.IsValid
		}
	}

	return results, nil
}

// ValidateRSAID validates a single SA ID using local Luhn check first,
// then the CheckID API. Returns (isValid, error).
func ValidateRSAID(idNumber string) (bool, error) {
	if !IsValidRSAID(idNumber) {
		return false, nil
	}
	return CheckIDValidate(idNumber)
}

// ValidateRSAIDsBulk validates a batch of SA IDs: local Luhn pre-filter,
// then bulk CheckID API call for those that pass locally.
// Returns a map of idNumber -> isValid.
func ValidateRSAIDsBulk(idNumbers []string) (map[string]bool, error) {
	results := make(map[string]bool, len(idNumbers))
	var passedLocal []string

	for _, id := range idNumbers {
		if IsValidRSAID(id) {
			passedLocal = append(passedLocal, id)
		} else {
			results[id] = false
		}
	}

	if len(passedLocal) == 0 {
		return results, nil
	}

	apiResults, err := CheckIDBulkValidate(passedLocal)
	if err != nil {
		return nil, err
	}

	for id, valid := range apiResults {
		results[id] = valid
	}

	return results, nil
}
