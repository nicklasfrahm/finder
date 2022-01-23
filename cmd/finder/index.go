package main

import (
	"github.com/nicklasfrahm/finder"
	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index <folder>",
	Short: "Build an index of all files within a folder",
	Long: `Running this command will create an index
database of all files and folders in the
specified directory.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return finder.Index(args[0])
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
