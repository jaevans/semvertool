/*
Copyright © 2025 James Evans
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
)

// scriptCmd represents the script command
var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Script utilities for semantic versioning",
	Long: `Provides utilities for scripting with semantic versions.
	
These commands are designed to be used in shell scripts, returning exit codes
that can be used in conditionals.`,
}

// CompareVersions compares two semantic versions and returns:
// -1 if v1 < v2
//
//	0 if v1 = v2
//	1 if v1 > v2
//
// Returns an error if either version is invalid
func CompareVersions(v1string, v2string string) (int, error) {
	v1, err := semver.NewVersion(v1string)
	if err != nil {
		return 0, fmt.Errorf("invalid version: %s", v1string)
	}

	v2, err := semver.NewVersion(v2string)
	if err != nil {
		return 0, fmt.Errorf("invalid version: %s", v2string)
	}

	if v1.LessThan(v2) {
		return 0, nil
	} else if v1.Equal(v2) {
		return 1, nil
	} else {
		return 2, nil
	}
}

// IsReleased checks if a version is a release version (no prerelease or metadata)
// Returns true for release versions, false for prerelease versions or those with metadata
// Returns an error if the version is invalid
func IsReleased(versionString string) (bool, error) {
	v, err := semver.NewVersion(versionString)
	if err != nil {
		return false, fmt.Errorf("invalid version: %s", versionString)
	}

	return v.Prerelease() == "" && v.Metadata() == "", nil
}

// compareCmd represents the compare subcommand
var compareCmd = &cobra.Command{
	Use:   "compare <version1> <version2>",
	Short: "Compare two semantic versions",
	Long: `Compare two semantic versions and return an exit code based on the comparison:
	
 0: version1 < version2
 1: version1 = version2
 2: version1 > version2
 
 If there is an error, the command will return 3.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := CompareVersions(args[0], args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(3)
		}
		os.Exit(result)
	},
}

// releasedCmd represents the released subcommand
var releasedCmd = &cobra.Command{
	Use:   "released <version>",
	Short: "Check if a version is a release version",
	Long: `Check if a version is a release version (not a prerelease and has no metadata).
	
Returns exit code 0 if the version is a release version (X.Y.Z only),
Returns exit code 1 if the version is a prerelease or has metadata.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		isReleased, err := IsReleased(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		if isReleased {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	},
}

func init() {
	scriptCmd.AddCommand(compareCmd)
	scriptCmd.AddCommand(releasedCmd)
}
