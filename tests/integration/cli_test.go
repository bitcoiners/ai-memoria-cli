package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "testing"
    "time"
)

// These tests require a running Rails API server at http://localhost:3000
// Make sure to run: cd ../api && rails server

const (
    DefaultAPIURL  = "http://localhost:3000"
    DefaultEmail   = "dev@ai-memoria.com"
    DefaultPassword = "dev123"
)

var apiRunning = false

func getAPIURL() string {
    if url := os.Getenv("AI_MEMORIA_API_URL"); url != "" {
        return url
    }
    return DefaultAPIURL
}

func getTestEmail() string {
    if email := os.Getenv("AI_MEMORIA_TEST_EMAIL"); email != "" {
        return email
    }
    return DefaultEmail
}

func getTestPassword() string {
    if password := os.Getenv("AI_MEMORIA_TEST_PASSWORD"); password != "" {
        return password
    }
    return DefaultPassword
}

// getBinaryPath returns the path to the CLI binary
func getBinaryPath() string {
    // First, try to find the binary by checking common locations
    possiblePaths := []string{
        "../../bin/mem",      // When running from tests/integration/
        "./bin/mem",          // When running from project root
        "../bin/mem",         // When running from tests/ directory
        "bin/mem",            // Relative to project root
    }
    
    // Also try to get from current working directory
    if cwd, err := os.Getwd(); err == nil {
        possiblePaths = append([]string{
            filepath.Join(cwd, "bin", "mem"),
            filepath.Join(cwd, "..", "bin", "mem"),
            filepath.Join(cwd, "../..", "bin", "mem"),
        }, possiblePaths...)
    }
    
    for _, path := range possiblePaths {
        if _, err := os.Stat(path); err == nil {
            return path
        }
    }
    
    return "./bin/mem"
}

// setupTestConfig creates an isolated test environment
func setupTestConfig(t *testing.T) string {
    tmpDir := t.TempDir()
    configDir := filepath.Join(tmpDir, ".ai-memoria")
    if err := os.MkdirAll(configDir, 0700); err != nil {
        t.Fatalf("Failed to create config dir: %v", err)
    }
    
    // Set HOME to tmpDir to isolate config
    os.Setenv("HOME", tmpDir)
    os.Setenv("AI_MEMORIA_API_URL", getAPIURL())
    
    t.Cleanup(func() {
        os.Unsetenv("HOME")
        os.Unsetenv("AI_MEMORIA_API_URL")
    })
    
    return tmpDir
}

// checkAPIServer verifies the API server is running using HTTP request
func checkAPIServer(t *testing.T) {
    if apiRunning {
        return
    }
    
    url := getAPIURL() + "/up"
    resp, err := http.Get(url)
    if err != nil {
        t.Skipf("API server not running at %s. Skipping integration tests.\nError: %v", getAPIURL(), err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        t.Skipf("API server at %s returned status %d. Skipping integration tests.", getAPIURL(), resp.StatusCode)
    }
    
    apiRunning = true
    t.Logf("✅ API server is running at %s", getAPIURL())
}

// checkBinary ensures the CLI binary exists
func checkBinary(t *testing.T) {
    binaryPath := getBinaryPath()
    if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
        t.Skipf("CLI binary not found. Run 'make build' first.\nSearched in: %s", binaryPath)
    }
}

func TestIntegrationLogin(t *testing.T) {
    checkAPIServer(t)
    checkBinary(t)
    setupTestConfig(t)
    
    binaryPath := getBinaryPath()
    cmd := exec.Command(binaryPath, "auth", "login", 
        "--email", getTestEmail(), 
        "--password", getTestPassword())
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Login failed: %v\nOutput: %s", err, output)
    }
    
    // Check success message
    if !bytes.Contains(output, []byte("Logged in")) {
        t.Errorf("Expected login success message, got: %s", output)
    }
    
    // Check config file was created
    home := os.Getenv("HOME")
    configPath := filepath.Join(home, ".ai-memoria", "config.json")
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        t.Errorf("Config file not created at %s", configPath)
    }
    
    // Verify config has token
    data, err := os.ReadFile(configPath)
    if err != nil {
        t.Fatalf("Failed to read config: %v", err)
    }
    
    var config map[string]interface{}
    if err := json.Unmarshal(data, &config); err != nil {
        t.Fatalf("Failed to parse config: %v", err)
    }
    
    devConfig, ok := config["development"].(map[string]interface{})
    if !ok {
        t.Fatal("development config not found")
    }
    
    if devConfig["api_key"] == nil || devConfig["api_key"] == "" {
        t.Error("API key not saved in config")
    }
    
    if devConfig["user_id"] == nil {
        t.Error("User ID not saved in config")
    }
}

func TestIntegrationWhoami(t *testing.T) {
    checkAPIServer(t)
    checkBinary(t)
    setupTestConfig(t)
    
    binaryPath := getBinaryPath()
    
    // First login to get a token
    loginCmd := exec.Command(binaryPath, "auth", "login", 
        "--email", getTestEmail(), 
        "--password", getTestPassword())
    if _, err := loginCmd.CombinedOutput(); err != nil {
        t.Fatalf("Login failed: %v", err)
    }
    
    // Test whoami
    cmd := exec.Command(binaryPath, "auth", "whoami")
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Whoami failed: %v\nOutput: %s", err, output)
    }
    
    if !bytes.Contains(output, []byte(getTestEmail())) {
        t.Errorf("Expected email %s in output, got: %s", getTestEmail(), output)
    }
    
    // Check for user fields
    expectedFields := []string{"ID", "Email", "Username", "Name"}
    for _, field := range expectedFields {
        if !bytes.Contains(output, []byte(field)) {
            t.Errorf("Expected field '%s' in output, got: %s", field, output)
        }
    }
}

func TestIntegrationCreateUser(t *testing.T) {
    checkAPIServer(t)
    checkBinary(t)
    setupTestConfig(t)
    
    binaryPath := getBinaryPath()
    
    // Create unique email using timestamp
    uniqueEmail := "test_" + time.Now().Format("20060102150405") + "@example.com"
    username := "testuser_" + time.Now().Format("20060102150405")
    
    cmd := exec.Command(binaryPath, "users", "create",
        "--email", uniqueEmail,
        "--username", username,
        "--name", "Integration Test User",
        "--password", "password123")
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Create user failed: %v\nOutput: %s", err, output)
    }
    
    if !bytes.Contains(output, []byte("User created")) {
        t.Errorf("Expected success message, got: %s", output)
    }
    
    if !bytes.Contains(output, []byte(uniqueEmail)) {
        t.Errorf("Expected email %s in output, got: %s", uniqueEmail, output)
    }
}

func TestIntegrationCreateUserWithExistingEmail(t *testing.T) {
    checkAPIServer(t)
    checkBinary(t)
    setupTestConfig(t)
    
    binaryPath := getBinaryPath()
    
    // Try to create user with existing email
    cmd := exec.Command(binaryPath, "users", "create",
        "--email", getTestEmail(),
        "--username", "duplicate",
        "--name", "Duplicate User",
        "--password", "password123")
    
    output, err := cmd.CombinedOutput()
    
    // Should fail with error
    if err == nil {
        t.Error("Expected error when creating user with existing email, but command succeeded")
    }
    
    if !bytes.Contains(output, []byte("Error")) {
        t.Errorf("Expected error message, got: %s", output)
    }
}

func TestIntegrationStatus(t *testing.T) {
    checkAPIServer(t)
    checkBinary(t)
    setupTestConfig(t)
    
    binaryPath := getBinaryPath()
    
    cmd := exec.Command(binaryPath, "status")
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Status failed: %v\nOutput: %s", err, output)
    }
    
    if !bytes.Contains(output, []byte("API Status")) {
        t.Errorf("Expected API Status in output, got: %s", output)
    }
}

func TestIntegrationLogout(t *testing.T) {
    checkAPIServer(t)
    checkBinary(t)
    setupTestConfig(t)
    
    binaryPath := getBinaryPath()
    
    // First login
    loginCmd := exec.Command(binaryPath, "auth", "login", 
        "--email", getTestEmail(), 
        "--password", getTestPassword())
    if _, err := loginCmd.CombinedOutput(); err != nil {
        t.Fatalf("Login failed: %v", err)
    }
    
    // Verify token exists in config
    home := os.Getenv("HOME")
    configPath := filepath.Join(home, ".ai-memoria", "config.json")
    data, err := os.ReadFile(configPath)
    if err != nil {
        t.Fatalf("Failed to read config: %v", err)
    }
    
    var config map[string]interface{}
    json.Unmarshal(data, &config)
    devConfig := config["development"].(map[string]interface{})
    
    if devConfig["api_key"] == nil || devConfig["api_key"] == "" {
        t.Fatal("API key not present before logout")
    }
    
    // Then logout
    logoutCmd := exec.Command(binaryPath, "auth", "logout")
    output, err := logoutCmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Logout failed: %v\nOutput: %s", err, output)
    }
    
    if !bytes.Contains(output, []byte("Logged out")) {
        t.Errorf("Expected logout success message, got: %s", output)
    }
    
    // Verify token was removed from config
    data, err = os.ReadFile(configPath)
    if err != nil {
        t.Fatalf("Failed to read config after logout: %v", err)
    }
    
    json.Unmarshal(data, &config)
    devConfig = config["development"].(map[string]interface{})
    
    if devConfig["api_key"] != nil && devConfig["api_key"] != "" {
        t.Errorf("API key still present after logout: %v", devConfig["api_key"])
    }
}

func TestIntegrationJSONOutput(t *testing.T) {
    checkAPIServer(t)
    checkBinary(t)
    setupTestConfig(t)
    
    binaryPath := getBinaryPath()
    
    // Login first
    loginCmd := exec.Command(binaryPath, "auth", "login", 
        "--email", getTestEmail(), 
        "--password", getTestPassword())
    if _, err := loginCmd.CombinedOutput(); err != nil {
        t.Fatalf("Login failed: %v", err)
    }
    
    // Test JSON output
    cmd := exec.Command(binaryPath, "--json", "auth", "whoami")
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("JSON output failed: %v\nOutput: %s", err, output)
    }
    
    // Verify output is valid JSON
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        t.Errorf("Output is not valid JSON: %v\nOutput: %s", err, output)
    }
    
    if result["email"] != getTestEmail() {
        t.Errorf("Expected email %s, got %v", getTestEmail(), result["email"])
    }
}

func TestIntegrationInvalidLogin(t *testing.T) {
    checkAPIServer(t)
    checkBinary(t)
    setupTestConfig(t)
    
    binaryPath := getBinaryPath()
    
    cmd := exec.Command(binaryPath, "auth", "login", 
        "--email", getTestEmail(), 
        "--password", "wrongpassword")
    
    output, err := cmd.CombinedOutput()
    
    // Should fail with error
    if err == nil {
        t.Error("Expected error for invalid credentials, but command succeeded")
    }
    
    if !bytes.Contains(output, []byte("Error")) {
        t.Errorf("Expected error message, got: %s", output)
    }
}

func TestIntegrationWhoamiWithoutLogin(t *testing.T) {
    checkAPIServer(t)
    checkBinary(t)
    setupTestConfig(t)
    
    binaryPath := getBinaryPath()
    
    cmd := exec.Command(binaryPath, "auth", "whoami")
    output, err := cmd.CombinedOutput()
    
    // Should fail with error
    if err == nil {
        t.Error("Expected error when not logged in, but command succeeded")
    }
    
    if !bytes.Contains(output, []byte("Not authenticated")) {
        t.Errorf("Expected authentication error, got: %s", output)
    }
}
