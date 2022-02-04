package main

import (
	"os"

	"github.com/nicklasfrahm/finder"
	"github.com/spf13/cobra"
)

var rebuild bool

var indexCmd = &cobra.Command{
	Use:   "index <folder>",
	Short: "Build an index of all files within a folder",
	Long: `Running this command will create an index
database of all files and folders in the
specified directory.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if rebuild {
			if _, err := os.Create("finder.sqlite3"); err != nil {
				return err
			}
		}

		return finder.Index(args[0])
	},
}

func init() {
	indexCmd.Flags().BoolVarP(&rebuild, "rebuild", "r", false, "rebuild the index from scratch")

	rootCmd.AddCommand(indexCmd)
}
