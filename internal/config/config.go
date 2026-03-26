package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type Config struct {
    APIKey  string `json:"api_key"`
    BaseURL string `json:"base_url"`
    Profile string `json:"profile"`
    UserID  int    `json:"user_id,omitempty"`
}

type ConfigFile struct {
    Development   *Config `json:"development,omitempty"`
    Production    *Config `json:"production,omitempty"`
    DefaultProfile string `json:"default_profile"`
}

// GetConfigPath is exported for testing
var GetConfigPath = getConfigPath

func Load(apiKey, baseURL, profile string) *Config {
    cfg := &Config{
        APIKey:  apiKey,
        BaseURL: baseURL,
        Profile: profile,
    }

    // Environment variables override
    if cfg.APIKey == "" {
        cfg.APIKey = os.Getenv("AI_MEMORIA_API_KEY")
    }
    if cfg.BaseURL == "" {
        cfg.BaseURL = os.Getenv("AI_MEMORIA_API_URL")
    }
    if cfg.Profile == "" {
        cfg.Profile = os.Getenv("AI_MEMORIA_PROFILE")
    }
    
    // Default profile
    if cfg.Profile == "" {
        cfg.Profile = "development"
    }

    // Load from config file
    configPath := getConfigPath()
    if data, err := os.ReadFile(configPath); err == nil {
        var fileConfig ConfigFile
        if err := json.Unmarshal(data, &fileConfig); err == nil {
            var profileConfig *Config
            if cfg.Profile == "development" {
                profileConfig = fileConfig.Development
            } else if cfg.Profile == "production" {
                profileConfig = fileConfig.Production
            }
            
            if profileConfig != nil {
                if cfg.APIKey == "" {
                    cfg.APIKey = profileConfig.APIKey
                }
                if cfg.BaseURL == "" {
                    cfg.BaseURL = profileConfig.BaseURL
                }
                if cfg.UserID == 0 {
                    cfg.UserID = profileConfig.UserID
                }
            }
        }
    }

    // Final defaults
    if cfg.BaseURL == "" {
        if cfg.Profile == "production" {
            cfg.BaseURL = "https://api.ai-memoria.com"
        } else {
            cfg.BaseURL = "http://localhost:3000"
        }
    }

    return cfg
}

func Save(cfg *Config) error {
    configPath := getConfigPath()
    
    // Create directory if it doesn't exist
    if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
        return err
    }
    
    // Read existing config if any
    var fileConfig ConfigFile
    if data, err := os.ReadFile(configPath); err == nil {
        json.Unmarshal(data, &fileConfig)
    }
    
    // Update the appropriate profile
    if cfg.Profile == "development" {
        fileConfig.Development = cfg
    } else if cfg.Profile == "production" {
        fileConfig.Production = cfg
    }
    
    // Save back to file
    data, err := json.MarshalIndent(fileConfig, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(configPath, data, 0600)
}

func getConfigPath() string {
    home, err := os.UserHomeDir()
    if err != nil {
        return ".ai-memoria/config.json"
    }
    return filepath.Join(home, ".ai-memoria", "config.json")
}
