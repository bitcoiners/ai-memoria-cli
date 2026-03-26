package auth

import (
    "fmt"
    "os"
    
    "github.com/bitcoiners/ai-memoria-cli/internal/api"
    "github.com/bitcoiners/ai-memoria-cli/internal/config"
    "github.com/bitcoiners/ai-memoria-cli/internal/utils"
)

func HandleLogin(cfg *config.Config, email, password string) {
    if email == "" || password == "" {
        utils.PrintError("Email and password are required")
        fmt.Println("Usage: ai-memoria-cli auth login --email <email> --password <password>")
        os.Exit(1)
    }
    
    client := api.NewClient(cfg)
    
    utils.PrintInfo(fmt.Sprintf("Authenticating as %s...", email))
    
    tokenResp, err := client.GetToken(email, password)
    if err != nil {
        utils.PrintError(err.Error())
        os.Exit(1)
    }
    
    // Convert ID to int
    userID, err := tokenResp.User.GetID()
    if err != nil {
        utils.PrintError(fmt.Sprintf("Failed to parse user ID: %v", err))
        os.Exit(1)
    }
    
    // Save token and user ID to config
    cfg.APIKey = tokenResp.Token
    cfg.UserID = userID
    if err := config.Save(cfg); err != nil {
        utils.PrintError(fmt.Sprintf("Failed to save config: %v", err))
        os.Exit(1)
    }
    
    if utils.CurrentFormat == utils.FormatJSON {
        utils.PrintJSON(map[string]interface{}{
            "success": true,
            "user":    tokenResp.User,
            "token":   tokenResp.Token,
        })
    } else {
        utils.PrintSuccess(fmt.Sprintf("Logged in as %s (%s)", tokenResp.User.Name, tokenResp.User.Email))
        utils.PrintInfo(fmt.Sprintf("Token saved to ~/.ai-memoria/config.json"))
    }
}

func HandleLogout(cfg *config.Config) {
    if cfg.APIKey == "" {
        utils.PrintInfo("Not logged in")
        return
    }
    
    client := api.NewClient(cfg)
    
    if err := client.RevokeToken(); err != nil {
        utils.PrintWarning(fmt.Sprintf("API error: %v", err))
    }
    
    // Clear token and user ID from config
    cfg.APIKey = ""
    cfg.UserID = 0
    if err := config.Save(cfg); err != nil {
        utils.PrintError(fmt.Sprintf("Failed to save config: %v", err))
        os.Exit(1)
    }
    
    if utils.CurrentFormat == utils.FormatJSON {
        utils.PrintJSON(map[string]interface{}{
            "success": true,
            "message": "Logged out successfully",
        })
    } else {
        utils.PrintSuccess("Logged out successfully")
    }
}

func HandleWhoami(cfg *config.Config) {
    if cfg.APIKey == "" {
        utils.PrintError("Not authenticated. Please run: ai-memoria-cli auth login")
        os.Exit(1)
    }
    
    client := api.NewClient(cfg)
    
    user, err := client.GetCurrentUser()
    if err != nil {
        utils.PrintError(err.Error())
        os.Exit(1)
    }
    
    if utils.CurrentFormat == utils.FormatJSON {
        utils.PrintJSON(user)
    } else {
        fmt.Println("\nCurrent User:")
        fmt.Printf("  ID:       %v\n", user.ID)
        fmt.Printf("  Email:    %s\n", user.Email)
        fmt.Printf("  Username: %s\n", user.Username)
        fmt.Printf("  Name:     %s\n", user.Name)
        fmt.Println()
    }
}
