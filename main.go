package main

import (
    "flag"
    "fmt"
    "os"
    
    "github.com/bitcoiners/ai-memoria-cli/cmd/auth"
    "github.com/bitcoiners/ai-memoria-cli/cmd/uninstall"
    "github.com/bitcoiners/ai-memoria-cli/cmd/users"
    "github.com/bitcoiners/ai-memoria-cli/cmd/status"
    "github.com/bitcoiners/ai-memoria-cli/internal/config"
    "github.com/bitcoiners/ai-memoria-cli/internal/utils"
)

// Version is set at build time with -ldflags
var Version = "dev"

var (
    apiKey  = flag.String("api-key", "", "API key for authentication")
    baseURL = flag.String("base-url", "", "API base URL")
    profile = flag.String("profile", "", "Configuration profile (development/production)")
    jsonOut = flag.Bool("json", false, "Output in JSON format")
    version = flag.Bool("version", false, "Show version information")
    help    = flag.Bool("help", false, "Show help")
)

func main() {
    flag.Parse()
    
    if *version {
        fmt.Printf("AI Memoria CLI version %s\n", Version)
        os.Exit(0)
    }
    
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
    
    // Handle commands
    command := flag.Arg(0)
    args := flag.Args()[1:]
    
    switch command {
    case "auth":
        handleAuth(args)
    case "users":
        handleUsers(args)
    case "status":
        handleStatus()
    case "uninstall":
        uninstall.Handle()
    default:
        fmt.Printf("Unknown command: %s\n\n", command)
        printHelp()
        os.Exit(1)
    }
}

func handleAuth(args []string) {
    if len(args) < 1 {
        fmt.Println("Usage: mem auth login|logout|whoami [options]")
        os.Exit(1)
    }
    
    subcommand := args[0]
    subArgs := args[1:]
    
    // Load configuration
    cfg := config.Load(*apiKey, *baseURL, *profile)
    
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

func handleUsers(args []string) {
    if len(args) < 1 {
        fmt.Println("Usage: mem users create [options]")
        os.Exit(1)
    }
    
    subcommand := args[0]
    subArgs := args[1:]
    
    // Load configuration
    cfg := config.Load(*apiKey, *baseURL, *profile)
    
    switch subcommand {
    case "create":
        users.HandleCreate(cfg, subArgs)
    default:
        fmt.Printf("Unknown users command: %s\n", subcommand)
        os.Exit(1)
    }
}

func handleStatus() {
    // Load configuration
    cfg := config.Load(*apiKey, *baseURL, *profile)
    status.Handle(cfg)
}

func printHelp() {
    fmt.Printf(`AI Memoria CLI version %s - Command line interface for AI Memoria

Usage:
  mem [global options] <command> [arguments] [options]

Global Options:
  --version        Show version information
  --help           Show this help message
  --api-key KEY    API key for authentication (or AI_MEMORIA_API_KEY env)
  --base-url URL   API base URL (or AI_MEMORIA_API_URL env, default: http://localhost:3000)
  --profile NAME   Configuration profile (development/production)
  --json           Output in JSON format

Commands:
  auth             Authentication commands
    auth login --email EMAIL --password PASS  Login and save token
    auth logout                                Logout and remove token
    auth whoami                                Show current user info
    
  users            Manage users
    users create --email EMAIL --username USER --name NAME --password PASS  Create a new user
    
  status           Check API connectivity and health
  
  uninstall        Remove AI Memoria CLI from your system

Examples:
  mem --version
  mem auth login --email dev@ai-memoria.com --password dev123
  mem auth whoami
  mem users create --email new@example.com --username newuser --name "New User" --password pass123
  mem status
  mem --profile production auth whoami
  mem uninstall`, Version)
}
