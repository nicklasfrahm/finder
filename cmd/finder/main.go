package main

import (
	"os"

	"github.com/spf13/cobra"
)

var help bool

var rootCmd = &cobra.Command{
	Use:   "finder",
	Short: "CLI to organize files",
	Long: `A command line interface to sort and organize large
amounts of files and folders.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if help {
			cmd.Help()
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "display help for command")
}

// main starts the invocation of the command line interface.
func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
