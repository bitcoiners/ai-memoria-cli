package users

import (
    "flag"
    "fmt"
    "os"
    
    "github.com/bitcoiners/ai-memoria-cli/internal/api"
    "github.com/bitcoiners/ai-memoria-cli/internal/config"
    "github.com/bitcoiners/ai-memoria-cli/internal/models"
    "github.com/bitcoiners/ai-memoria-cli/internal/utils"
)

func Handle(cfg *config.Config, args []string) {
    if len(args) < 1 {
        fmt.Println("Usage: ai-memoria-cli users create [options]")
        os.Exit(1)
    }
    
    subcommand := args[0]
    subArgs := args[1:]
    
    switch subcommand {
    case "create":
        createUser(cfg, subArgs)
    default:
        fmt.Printf("Unknown users command: %s\n", subcommand)
        os.Exit(1)
    }
}

func createUser(cfg *config.Config, args []string) {
    fs := flag.NewFlagSet("create", flag.ExitOnError)
    email := fs.String("email", "", "User email address")
    username := fs.String("username", "", "Username")
    name := fs.String("name", "", "Full name")
    password := fs.String("password", "", "Password")
    fs.Parse(args)
    
    if *email == "" || *username == "" || *name == "" || *password == "" {
        utils.PrintError("All fields are required: --email, --username, --name, --password")
        fs.Usage()
        os.Exit(1)
    }
    
    client := api.NewClient(cfg)
    
    req := &models.CreateUserRequest{
        Email:                *email,
        Username:             *username,
        Name:                 *name,
        Password:             *password,
        PasswordConfirmation: *password,
    }
    
    user, err := client.CreateUser(req)
    if err != nil {
        utils.PrintError(err.Error())
        os.Exit(1)
    }
    
    if utils.CurrentFormat == utils.FormatJSON {
        utils.PrintJSON(user)
    } else {
        utils.PrintSuccess(fmt.Sprintf("User created: %s (ID: %d)", user.Email, user.ID))
    }
}
