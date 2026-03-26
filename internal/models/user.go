package models

import "fmt"

type User struct {
    ID       interface{} `json:"id"` // Can be string or int
    Email    string      `json:"email"`
    Username string      `json:"username"`
    Name     string      `json:"name"`
}

// GetID returns the user ID as an int
func (u *User) GetID() (int, error) {
    switch v := u.ID.(type) {
    case float64:
        return int(v), nil
    case int:
        return v, nil
    case string:
        // Try to parse string to int
        var id int
        _, err := fmt.Sscan(v, &id)
        return id, err
    default:
        return 0, fmt.Errorf("unexpected ID type: %T", u.ID)
    }
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
