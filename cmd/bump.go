/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	// "golang.org/x/mod/semver"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump a semver version",
	Long: `
	Bump a semver version to the next major, minor, patch, or prerelease version.

	Examples:
	semvertool bump --major 0.1.0-alpha.1+build.1
	1.0.0
	
	semvertool bump --minor 1.0.0
	1.1.0
	
	semvertool bump --patch 1.1.0
	1.1.1
	
	semvertool bump --prerelease 1.1.1-alpha.0
	1.1.1-alpha.1

	semvertool bump --prerelease 1.1.1-alpha.1.0
	1.1.1-alpha.1.1
	
	semvertool bump 1.1.1-alpha.2+3f6d1270
	1.1.1

	semvertool bump --prerelease --prerelease-prefix snapshot 1.1.1
	1.1.2-snapshot.1
	`,
	Run: runBump,
}

func init() {
	bumpCmd.AddCommand(gitCmd)

	cf := getCommonBumpFlags()
	bumpCmd.Flags().AddFlagSet(cf)
	bumpCmd.Flags().StringP("metadata", "", "", "Append the given string to the version as metadata.")
	bumpCmd.MarkFlagsMutuallyExclusive("major", "minor", "patch", "prerelease", "from-message")
}

func runBump(cmd *cobra.Command, args []string) {
	// Flags can only be bound once, so it needs to be done in the Run function
	// The also need to be done one at a time, so we can't use BindPFlags
	// See https://github.com/spf13/viper/issues/375#issuecomment-794668149
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(flag.Name, flag)
	})
	if len(args) < 1 {
		_ = cmd.Help()
		return
	}
	if len(args) > 1 {
		fmt.Println("Too many arguments")
		fmt.Println()
		_ = cmd.Help()
		return
	}
	oldV := args[0]
	bumpType := getBumpType()
	newV, err := doBump(oldV, bumpType)
	if err != nil {
		fmt.Println("Error bumping version:", err)
		return
	}

	fmt.Println(newV)
}
