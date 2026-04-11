package utils

import (
	"api/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	verifyNowBaseURL = "https://www.verifynow.co.za"
	verifyNowClient  = &http.Client{Timeout: 45 * time.Second}
)

// VerifyBankAccountRequest is the request payload for VerifyNow bank account verification.
type VerifyBankAccountRequest struct {
	FirstName         string `json:"firstName"`
	Surname           string `json:"surname"`
	IdentityNumber    string `json:"identityNumber"`
	IdentityType      string `json:"identityType"`
	BankAccountNumber string `json:"bankAccountNumber"`
	BankBranchCode    string `json:"bankBranchCode"`
	BankAccountType   string `json:"bankAccountType"`
	Mode              string `json:"mode"`
}

// VerifyBankAccountResponse is the top-level response from VerifyNow.
type VerifyBankAccountResponse struct {
	Success   bool                        `json:"success"`
	RequestID string                      `json:"requestId"`
	Service   string                      `json:"service"`
	Results   VerifyBankAccountResults    `json:"results"`
}

// VerifyBankAccountResults contains the verification outcome.
type VerifyBankAccountResults struct {
	IdentityAndAccountVerified bool                          `json:"identity_and_account_verified"`
	Summary                    string                        `json:"summary"`
	VerificationResults        BankVerificationResultDetails `json:"verification_results"`
}

// BankVerificationResultDetails contains granular verification checks.
type BankVerificationResultDetails struct {
	Status           string `json:"Status"`
	AccountFound     string `json:"accountFound"`
	AccountOpen      string `json:"accountOpen"`
	IdentityMatch    string `json:"identityMatch"`
	AccountTypeMatch string `json:"accountTypeMatch"`
	AcceptsCredits   string `json:"acceptsCredits"`
	AcceptsDebits    string `json:"acceptsDebits"`
}

// VerifyBankAccount calls the VerifyNow API to verify a South African bank account.
func VerifyBankAccount(req VerifyBankAccountRequest) (*VerifyBankAccountResponse, error) {
	if config.VerifyNowApiKey == "" {
		return nil, fmt.Errorf("VERIFYNOW_API_KEY is not configured")
	}

	if req.IdentityType == "" {
		req.IdentityType = "IDNumber"
	}
	if req.Mode == "" {
		req.Mode = config.VerifyNowMode
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("verifynow marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/api/external/bank-account-verification", verifyNowBaseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("verifynow request error: %w", err)
	}
	httpReq.Header.Set("x-api-key", config.VerifyNowApiKey)
	httpReq.Header.Set("Idempotency-Key", uuid.New().String())
	httpReq.Header.Set("Content-Type", "application/json")

	log.Printf("verifynow: calling bank-account-verification for ID %s", req.IdentityNumber)

	resp, err := verifyNowClient.Do(httpReq)
	if err != nil {
		log.Printf("verifynow: request failed: %v", err)
		return nil, fmt.Errorf("verifynow API unreachable: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("verifynow read body error: %w", err)
	}

	log.Printf("verifynow: status=%d body=%s", resp.StatusCode, string(respBody))

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("verifynow API: invalid API key")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("verifynow API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result VerifyBankAccountResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("verifynow response decode error: %w", err)
	}
	return &result, nil
}
