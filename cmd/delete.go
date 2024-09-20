package cmd

import (
    "fmt"
    "os"

    "github.com/joelong01/repo-config/internal/config"
    "github.com/spf13/cobra"
)

var (
    deleteJSONFile string
    deleteSilent   bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
    Use:   "delete",
    Short: "Delete repository configuration output files",
    Run: func(cmd *cobra.Command, args []string) {
        result := ""
        err := config.DeleteConfig(deleteJSONFile, deleteSilent, os.Stdin)
        if err == nil {
			jsonOutputFile, envOutputFile, _ := config.GetOutputFilePaths(collectJSONFile)
			result = config.CreateSuccessOutput(jsonOutputFile, envOutputFile)

		} else {
			result = config.CreateErrorOutput(err)
		}

		fmt.Print(result)
		os.Exit(1)
       
    },
}

func init() {
    rootCmd.AddCommand(deleteCmd)

    // Define the --json flag as required
    deleteCmd.Flags().StringVarP(&deleteJSONFile, "json", "j", "", "Path to the JSON configuration file (required)")
    deleteCmd.MarkFlagRequired("json")

    // Define the --silent flag as optional
    deleteCmd.Flags().BoolVarP(&deleteSilent, "silent", "s", false, "Run in silent mode")
}
