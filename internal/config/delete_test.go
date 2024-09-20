package config

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestDeleteConfig_PromptYes(t *testing.T) {
    testDeleteConfigPrompt(t, "yes\n", true)
}

func TestDeleteConfig_PromptNo(t *testing.T) {
    testDeleteConfigPrompt(t, "no\n", false)
}

// testDeleteConfigPrompt is a helper function that tests DeleteConfig with given input and expectation.
// Parameters:
// - t: Testing object
// - userInput: Simulated user input ("yes\n" or "no\n")
// - expectDeletion: Boolean indicating whether files should be deleted
func testDeleteConfigPrompt(t *testing.T, userInput string, expectDeletion bool) {
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

    // Create mock input JSON file
    inputJSONFile := filepath.Join(dir, "test-config.json")
    if err := os.WriteFile(inputJSONFile, []byte("{}"), 0644); err != nil {
        t.Fatalf("Failed to write input JSON file: %v", err)
    }

    // Create mock output files
    jsonOutputFile, envOutputFile, err := getOutputFilePaths(inputJSONFile)
    if err != nil {
        t.Fatalf("Failed to get output file paths: %v", err)
    }

    if err := os.WriteFile(jsonOutputFile, []byte("{}"), 0644); err != nil {
        t.Fatalf("Failed to write JSON output file: %v", err)
    }

    if err := os.WriteFile(envOutputFile, []byte("KEY=value"), 0644); err != nil {
        t.Fatalf("Failed to write env output file: %v", err)
    }

    // Simulate user input
    reader := strings.NewReader(userInput)

    // Call DeleteConfig with reader
    err = DeleteConfig(inputJSONFile, false, reader)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    // Check if files are deleted based on expectation
    if expectDeletion {
        // Files should be deleted
        if _, err := os.Stat(jsonOutputFile); !os.IsNotExist(err) {
            t.Errorf("Expected JSON output file to be deleted")
        }
        if _, err := os.Stat(envOutputFile); !os.IsNotExist(err) {
            t.Errorf("Expected env output file to be deleted")
        }
    } else {
        // Files should not be deleted
        if _, err := os.Stat(jsonOutputFile); os.IsNotExist(err) {
            t.Errorf("Expected JSON output file to exist")
        }
        if _, err := os.Stat(envOutputFile); os.IsNotExist(err) {
            t.Errorf("Expected env output file to exist")
        }
    }
}
