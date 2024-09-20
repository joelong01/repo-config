package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// DeleteConfig deletes the output files (.json and .env) derived from the input JSON file.
// Parameters:
// - inputJSONFile: Path to the input JSON configuration file.
// - silent: Boolean indicating whether to run in silent mode.
// - inputReader: io.Reader for user input (useful for testing).
// Returns an error if the operation fails.
func DeleteConfig(inputJSONFile string, silent bool, inputReader io.Reader) error {
    // Determine the output file paths
    jsonOutputFile, envOutputFile, err := getOutputFilePaths(inputJSONFile)
    if err != nil {
        return err
    }

    // Collect the files to delete
    filesToDelete := []string{}

    if _, err := os.Stat(jsonOutputFile); err == nil {
        filesToDelete = append(filesToDelete, jsonOutputFile)
    }

    if _, err := os.Stat(envOutputFile); err == nil {
        filesToDelete = append(filesToDelete, envOutputFile)
    }

    if len(filesToDelete) == 0 {
        fmt.Println("No output files to delete.")
        return nil
    }

    if silent {
        // Silent mode, delete without prompting
        return deleteFiles(filesToDelete, silent)
    } else {
        // Prompt the user for confirmation
        return promptAndDeleteFiles(filesToDelete, inputReader, silent)
    }
}


// promptAndDeleteFiles prompts the user for confirmation before deleting files.
// Parameters:
// - files: Slice of file paths to delete.
// - inputReader: io.Reader for user input (useful for testing).
// - silent: Boolean indicating whether to suppress output messages.
// Returns an error if any file deletion fails or if the user cancels.
func promptAndDeleteFiles(files []string, inputReader io.Reader, silent bool) error {
    fmt.Printf("Do you want to delete the following files?\n")
    for _, file := range files {
        fmt.Printf(" - %s\n", file)
    }
    fmt.Print("Type 'yes' to confirm: ")

    reader := bufio.NewReader(inputReader)
    input, err := reader.ReadString('\n')
    if err != nil {
        return fmt.Errorf("failed to read input: %v", err)
    }
    input = strings.TrimSpace(input)

    if strings.ToLower(input) == "yes" {
        return deleteFiles(files, silent)
    } else {
        fmt.Println("Deletion canceled.")
        return nil
    }
}


// deleteFiles deletes the specified files.
// Parameters:
// - files: Slice of file paths to delete.
// - silent: if false, prints the names of the files deleted to stdout
// Returns an error if any file deletion fails.
func deleteFiles(files []string, silent bool) error {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return fmt.Errorf("failed to delete file '%s': %v", file, err)
		}
		if !silent {
			fmt.Printf("Deleted file: %s\n", file)
		}
	}
	return nil
}
