/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "semvertool",
	Short: "Manipulate semver version strings.",
	Long:  `semvertool is a CLI tool to manage semver versions. It can bump major, minor, patch, and prerelease versions.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: rootCmd.Run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bump.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.AddCommand(deprecatedGitCmd)
	rootCmd.AddCommand(bumpCmd)

	// Add the sort subcommand to the root command
	rootCmd.AddCommand(SortCmd)

	rootCmd.AddCommand(scriptCmd)
}
