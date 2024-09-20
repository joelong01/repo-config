// config_test.go

package config

import (
    "bytes"
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestLoadConfigFile(t *testing.T) {
    // Setup temporary directory
    dir, err := os.MkdirTemp("", "config_test")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(dir)

    // Valid JSON content
    validJSON := `{
        "key1": {
            "description": "Test key 1",
            "default": "value1"
        },
        "key2": {
            "description": "Test key 2",
            "default": "value2"
        }
    }`

    // Write valid JSON to a file
    validJSONPath := filepath.Join(dir, "valid.json")
    if err := os.WriteFile(validJSONPath, []byte(validJSON), 0644); err != nil {
        t.Fatalf("Failed to write valid JSON file: %v", err)
    }

    // Test with valid JSON file
    configMap, err := loadConfigFile(validJSONPath)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if len(configMap) != 2 {
        t.Errorf("Expected 2 items in configMap, got %d", len(configMap))
    }

    // Test with invalid JSON file
    invalidJSONPath := filepath.Join(dir, "invalid.json")
    if err := os.WriteFile(invalidJSONPath, []byte("{invalid json"), 0644); err != nil {
        t.Fatalf("Failed to write invalid JSON file: %v", err)
    }
    _, err = loadConfigFile(invalidJSONPath)
    if err == nil {
        t.Error("Expected error for invalid JSON, got nil")
    }

    // Test with non-existent file
    _, err = loadConfigFile(filepath.Join(dir, "nonexistent.json"))
    if err == nil {
        t.Error("Expected error for non-existent file, got nil")
    }
}

func TestSaveConfig(t *testing.T) {
    // Setup temporary directory
    dir, err := os.MkdirTemp("", "config_test")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(dir)

    // Mock home directory
    originalHome := os.Getenv("HOME")
    os.Setenv("HOME", dir)
    defer os.Setenv("HOME", originalHome)

    inputJSONFile := filepath.Join(dir, "test-config.json")
    if err := os.WriteFile(inputJSONFile, []byte("{}"), 0644); err != nil {
        t.Fatalf("Failed to write input JSON file: %v", err)
    }

    configMap := map[string]ItemConfig{
        "key1": {
            Description:                 "Test key 1",
            Default:                     "value1",
            TempEnvironmentVariableName: "KEY1",
        },
        "key2": {
            Description: "Test key 2",
            Default:     "value2",
        },
    }

    // Test saving the config
    err = saveConfig(inputJSONFile, configMap)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    // Verify JSON output file
    outputDir := filepath.Join(dir, ".repo-config")
    jsonOutputFile := filepath.Join(outputDir, ".test-config-values.json")
    if _, err := os.Stat(jsonOutputFile); os.IsNotExist(err) {
        t.Errorf("Expected JSON output file to exist")
    }

    // Verify .env output file
    envOutputFile := filepath.Join(outputDir, ".test-config-values.env")
    if _, err := os.Stat(envOutputFile); os.IsNotExist(err) {
        t.Errorf("Expected .env output file to exist")
    }

    // Read and verify contents of .env file
    envContent, err := os.ReadFile(envOutputFile)
    if err != nil {
        t.Fatalf("Failed to read .env output file: %v", err)
    }
    expectedEnvContent := "KEY1=value1"
    if strings.TrimSpace(string(envContent)) != expectedEnvContent {
        t.Errorf("Expected .env content '%s', got '%s'", expectedEnvContent, strings.TrimSpace(string(envContent)))
    }
}

func TestCompareConfigs(t *testing.T) {
    // Test case: No changes
    configMap := map[string]ItemConfig{
        "key1": {Default: "value1"},
        "key2": {Default: "value2"},
    }
    existingValues := map[string]string{
        "key1": "value1",
        "key2": "value2",
    }
    hasChanges := compareConfigs(configMap, existingValues)
    if hasChanges {
        t.Error("Expected no changes, but changes were detected")
    }

    // Test case: New setting added
    configMap["key3"] = ItemConfig{Default: "value3"}
    hasChanges = compareConfigs(configMap, existingValues)
    if !hasChanges {
        t.Error("Expected changes due to new setting, but no changes were detected")
    }

    // Test case: Setting deleted
    delete(configMap, "key1")
    hasChanges = compareConfigs(configMap, existingValues)
    if !hasChanges {
        t.Error("Expected changes due to deleted setting, but no changes were detected")
    }
}
func TestCheckForMissingValues_AllValuesPresent(t *testing.T) {
    configMap := map[string]ItemConfig{
        "key1": {Default: "value1"},
        "key2": {Default: "value2"},
    }

    missingValues := checkForMissingValues(configMap)
    if len(missingValues) != 0 {
        t.Errorf("Expected no missing values, got %v", missingValues)
    }
}
func TestCheckForMissingValues_NoExistingValues(t *testing.T) {
    configMap := map[string]ItemConfig{
        "key1": {Default: ""},
        "key2": {Default: ""},
    }
    existingValues := map[string]string{}

    updateConfigMapWithExistingValues(configMap, existingValues)

    missingValues := checkForMissingValues(configMap)
    if len(missingValues) != 2 {
        t.Errorf("Expected 2 missing values, got %v", missingValues)
    }
}

func TestCheckForMissingValues(t *testing.T) {
    // Prepare configMap with empty Default values
    configMap := map[string]ItemConfig{
        "key1": {Default: ""},
        "key2": {Default: ""},
    }

    // Existing values to update configMap with
    existingValues := map[string]string{
        "key1": "value1",
    }

    // Update configMap with existingValues
    updateConfigMapWithExistingValues(configMap, existingValues)

    // Now check for missing values
    missingValues := checkForMissingValues(configMap)
    if len(missingValues) != 1 || missingValues[0] != "key2" {
        t.Errorf("Expected missing value for 'key2', got %v", missingValues)
    }

    // Verify that configMap is updated with existing values
    if configMap["key1"].Default != "value1" {
        t.Errorf("Expected configMap['key1'].Default to be 'value1', got '%s'", configMap["key1"].Default)
    }

    // Verify that configMap["key2"].Default is still empty
    if configMap["key2"].Default != "" {
        t.Errorf("Expected configMap['key2'].Default to be '', got '%s'", configMap["key2"].Default)
    }
}


func TestInteractiveConfig(t *testing.T) {
    // Prepare configMap with one item
    configMap := map[string]ItemConfig{
        "key1": {
            Description: "Test key 1",
            Default:     "default1",
        },
    }

    // Simulate user input: Update value and then save
    userInput := bytes.NewBufferString("1\nnewvalue\ns\n")

    // Mock inputJSONFile path
    inputJSONFile := "/dev/null"

    // Run interactiveConfig
    err := interactiveConfig(configMap, inputJSONFile, userInput)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    // Verify that the value was updated
    if configMap["key1"].Default != "newvalue" {
        t.Errorf("Expected 'key1' default to be 'newvalue', got '%s'", configMap["key1"].Default)
    }
}
