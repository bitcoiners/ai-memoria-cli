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
    
    // Token endpoint (POST for login, DELETE for logout)
    mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
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
        } else if r.Method == http.MethodDelete {
            auth := r.Header.Get("Authorization")
            if auth == "Bearer test-token-123" {
                w.WriteHeader(http.StatusOK)
                json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
            } else {
                w.WriteHeader(http.StatusUnauthorized)
                json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
            }
        } else {
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    })
    
    // Users endpoint
    mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            var req map[string]interface{}
            if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            
            userData, ok := req["user"].(map[string]interface{})
            if !ok {
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user data"})
                return
            }
            
            resp := models.UserResponse{
                User: models.User{
                    ID:       "456",
                    Email:    userData["email"].(string),
                    Username: userData["username"].(string),
                    Name:     userData["name"].(string),
                },
            }
            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(resp)
        }
    })
    
    // User by ID endpoint
    mux.HandleFunc("/users/123", func(w http.ResponseWriter, r *http.Request) {
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
    
    // Health check
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

func TestGetCurrentUser(t *testing.T) {
    server, cfg := setupMockServer(t)
    defer server.Close()
    
    client := api.NewClient(cfg)
    
    user, err := client.GetCurrentUser()
    if err != nil {
        t.Fatalf("GetCurrentUser failed: %v", err)
    }
    
    if user.Email != "test@example.com" {
        t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
    }
}

func TestGetCurrentUserNoID(t *testing.T) {
    cfg := &config.Config{
        BaseURL: "http://localhost:3000",
        APIKey:  "test-token-123",
        UserID:  0,
    }
    
    client := api.NewClient(cfg)
    
    _, err := client.GetCurrentUser()
    if err == nil {
        t.Error("Expected error for missing user ID, got nil")
    }
}

func TestCreateUser(t *testing.T) {
    server, cfg := setupMockServer(t)
    defer server.Close()
    
    client := api.NewClient(cfg)
    
    req := &models.CreateUserRequest{
        Email:                "newuser@example.com",
        Username:             "newuser",
        Name:                 "New User",
        Password:             "password123",
        PasswordConfirmation: "password123",
    }
    
    user, err := client.CreateUser(req)
    if err != nil {
        t.Fatalf("CreateUser failed: %v", err)
    }
    
    if user.Email != "newuser@example.com" {
        t.Errorf("Expected email 'newuser@example.com', got '%s'", user.Email)
    }
}

func TestCreateUserError(t *testing.T) {
    // Create a server that returns error on create
    mux := http.NewServeMux()
    mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "Invalid data"})
        }
    })
    server := httptest.NewServer(mux)
    defer server.Close()
    
    cfg := &config.Config{
        BaseURL: server.URL,
        APIKey:  "test-token",
        UserID:  123,
    }
    
    client := api.NewClient(cfg)
    
    req := &models.CreateUserRequest{
        Email:                "invalid@example.com",
        Username:             "invalid",
        Name:                 "Invalid User",
        Password:             "pass",
        PasswordConfirmation: "pass",
    }
    
    _, err := client.CreateUser(req)
    if err == nil {
        t.Error("Expected error for invalid data, got nil")
    }
}

func TestCreateUserWithoutAuth(t *testing.T) {
    server, _ := setupMockServer(t)
    defer server.Close()
    
    // Create client without API key
    cfgNoAuth := &config.Config{
        BaseURL: server.URL,
        APIKey:  "",
        UserID:  0,
    }
    
    client := api.NewClient(cfgNoAuth)
    
    req := &models.CreateUserRequest{
        Email:                "public@example.com",
        Username:             "public",
        Name:                 "Public User",
        Password:             "password123",
        PasswordConfirmation: "password123",
    }
    
    user, err := client.CreateUser(req)
    if err != nil {
        t.Fatalf("CreateUser without auth failed: %v", err)
    }
    
    if user.Email != "public@example.com" {
        t.Errorf("Expected email 'public@example.com', got '%s'", user.Email)
    }
}

func TestRevokeToken(t *testing.T) {
    server, cfg := setupMockServer(t)
    defer server.Close()
    
    client := api.NewClient(cfg)
    
    err := client.RevokeToken()
    if err != nil {
        t.Fatalf("RevokeToken failed: %v", err)
    }
}

func TestRevokeTokenError(t *testing.T) {
    // Create a server that returns error on revoke
    mux := http.NewServeMux()
    mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodDelete {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
        }
    })
    server := httptest.NewServer(mux)
    defer server.Close()
    
    cfg := &config.Config{
        BaseURL: server.URL,
        APIKey:  "invalid-token",
        UserID:  123,
    }
    
    client := api.NewClient(cfg)
    
    err := client.RevokeToken()
    if err == nil {
        t.Error("Expected error for invalid token, got nil")
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

func TestAPIErrorHandling(t *testing.T) {
    // Test with invalid server URL
    cfg := &config.Config{
        BaseURL: "http://localhost:9999",
        APIKey:  "test-token",
        UserID:  123,
    }
    
    client := api.NewClient(cfg)
    
    _, err := client.GetUser(123)
    if err == nil {
        t.Error("Expected error for non-existent server, got nil")
    }
}
