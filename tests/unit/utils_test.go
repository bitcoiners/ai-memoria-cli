package unit

import (
    "bytes"
    "encoding/json"
    "os"
    "testing"
    
    "github.com/bitcoiners/ai-memoria-cli/internal/utils"
)

func TestPrintJSON(t *testing.T) {
    // Capture stdout
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    
    data := map[string]string{"key": "value"}
    utils.PrintJSON(data)
    
    w.Close()
    os.Stdout = oldStdout
    
    var buf bytes.Buffer
    buf.ReadFrom(r)
    output := buf.String()
    
    var result map[string]string
    if err := json.Unmarshal([]byte(output), &result); err != nil {
        t.Errorf("Failed to parse JSON output: %v", err)
    }
    
    if result["key"] != "value" {
        t.Errorf("Expected value 'value', got '%s'", result["key"])
    }
}

func TestPrintSuccess(t *testing.T) {
    // Capture stdout
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    
    utils.PrintSuccess("test message")
    
    w.Close()
    os.Stdout = oldStdout
    
    var buf bytes.Buffer
    buf.ReadFrom(r)
    output := buf.String()
    
    expected := "✅ test message\n"
    if output != expected {
        t.Errorf("Expected '%s', got '%s'", expected, output)
    }
}

func TestPrintError(t *testing.T) {
    // Capture stderr
    oldStderr := os.Stderr
    r, w, _ := os.Pipe()
    os.Stderr = w
    
    utils.PrintError("test error")
    
    w.Close()
    os.Stderr = oldStderr
    
    var buf bytes.Buffer
    buf.ReadFrom(r)
    output := buf.String()
    
    expected := "❌ Error: test error\n"
    if output != expected {
        t.Errorf("Expected '%s', got '%s'", expected, output)
    }
}

func TestPrintInfo(t *testing.T) {
    // Capture stdout
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    
    utils.PrintInfo("test info")
    
    w.Close()
    os.Stdout = oldStdout
    
    var buf bytes.Buffer
    buf.ReadFrom(r)
    output := buf.String()
    
    expected := "ℹ️  test info\n"
    if output != expected {
        t.Errorf("Expected '%s', got '%s'", expected, output)
    }
}

func TestPrintWarning(t *testing.T) {
    // Capture stdout
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    
    utils.PrintWarning("test warning")
    
    w.Close()
    os.Stdout = oldStdout
    
    var buf bytes.Buffer
    buf.ReadFrom(r)
    output := buf.String()
    
    expected := "⚠️  test warning\n"
    if output != expected {
        t.Errorf("Expected '%s', got '%s'", expected, output)
    }
}
