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

func HandleCreate(cfg *config.Config, args []string) {
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
        utils.PrintSuccess(fmt.Sprintf("User created: %s (ID: %v)", user.Email, user.ID))
    }
}
