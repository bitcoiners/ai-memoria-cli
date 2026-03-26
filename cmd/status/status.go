package status

import (
    "fmt"
    "os"
    
    "github.com/bitcoiners/ai-memoria-cli/internal/api"
    "github.com/bitcoiners/ai-memoria-cli/internal/config"
    "github.com/bitcoiners/ai-memoria-cli/internal/utils"
)

func Handle(cfg *config.Config, args []string) {
    client := api.NewClient(cfg)
    
    status, err := client.CheckStatus()
    if err != nil {
        utils.PrintError(fmt.Sprintf("Failed to connect: %v", err))
        os.Exit(1)
    }
    
    if utils.CurrentFormat == utils.FormatJSON {
        utils.PrintJSON(status)
    } else {
        fmt.Println("✅ API Status:")
        for key, value := range status {
            fmt.Printf("  %s: %v\n", key, value)
        }
    }
}
