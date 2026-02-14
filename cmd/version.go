/*
Copyright © 2025 James Evans
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long: `Print the version information for semvertool.

If the tool was built from a tagged release, it shows the version tag.
If it was built from an untagged commit (snapshot), it includes commit information.`,
	Run: runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	// If version is "dev", this is a development build (not from goreleaser)
	if version == "dev" {
		fmt.Printf("semvertool %s (snapshot)\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", date)
	} else {
		// This is a goreleaser build
		// Check if it's a snapshot build (goreleaser adds commit info for non-tag builds)
		// commit will be "none" or empty for tagged releases
		if commit != "none" && commit != "" {
			// Snapshot build - show version with commit info
			fmt.Printf("semvertool %s (snapshot)\n", version)
			fmt.Printf("  commit: %s\n", commit)
			fmt.Printf("  built:  %s\n", date)
		} else {
			// Tagged release - show clean version
			fmt.Printf("semvertool %s\n", version)
		}
	}
}
