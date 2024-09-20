package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

// ItemConfig represents the structure of each configuration item.
type ItemConfig struct {
	Description                 string `json:"description"`
	ShellScript                 string `json:"shellscript"`
	Default                     string `json:"default"`
	TempEnvironmentVariableName string `json:"tempEnvironmentVariableName"`
	RequiredAsEnv               bool   `json:"requiredAsEnv"`
	// Add any additional fields if necessary.
}

// CollectConfig loads and processes the configuration based on the input JSON file and silent mode.
func CollectConfig(inputJSONFile string, silent bool) error {
	// Check if the input JSON file exists.
	if _, err := os.Stat(inputJSONFile); os.IsNotExist(err) {
		return fmt.Errorf("JSON file '%s' not found", inputJSONFile)
	}

	// Load the input JSON file.
	configMap, err := loadConfigFile(inputJSONFile)
	if err != nil {
		return err
	}

	// Determine the output file paths.
	jsonOutputFile, _, err := GetOutputFilePaths(inputJSONFile)
	if err != nil {
		return err
	}

	// Load existing values from the output JSON file if it exists.
	existingValues := make(map[string]string)
	outputFileExists := false
	if _, err := os.Stat(jsonOutputFile); err == nil {
		outputFileExists = true
		existingValues, err = loadExistingValues(jsonOutputFile)
		if err != nil {
			return err
		}
	}

	// Update configMap with existing values
	updateConfigMapWithExistingValues(configMap, existingValues)

	// Handle silent mode.
	if silent {
		// Get modification times.
		inputJSONModTime, err := getFileModTime(inputJSONFile)
		if err != nil {
			return err
		}

		outputJSONModTime, err := getFileModTime(jsonOutputFile)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		// Check if input JSON is newer than output JSON or output doesn't exist.
		if !outputFileExists || inputJSONModTime.After(outputJSONModTime) {
			// Check for new or deleted settings.
			hasChanges := compareConfigs(configMap, existingValues)
			if hasChanges {
				// New settings found or settings deleted, proceed interactively.
				fmt.Fprintln(os.Stderr, "Configuration changes detected. Proceeding interactively.")
				return interactiveConfig(configMap, inputJSONFile, os.Stdin)
			} else {
				// No changes, save updated config silently.
				return saveConfig(inputJSONFile, configMap)
			}
		} else {
			// Input JSON is older or same age.
			// Verify all settings have values.
			missingValues := checkForMissingValues(configMap)
			if len(missingValues) > 0 {
				// Missing values found, proceed interactively.
				fmt.Fprintln(os.Stderr, "Missing values detected. Proceeding interactively.")
				return interactiveConfig(configMap, inputJSONFile, os.Stdin)
			} else {
				// All values present, return without action.
				return nil
			}
		}
	}

	// Not in silent mode or silent mode overridden, proceed interactively.
	return interactiveConfig(configMap, inputJSONFile, os.Stdin)
}

// updateConfigMapWithExistingValues updates configMap with values from existingValues.
func updateConfigMapWithExistingValues(configMap map[string]ItemConfig, existingValues map[string]string) {
	for key, item := range configMap {
		if value, exists := existingValues[key]; exists {
			item.Default = value
			configMap[key] = item
		}
	}
}

// loadConfigFile loads and parses the JSON configuration file into a map.
func loadConfigFile(jsonFile string) (map[string]ItemConfig, error) {
	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configMap := make(map[string]ItemConfig)
	if err := decoder.Decode(&configMap); err != nil {
		return nil, fmt.Errorf("failed to parse JSON file: %v", err)
	}
	return configMap, nil
}


// GetOutputFilePaths determines the output file paths based on the input JSON file.
// Now derives the project name from the Git repository.
func GetOutputFilePaths(inputJSONFile string) (string, string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("failed to get home directory: %v", err)
	}

	// Get the absolute path of the input JSON file
	absInputJSONFile, err := filepath.Abs(inputJSONFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to get absolute path of input JSON file: %v", err)
	}

	// Determine the project name from the Git repository
	projectName, err := getProjectName(absInputJSONFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to get Git project name: %v", err)
	}

	// Build the output directory path including the project name
	outputDir := filepath.Join(homeDir, ".repo-config", projectName)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return "", "", fmt.Errorf("failed to create directory '%s': %v", outputDir, err)
		}
	}

	baseFilename := filepath.Base(inputJSONFile)
	ext := filepath.Ext(baseFilename)
	nameWithoutExt := strings.TrimSuffix(baseFilename, ext)
	jsonOutputFile := filepath.Join(outputDir, fmt.Sprintf(".%s-values.json", nameWithoutExt))
	envOutputFile := filepath.Join(outputDir, fmt.Sprintf(".%s-values.env", nameWithoutExt))

	return jsonOutputFile, envOutputFile, nil
}

// getProjectName determines the project name based on the Git repository.
// if it is not a git respository, it returns the current directofy as
// the project name.
// TODO: have the container set an environment variable that can be used
//
//	as the project name and also add a --project-name as an optional
//	flag that will override this.
//
// It returns the base name of the Git repository's root directory.
func getProjectName(startPath string) (string, error) {
	// Prepare the command to get the top-level directory of the Git repository
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = filepath.Dir(startPath)
	gitRootPath := ""
	// Execute the command and capture the output
	output, err := cmd.Output()

	if err != nil {
		gitRootPath = os.Getenv("PWD")

	} else {
		// Get the path to the top-level directory
		gitRootPath = strings.TrimSpace(string(output))
	}

	// Extract the base name as the project name

	projectName := filepath.Base(gitRootPath)
	return projectName, nil
}

// loadExistingValues loads existing values from the output JSON file.
func loadExistingValues(jsonOutputFile string) (map[string]string, error) {
	file, err := os.Open(jsonOutputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open existing values file: %v", err)
	}
	defer file.Close()

	var nameValuePairs []map[string]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&nameValuePairs); err != nil {
		return nil, fmt.Errorf("failed to parse existing values file: %v", err)
	}

	existingValues := make(map[string]string)
	for _, pair := range nameValuePairs {
		existingValues[pair["Name"]] = pair["Value"]
	}
	return existingValues, nil
}

// getFileModTime gets the modification time of a file.
func getFileModTime(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}

// compareConfigs checks for new or deleted settings between the input config and existing values.
func compareConfigs(configMap map[string]ItemConfig, existingValues map[string]string) bool {
	hasChanges := false

	// Check for new settings.
	for key := range configMap {
		if _, exists := existingValues[key]; !exists {
			hasChanges = true
			break
		}
	}

	// Check for deleted settings.
	for key := range existingValues {
		if _, exists := configMap[key]; !exists {
			hasChanges = true
			break
		}
	}

	return hasChanges
}

// checkForMissingValues checks if any settings are missing values in configMap.
func checkForMissingValues(configMap map[string]ItemConfig) []string {
	missingValues := []string{}
	for key, item := range configMap {
		if item.Default == "" {
			missingValues = append(missingValues, key)
		}
	}
	return missingValues
}

// interactiveConfig handles the interactive prompt for updating settings.
func interactiveConfig(configMap map[string]ItemConfig, inputJSONFile string, inputReader io.Reader) error {
	keys := make([]string, 0, len(configMap))
	for key := range configMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	reader := bufio.NewReader(inputReader)

	for {
		// Use tabwriter for formatting.
		writer := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
		fmt.Fprintln(writer, "Index\tDescription\tDefault Value")
		fmt.Fprintln(writer, "-----\t-----------\t-------------")
		for i, key := range keys {
			item := configMap[key]
			fmt.Fprintf(writer, "%d\t%s\t%s\n", i+1, item.Description, item.Default)
		}
		writer.Flush()

		n := len(keys)
		fmt.Fprintf(os.Stderr, "\nEnter a number [1..%d] to update. Press 's' to save, or 'c' to continue: ", n)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %v", err)
		}
		input = strings.TrimSpace(input)

		switch strings.ToLower(input) {
		case "s":
			// Save the updated configuration.
			if err := saveConfig(inputJSONFile, configMap); err != nil {
				return fmt.Errorf("failed to save configuration: %v", err)
			}
			fmt.Fprintln(os.Stderr, "Configuration saved.")
		case "c":
			return nil
		default:
			// Check if input is a valid number.
			index, err := strconv.Atoi(input)
			if err != nil || index < 1 || index > n {
				fmt.Fprintf(os.Stderr, "Invalid input. Please enter a number between 1 and %d, 's', 'x', or 'c'.\n", n)
				continue
			}
			// Update the selected item.
			key := keys[index-1]
			item := configMap[key]
			fmt.Fprint(os.Stderr, "Enter new default value (leave empty to keep current value): ")
			newValue, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %v", err)
			}
			newValue = strings.TrimSpace(newValue)
			if newValue != "" {
				item.Default = newValue
				configMap[key] = item
				fmt.Fprintf(os.Stderr, "Default value for '%s' updated to '%s'\n", item.Description, newValue)
			} else {
				fmt.Fprintf(os.Stderr, "Default value for '%s' remains '%s'\n", item.Description, item.Default)
			}
		}
	}
}

// saveConfig saves the updated configuration to the appropriate files.
func saveConfig(inputJSONFile string, configMap map[string]ItemConfig) error {
	// Use getOutputFilePaths to determine the output file paths
	jsonOutputFile, envOutputFile, err := GetOutputFilePaths(inputJSONFile)
	if err != nil {
		return err
	}

	// Prepare data for .json file (array of name-value pairs).
	var nameValuePairs []map[string]string
	for key, item := range configMap {
		nameValuePair := map[string]string{
			"Name":  key,
			"Value": item.Default,
		}
		nameValuePairs = append(nameValuePairs, nameValuePair)
	}

	// Write to the .json file.
	jsonFile, err := os.Create(jsonOutputFile)
	if err != nil {
		return fmt.Errorf("failed to create JSON output file: %v", err)
	}
	defer jsonFile.Close()

	jsonEncoder := json.NewEncoder(jsonFile)
	jsonEncoder.SetIndent("", "    ") // Format JSON with indentation.
	if err := jsonEncoder.Encode(nameValuePairs); err != nil {
		return fmt.Errorf("failed to write JSON output file: %v", err)
	}

	// Prepare data for .env file.
	var envLines []string
	for _, item := range configMap {
		if item.TempEnvironmentVariableName != "" {
			line := fmt.Sprintf("%s=%s", item.TempEnvironmentVariableName, item.Default)
			envLines = append(envLines, line)
		}
	}

	// Write to the .env file if there are any environment variables.
	if len(envLines) > 0 {
		envFile, err := os.Create(envOutputFile)
		if err != nil {
			return fmt.Errorf("failed to create env output file: %v", err)
		}
		defer envFile.Close()

		envContent := strings.Join(envLines, "\n")
		if _, err := envFile.WriteString(envContent); err != nil {
			return fmt.Errorf("failed to write env output file: %v", err)
		}
	}

	return nil
}
