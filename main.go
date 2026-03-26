package main

import (
    "flag"
    "fmt"
    "os"
    
    "github.com/bitcoiners/ai-memoria-cli/cmd/auth"
    "github.com/bitcoiners/ai-memoria-cli/cmd/users"
    "github.com/bitcoiners/ai-memoria-cli/cmd/status"
    "github.com/bitcoiners/ai-memoria-cli/internal/config"
    "github.com/bitcoiners/ai-memoria-cli/internal/utils"
)

var (
    apiKey  = flag.String("api-key", "", "API key for authentication")
    baseURL = flag.String("base-url", "", "API base URL")
    profile = flag.String("profile", "", "Configuration profile (development/production)")
    jsonOut = flag.Bool("json", false, "Output in JSON format")
    help    = flag.Bool("help", false, "Show help")
)

func main() {
    flag.Parse()
    
    if *jsonOut {
        utils.CurrentFormat = utils.FormatJSON
    }
    
    if *help {
        printHelp()
        os.Exit(0)
    }
    
    if flag.NArg() == 0 {
        printHelp()
        os.Exit(1)
    }
    
    // Load configuration
    cfg := config.Load(*apiKey, *baseURL, *profile)
    
    // Handle commands
    command := flag.Arg(0)
    args := flag.Args()[1:]
    
    switch command {
    case "auth":
        handleAuth(cfg, args)
    case "users":
        users.Handle(cfg, args)
    case "status":
        status.Handle(cfg, args)
    default:
        fmt.Printf("Unknown command: %s\n\n", command)
        printHelp()
        os.Exit(1)
    }
}

func handleAuth(cfg *config.Config, args []string) {
    if len(args) < 1 {
        fmt.Println("Usage: ai-memoria-cli auth login|logout|whoami [options]")
        os.Exit(1)
    }
    
    subcommand := args[0]
    subArgs := args[1:]
    
    switch subcommand {
    case "login":
        email := ""
        password := ""
        
        fs := flag.NewFlagSet("login", flag.ExitOnError)
        fs.StringVar(&email, "email", "", "User email")
        fs.StringVar(&password, "password", "", "User password")
        fs.Parse(subArgs)
        
        auth.HandleLogin(cfg, email, password)
        
    case "logout":
        auth.HandleLogout(cfg)
        
    case "whoami":
        auth.HandleWhoami(cfg)
        
    default:
        fmt.Printf("Unknown auth command: %s\n", subcommand)
        os.Exit(1)
    }
}

func printHelp() {
    fmt.Println(`AI Memoria CLI - Command line interface for AI Memoria

Usage:
  ai-memoria-cli [global options] <command> [arguments] [options]

Global Options:
  --api-key KEY    API key for authentication (or AI_MEMORIA_API_KEY env)
  --base-url URL   API base URL (or AI_MEMORIA_API_URL env, default: http://localhost:3000)
  --profile NAME   Configuration profile (development/production)
  --json           Output in JSON format
  --help           Show this help message

Commands:
  auth             Authentication commands
    auth login --email EMAIL --password PASS  Login and save token
    auth logout                                Logout and remove token
    auth whoami                                Show current user info
    
  users            Manage users
    users create --email EMAIL --username USER --name NAME --password PASS  Create a new user
    
  status           Check API connectivity and health

Examples:
  ai-memoria-cli auth login --email dev@ai-memoria.com --password dev123
  ai-memoria-cli auth whoami
  ai-memoria-cli users create --email new@example.com --username newuser --name "New User" --password pass123
  ai-memoria-cli status
  ai-memoria-cli --profile production auth whoami`)
}
