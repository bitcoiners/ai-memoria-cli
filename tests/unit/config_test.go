package unit

import (
    "os"
    "path/filepath"
    "testing"
    
    "github.com/bitcoiners/ai-memoria-cli/internal/config"
)

// Helper to run test with isolated config path
func withTempConfig(t *testing.T, fn func(string)) {
    tmpDir := t.TempDir()
    configPath := filepath.Join(tmpDir, "config.json")
    
    // Save original
    original := config.GetConfigPath
    config.GetConfigPath = func() string {
        return configPath
    }
    defer func() { config.GetConfigPath = original }()
    
    fn(configPath)
}

func TestLoadConfig(t *testing.T) {
    withTempConfig(t, func(configPath string) {
        cfg := config.Load("", "", "")
        if cfg.BaseURL != "http://localhost:3000" {
            t.Errorf("Expected BaseURL http://localhost:3000, got %s", cfg.BaseURL)
        }
        if cfg.Profile != "development" {
            t.Errorf("Expected Profile development, got %s", cfg.Profile)
        }
    })
}

func TestLoadConfigWithEnvVars(t *testing.T) {
    withTempConfig(t, func(configPath string) {
        os.Setenv("AI_MEMORIA_API_KEY", "test-key-123")
        os.Setenv("AI_MEMORIA_API_URL", "https://test.example.com")
        os.Setenv("AI_MEMORIA_PROFILE", "production")
        defer func() {
            os.Unsetenv("AI_MEMORIA_API_KEY")
            os.Unsetenv("AI_MEMORIA_API_URL")
            os.Unsetenv("AI_MEMORIA_PROFILE")
        }()
        
        cfg := config.Load("", "", "")
        if cfg.APIKey != "test-key-123" {
            t.Errorf("Expected APIKey test-key-123, got %s", cfg.APIKey)
        }
        if cfg.BaseURL != "https://test.example.com" {
            t.Errorf("Expected BaseURL https://test.example.com, got %s", cfg.BaseURL)
        }
        if cfg.Profile != "production" {
            t.Errorf("Expected Profile production, got %s", cfg.Profile)
        }
    })
}

func TestLoadConfigWithArgs(t *testing.T) {
    withTempConfig(t, func(configPath string) {
        cfg := config.Load("arg-key", "https://arg.example.com", "production")
        
        if cfg.APIKey != "arg-key" {
            t.Errorf("Expected APIKey arg-key, got %s", cfg.APIKey)
        }
        if cfg.BaseURL != "https://arg.example.com" {
            t.Errorf("Expected BaseURL https://arg.example.com, got %s", cfg.BaseURL)
        }
        if cfg.Profile != "production" {
            t.Errorf("Expected Profile production, got %s", cfg.Profile)
        }
    })
}

func TestSaveAndLoadConfig(t *testing.T) {
    withTempConfig(t, func(configPath string) {
        // Save config
        cfg := &config.Config{
            APIKey:  "test-save-key",
            BaseURL: "https://save.example.com",
            Profile: "development",
            UserID:  123,
        }
        
        if err := config.Save(cfg); err != nil {
            t.Fatalf("Failed to save config: %v", err)
        }
        
        // Load config
        loaded := config.Load("", "", "")
        
        if loaded.APIKey != cfg.APIKey {
            t.Errorf("Expected APIKey %s, got %s", cfg.APIKey, loaded.APIKey)
        }
        if loaded.BaseURL != cfg.BaseURL {
            t.Errorf("Expected BaseURL %s, got %s", cfg.BaseURL, loaded.BaseURL)
        }
        if loaded.UserID != cfg.UserID {
            t.Errorf("Expected UserID %d, got %d", cfg.UserID, loaded.UserID)
        }
    })
}

func TestConfigWithProfile(t *testing.T) {
    withTempConfig(t, func(configPath string) {
        // Save dev config
        devCfg := &config.Config{
            APIKey:  "dev-key",
            BaseURL: "http://localhost:3000",
            Profile: "development",
            UserID:  1,
        }
        config.Save(devCfg)
        
        // Save prod config
        prodCfg := &config.Config{
            APIKey:  "prod-key",
            BaseURL: "https://api.example.com",
            Profile: "production",
            UserID:  1,
        }
        
        // Manually set production profile in config file
        config.GetConfigPath = func() string { return configPath }
        config.Save(prodCfg)
        
        // Test development profile
        cfg := config.Load("", "", "development")
        if cfg.APIKey != "dev-key" {
            t.Errorf("Expected dev-key, got %s", cfg.APIKey)
        }
        
        // Test production profile
        cfg = config.Load("", "", "production")
        if cfg.APIKey != "prod-key" {
            t.Errorf("Expected prod-key, got %s", cfg.APIKey)
        }
    })
}
