package unit

import (
    "testing"
    
    "github.com/bitcoiners/ai-memoria-cli/internal/models"
)

func TestUserGetID(t *testing.T) {
    tests := []struct {
        name     string
        user     models.User
        expected int
        hasError bool
    }{
        {
            name:     "int ID",
            user:     models.User{ID: 123},
            expected: 123,
            hasError: false,
        },
        {
            name:     "float64 ID",
            user:     models.User{ID: float64(456)},
            expected: 456,
            hasError: false,
        },
        {
            name:     "string ID",
            user:     models.User{ID: "789"},
            expected: 789,
            hasError: false,
        },
        {
            name:     "invalid string ID",
            user:     models.User{ID: "invalid"},
            expected: 0,
            hasError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            id, err := tt.user.GetID()
            if tt.hasError && err == nil {
                t.Errorf("Expected error but got none")
            }
            if !tt.hasError && err != nil {
                t.Errorf("Unexpected error: %v", err)
            }
            if id != tt.expected {
                t.Errorf("Expected ID %d, got %d", tt.expected, id)
            }
        })
    }
}
