package etherscan

import (
	"crypto-acc-tracking/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	BaseURL        = "https://api.etherscan.io/api"
	DefaultTimeout = 30 * time.Second
	MaxRetries     = 3
	RetryDelay     = 5 * time.Second
)

// Client represents the Etherscan API client
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// New creates a new Etherscan client
func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL: BaseURL,
	}
}

// GetNormalTransactions fetches normal transactions for an address
func (c *Client) GetNormalTransactions(address string, startBlock, endBlock int, page, offset int) ([]models.EtherscanNormalTx, error) {
	params := url.Values{
		"module":     []string{"account"},
		"action":     []string{"txlist"},
		"address":    []string{address},
		"startblock": []string{strconv.Itoa(startBlock)},
		"endblock":   []string{strconv.Itoa(endBlock)},
		"page":       []string{strconv.Itoa(page)},
		"offset":     []string{strconv.Itoa(offset)},
		"sort":       []string{"desc"},
	}

	if c.apiKey != "" {
		params.Set("apikey", c.apiKey)
	}

	var response models.EtherscanResponse
	if err := c.makeRequest(params, &response); err != nil {
		return nil, err
	}

	if response.Status != "1" {
		if response.Message == "NOTOK" {
			return nil, fmt.Errorf("Etherscan API error: Rate limit exceeded or invalid request. Consider using an API key with -k flag")
		}
		return nil, fmt.Errorf("Etherscan API error: %s", response.Message)
	}

	// Convert interface{} to []models.EtherscanNormalTx
	resultBytes, err := json.Marshal(response.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	var transactions []models.EtherscanNormalTx
	if err := json.Unmarshal(resultBytes, &transactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transactions: %w", err)
	}

	return transactions, nil
}

// GetInternalTransactions fetches internal transactions for an address
func (c *Client) GetInternalTransactions(address string, startBlock, endBlock int, page, offset int) ([]models.EtherscanInternalTx, error) {
	params := url.Values{
		"module":     []string{"account"},
		"action":     []string{"txlistinternal"},
		"address":    []string{address},
		"startblock": []string{strconv.Itoa(startBlock)},
		"endblock":   []string{strconv.Itoa(endBlock)},
		"page":       []string{strconv.Itoa(page)},
		"offset":     []string{strconv.Itoa(offset)},
		"sort":       []string{"desc"},
	}

	if c.apiKey != "" {
		params.Set("apikey", c.apiKey)
	}

	var response models.EtherscanResponse
	if err := c.makeRequest(params, &response); err != nil {
		return nil, err
	}

	if response.Status != "1" {
		if response.Message == "NOTOK" {
			return nil, fmt.Errorf("Etherscan API error: Rate limit exceeded or invalid request. Consider using an API key with -k flag")
		}
		if response.Message == "No transactions found" {
			return []models.EtherscanInternalTx{}, nil // Return empty slice for no transactions
		}
		return nil, fmt.Errorf("Etherscan API error: %s", response.Message)
	}

	// Convert interface{} to []models.EtherscanInternalTx
	resultBytes, err := json.Marshal(response.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	var transactions []models.EtherscanInternalTx
	if err := json.Unmarshal(resultBytes, &transactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal internal transactions: %w", err)
	}

	return transactions, nil
}

// GetTokenTransactions fetches ERC-20 token transactions for an address
func (c *Client) GetTokenTransactions(address string, startBlock, endBlock int, page, offset int) ([]models.EtherscanTokenTx, error) {
	params := url.Values{
		"module":     []string{"account"},
		"action":     []string{"tokentx"},
		"address":    []string{address},
		"startblock": []string{strconv.Itoa(startBlock)},
		"endblock":   []string{strconv.Itoa(endBlock)},
		"page":       []string{strconv.Itoa(page)},
		"offset":     []string{strconv.Itoa(offset)},
		"sort":       []string{"desc"},
	}

	if c.apiKey != "" {
		params.Set("apikey", c.apiKey)
	}

	var response models.EtherscanResponse
	if err := c.makeRequest(params, &response); err != nil {
		return nil, err
	}

	if response.Status != "1" {
		if response.Message == "NOTOK" {
			return nil, fmt.Errorf("Etherscan API error: Rate limit exceeded or invalid request. Consider using an API key with -k flag")
		}
		if response.Message == "No transactions found" {
			return []models.EtherscanTokenTx{}, nil // Return empty slice for no transactions
		}
		return nil, fmt.Errorf("Etherscan API error: %s", response.Message)
	}

	// Convert interface{} to []models.EtherscanTokenTx
	resultBytes, err := json.Marshal(response.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	var transactions []models.EtherscanTokenTx
	if err := json.Unmarshal(resultBytes, &transactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token transactions: %w", err)
	}

	return transactions, nil
}

// GetNFTTransactions fetches ERC-721 NFT transactions for an address
func (c *Client) GetNFTTransactions(address string, startBlock, endBlock int, page, offset int) ([]models.EtherscanNFTTx, error) {
	params := url.Values{
		"module":     []string{"account"},
		"action":     []string{"tokennfttx"},
		"address":    []string{address},
		"startblock": []string{strconv.Itoa(startBlock)},
		"endblock":   []string{strconv.Itoa(endBlock)},
		"page":       []string{strconv.Itoa(page)},
		"offset":     []string{strconv.Itoa(offset)},
		"sort":       []string{"desc"},
	}

	if c.apiKey != "" {
		params.Set("apikey", c.apiKey)
	}

	var response models.EtherscanResponse
	if err := c.makeRequest(params, &response); err != nil {
		return nil, err
	}

	if response.Status != "1" {
		if response.Message == "NOTOK" {
			return nil, fmt.Errorf("Etherscan API error: Rate limit exceeded or invalid request. Consider using an API key with -k flag")
		}
		if response.Message == "No transactions found" {
			return []models.EtherscanNFTTx{}, nil // Return empty slice for no transactions
		}
		return nil, fmt.Errorf("Etherscan API error: %s", response.Message)
	}

	// Convert interface{} to []models.EtherscanNFTTx
	resultBytes, err := json.Marshal(response.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	var transactions []models.EtherscanNFTTx
	if err := json.Unmarshal(resultBytes, &transactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal NFT transactions: %w", err)
	}

	return transactions, nil
}

// makeRequest performs HTTP request to Etherscan API with retry logic
func (c *Client) makeRequest(params url.Values, response interface{}) error {
	requestURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	var lastErr error
	for attempt := 0; attempt < MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(RetryDelay)
		}

		resp, err := c.httpClient.Get(requestURL)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
			continue
		}

		if err := json.Unmarshal(body, response); err != nil {
			lastErr = fmt.Errorf("failed to unmarshal response: %w", err)
			continue
		}

		return nil
	}

	return fmt.Errorf("request failed after %d attempts: %w", MaxRetries, lastErr)
}
