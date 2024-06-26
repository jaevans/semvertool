/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	// "golang.org/x/mod/semver"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump a semver version",
	Long: `Bump a semver version to the next major, minor, patch, or prerelease version.

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
	
	semvertool bump --build $(git rev-parse --short HEAD) 1.1.1-alpha.1
	1.1.1-alpha.2+3f6d1270
	
	semvertool bump 1.1.1-alpha.2+3f6d1270`,
	Run: runBump,
}

func init() {
	rootCmd.AddCommand(bumpCmd)

	addCommonBumpFlags(bumpCmd)
	bumpCmd.Flags().String("metadata", "", "Append the given string to the version as metadata.")

	viper.BindPFlags(bumpCmd.Flags())
}

func runBump(cmd *cobra.Command, args []string) {

	if len(args) < 1 {
		cmd.Help()
		return
	}
	if len(args) > 1 {
		fmt.Println("Too many arguments")
		fmt.Println()
		cmd.Help()
		return
	}
	oldV := args[0]
	newV, err := doBump(oldV, getBumpType())
	if err != nil {
		fmt.Println("Error bumping version:", err)
		return
	}

	fmt.Println(newV)
}
