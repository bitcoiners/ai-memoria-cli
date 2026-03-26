package unit

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/bitcoiners/ai-memoria-cli/internal/api"
    "github.com/bitcoiners/ai-memoria-cli/internal/config"
    "github.com/bitcoiners/ai-memoria-cli/internal/models"
)

func setupMockServer(t *testing.T) (*httptest.Server, *config.Config) {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        
        var req models.TokenRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        
        if req.Email == "test@example.com" && req.Password == "password" {
            resp := models.TokenResponse{
                Token: "test-token-123",
                User: models.User{
                    ID:       "123",
                    Email:    "test@example.com",
                    Username: "testuser",
                    Name:     "Test User",
                },
            }
            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(resp)
        } else {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
        }
    })
    
    mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        
        auth := r.Header.Get("Authorization")
        if auth != "Bearer test-token-123" {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
            return
        }
        
        resp := models.UserResponse{
            User: models.User{
                ID:       "123",
                Email:    "test@example.com",
                Username: "testuser",
                Name:     "Test User",
            },
        }
        json.NewEncoder(w).Encode(resp)
    })
    
    mux.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
    })
    
    server := httptest.NewServer(mux)
    
    cfg := &config.Config{
        BaseURL: server.URL,
        APIKey:  "test-token-123",
        UserID:  123,
    }
    
    return server, cfg
}

func TestGetToken(t *testing.T) {
    server, cfg := setupMockServer(t)
    defer server.Close()
    
    client := api.NewClient(cfg)
    
    resp, err := client.GetToken("test@example.com", "password")
    if err != nil {
        t.Fatalf("GetToken failed: %v", err)
    }
    
    if resp.Token != "test-token-123" {
        t.Errorf("Expected token 'test-token-123', got '%s'", resp.Token)
    }
    
    if resp.User.Email != "test@example.com" {
        t.Errorf("Expected email 'test@example.com', got '%s'", resp.User.Email)
    }
}

func TestGetTokenInvalid(t *testing.T) {
    server, cfg := setupMockServer(t)
    defer server.Close()
    
    client := api.NewClient(cfg)
    
    _, err := client.GetToken("test@example.com", "wrong")
    if err == nil {
        t.Error("Expected error for invalid credentials, got nil")
    }
}

func TestGetUser(t *testing.T) {
    server, cfg := setupMockServer(t)
    defer server.Close()
    
    client := api.NewClient(cfg)
    
    user, err := client.GetUser(123)
    if err != nil {
        t.Fatalf("GetUser failed: %v", err)
    }
    
    if user.Email != "test@example.com" {
        t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
    }
}

func TestCheckStatus(t *testing.T) {
    server, cfg := setupMockServer(t)
    defer server.Close()
    
    client := api.NewClient(cfg)
    
    status, err := client.CheckStatus()
    if err != nil {
        t.Fatalf("CheckStatus failed: %v", err)
    }
    
    if status["status"] != "ok" {
        t.Errorf("Expected status 'ok', got '%v'", status["status"])
    }
}
