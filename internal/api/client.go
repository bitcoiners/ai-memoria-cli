package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
    
    "github.com/bitcoiners/ai-memoria-cli/internal/config"
    "github.com/bitcoiners/ai-memoria-cli/internal/models"
)

type Client struct {
    config *config.Config
    http   *http.Client
}

func NewClient(cfg *config.Config) *Client {
    return &Client{
        config: cfg,
        http: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (c *Client) GetToken(email, password string) (*models.TokenResponse, error) {
    reqBody := models.TokenRequest{
        Email:    email,
        Password: password,
    }
    
    jsonBody, err := json.Marshal(reqBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    url := fmt.Sprintf("%s/token", c.config.BaseURL)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    
    resp, err := c.http.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
    }
    
    var tokenResp models.TokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &tokenResp, nil
}

func (c *Client) RevokeToken() error {
    url := fmt.Sprintf("%s/token", c.config.BaseURL)
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }
    
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
    req.Header.Set("Accept", "application/json")
    
    resp, err := c.http.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
    }
    
    return nil
}

func (c *Client) GetUser(userID int) (*models.User, error) {
    url := fmt.Sprintf("%s/users/%d", c.config.BaseURL, userID)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
    req.Header.Set("Accept", "application/json")
    
    resp, err := c.http.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
    }
    
    var userResp models.UserResponse
    if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &userResp.User, nil
}

func (c *Client) GetCurrentUser() (*models.User, error) {
    if c.config.UserID == 0 {
        return nil, fmt.Errorf("user ID not found. Please log in again.")
    }
    return c.GetUser(c.config.UserID)
}

func (c *Client) CreateUser(user *models.CreateUserRequest) (*models.User, error) {
    jsonBody, err := json.Marshal(map[string]interface{}{
        "user": user,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    url := fmt.Sprintf("%s/users", c.config.BaseURL)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    
    // Add token if available (optional for public signup)
    if c.config.APIKey != "" {
        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
    }
    
    resp, err := c.http.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
    }
    
    var userResp models.UserResponse
    if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &userResp.User, nil
}

func (c *Client) CheckStatus() (map[string]interface{}, error) {
    url := fmt.Sprintf("%s/up", c.config.BaseURL)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    resp, err := c.http.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to API: %w", err)
    }
    defer resp.Body.Close()
    
    var status map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
        // If response is not JSON, just return basic status
        return map[string]interface{}{
            "status": resp.Status,
            "code":   resp.StatusCode,
        }, nil
    }
    
    return status, nil
}
