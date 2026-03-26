package models

type User struct {
    ID       int    `json:"id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Name     string `json:"name"`
}

type UserResponse struct {
    User User `json:"user"`
}

type CreateUserRequest struct {
    Email                string `json:"email"`
    Password             string `json:"password"`
    PasswordConfirmation string `json:"password_confirmation"`
    Username             string `json:"username"`
    Name                 string `json:"name"`
}

type TokenRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type TokenResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}
