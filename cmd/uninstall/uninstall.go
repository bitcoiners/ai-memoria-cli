package uninstall

import (
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/bitcoiners/ai-memoria-cli/internal/utils"
)

func Handle() {
    utils.PrintWarning("This will remove the AI Memoria CLI and all its data.")
    fmt.Println()
    fmt.Println("The following will be deleted:")
    fmt.Println("  • Binary: ~/.local/bin/mem")
    fmt.Println("  • Configuration: ~/.ai-memoria/ (including API keys and tokens)")
    fmt.Println()
    
    fmt.Print("Are you sure you want to uninstall? (y/N): ")
    var response string
    fmt.Scanln(&response)
    
    if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
        utils.PrintInfo("Uninstall cancelled.")
        return
    }
    
    // Get binary path
    execPath, err := os.Executable()
    if err == nil {
        // Try to delete the binary if we can find it
        execPath, _ = filepath.EvalSymlinks(execPath)
        if execPath != "" && execPath != "/dev/null" {
            fmt.Printf("Removing binary: %s\n", execPath)
            if err := os.Remove(execPath); err != nil {
                utils.PrintWarning(fmt.Sprintf("Could not remove binary: %v", err))
                fmt.Println("You may need to manually remove it from ~/.local/bin/mem")
            }
        }
    } else {
        // Fallback to known location
        home, _ := os.UserHomeDir()
        binaryPath := filepath.Join(home, ".local", "bin", "mem")
        fmt.Printf("Removing binary: %s\n", binaryPath)
        if err := os.Remove(binaryPath); err != nil {
            utils.PrintWarning(fmt.Sprintf("Could not remove binary: %v", err))
        }
    }
    
    // Remove configuration directory
    home, _ := os.UserHomeDir()
    configDir := filepath.Join(home, ".ai-memoria")
    
    if _, err := os.Stat(configDir); err == nil {
        fmt.Printf("Removing configuration: %s\n", configDir)
        if err := os.RemoveAll(configDir); err != nil {
            utils.PrintWarning(fmt.Sprintf("Could not remove config directory: %v", err))
        }
    }
    
    fmt.Println()
    utils.PrintSuccess("AI Memoria CLI has been uninstalled.")
    fmt.Println()

}
