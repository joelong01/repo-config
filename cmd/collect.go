package cmd

import (
	"fmt"
	"os"

	"github.com/joelong01/repo-config/internal/config"
	"github.com/spf13/cobra"
)

// Variables for flags.
var (
	collectJSONFile string
	collectSilent   bool
)

// collectCmd represents the collect command.
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect repository configurations",
	Run: func(cmd *cobra.Command, args []string) {
		result := ""
		err := config.CollectConfig(collectJSONFile, collectSilent)
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
	rootCmd.AddCommand(collectCmd)

	// Define the --json flag as required.
	collectCmd.Flags().StringVarP(&collectJSONFile, "json", "j", "", "Path to the JSON configuration file (required)")
	collectCmd.MarkFlagRequired("json")

	// Define the --silent flag as optional.
	collectCmd.Flags().BoolVarP(&collectSilent, "silent", "s", false, "Run in silent mode")
}
