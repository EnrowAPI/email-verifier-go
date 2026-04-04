package emailverifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.enrow.io"

// VerificationResult represents the response from a single email verification.
type VerificationResult struct {
	ID            string `json:"id"`
	Email         string `json:"email,omitempty"`
	Qualification string `json:"qualification,omitempty"`
	Status        string `json:"status,omitempty"`
	Message       string `json:"message,omitempty"`
	CreditsUsed   int    `json:"creditsUsed,omitempty"`
}

// BulkVerificationResult represents the response from initiating a bulk verification.
type BulkVerificationResult struct {
	BatchID     string `json:"batchId"`
	Total       int    `json:"total"`
	Status      string `json:"status"`
	CreditsUsed int    `json:"creditsUsed,omitempty"`
}

// BulkVerificationResults represents the response from retrieving bulk verification results.
type BulkVerificationResults struct {
	BatchID     string               `json:"batchId"`
	Status      string               `json:"status"`
	Total       int                  `json:"total"`
	Completed   int                  `json:"completed,omitempty"`
	CreditsUsed int                  `json:"creditsUsed,omitempty"`
	Results     []VerificationResult `json:"results,omitempty"`
}

func request(apiKey, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Message != "" {
			return nil, fmt.Errorf("%s", errResp.Message)
		}
		return nil, fmt.Errorf("API error %d", resp.StatusCode)
	}

	return respBody, nil
}

type settings struct {
	Webhook string `json:"webhook,omitempty"`
}

// VerifyEmail starts a single email verification and returns the initial result.
func VerifyEmail(apiKey, email string, webhook string) (*VerificationResult, error) {
	body := map[string]interface{}{
		"email": email,
	}
	if webhook != "" {
		body["settings"] = settings{Webhook: webhook}
	}

	data, err := request(apiKey, http.MethodPost, "/email/verify/single", body)
	if err != nil {
		return nil, err
	}

	var result VerificationResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

// GetVerificationResult retrieves the result of a previously started single email verification.
func GetVerificationResult(apiKey, id string) (*VerificationResult, error) {
	data, err := request(apiKey, http.MethodGet, "/email/verify/single?id="+id, nil)
	if err != nil {
		return nil, err
	}

	var result VerificationResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

// VerifyEmails starts a bulk email verification and returns the batch info.
func VerifyEmails(apiKey string, emails []string, webhook string) (*BulkVerificationResult, error) {
	verifications := make([]map[string]string, len(emails))
	for i, email := range emails {
		verifications[i] = map[string]string{"email": email}
	}

	body := map[string]interface{}{
		"verifications": verifications,
	}
	if webhook != "" {
		body["settings"] = settings{Webhook: webhook}
	}

	data, err := request(apiKey, http.MethodPost, "/email/verify/bulk", body)
	if err != nil {
		return nil, err
	}

	var result BulkVerificationResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

// GetVerificationResults retrieves the results of a previously started bulk email verification.
func GetVerificationResults(apiKey, id string) (*BulkVerificationResults, error) {
	data, err := request(apiKey, http.MethodGet, "/email/verify/bulk?id="+id, nil)
	if err != nil {
		return nil, err
	}

	var result BulkVerificationResults
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}
