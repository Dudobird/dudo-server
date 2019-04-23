package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version for current server
const Version = "v0.1.0"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number of dudo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("current version is: %s -- HEAD\n", Version)
	},
}
