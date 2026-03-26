package utils

import (
    "encoding/json"
    "fmt"
    "os"
)

type OutputFormat string

const (
    FormatJSON  OutputFormat = "json"
    FormatTable OutputFormat = "table"
)

var CurrentFormat = FormatTable

func PrintJSON(data interface{}) {
    encoder := json.NewEncoder(os.Stdout)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(data); err != nil {
        fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
    }
}

func PrintSuccess(message string) {
    fmt.Printf("✅ %s\n", message)
}

func PrintError(message string) {
    fmt.Fprintf(os.Stderr, "❌ Error: %s\n", message)
}

func PrintInfo(message string) {
    fmt.Printf("ℹ️  %s\n", message)
}

func PrintWarning(message string) {
    fmt.Printf("⚠️  %s\n", message)
}
