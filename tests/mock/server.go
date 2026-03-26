package mock

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    
    "github.com/bitcoiners/ai-memoria-cli/internal/models"
)

type MockServer struct {
    Server *httptest.Server
    Token  string
}

func NewMockServer() *MockServer {
    token := "mock-token-12345"
    
    mux := http.NewServeMux()
    
    // Token endpoint
    mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        
        var req map[string]string
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        
        if req["email"] == "dev@ai-memoria.com" && req["password"] == "dev123" {
            resp := models.TokenResponse{
                Token: token,
                User: models.User{
                    ID:       "123",
                    Email:    "dev@ai-memoria.com",
                    Username: "developer",
                    Name:     "Developer",
                },
            }
            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(resp)
        } else {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
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
            
            userData := req["user"].(map[string]interface{})
            resp := map[string]interface{}{
                "user": map[string]interface{}{
                    "id":       "456",
                    "email":    userData["email"],
                    "username": userData["username"],
                    "name":     userData["name"],
                },
            }
            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(resp)
        }
    })
    
    // Users/me endpoint
    mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        
        auth := r.Header.Get("Authorization")
        if auth != "Bearer "+token {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
            return
        }
        
        resp := models.UserResponse{
            User: models.User{
                ID:       "123",
                Email:    "dev@ai-memoria.com",
                Username: "developer",
                Name:     "Developer",
            },
        }
        json.NewEncoder(w).Encode(resp)
    })
    
    // Health check
    mux.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
    })
    
    server := httptest.NewServer(mux)
    
    return &MockServer{
        Server: server,
        Token:  token,
    }
}

func (m *MockServer) Close() {
    m.Server.Close()
}

func (m *MockServer) URL() string {
    return m.Server.URL
}
