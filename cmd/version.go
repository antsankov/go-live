package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VERSION of the Package, update on release.
const VERSION = "1.0.0"

var printVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of go-live",
	Long:  "Print the version of go-live",
	RunE: func(cmd *cobra.Command, args []string) error {
		printVersion()
		return nil
	},
}

func printVersion() {
	fmt.Println(VERSION)
}

func init() {
	rootCmd.AddCommand(printVersionCmd)
}
