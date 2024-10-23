/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare two semver versions",
	Long: `
	Compare two semver versions and return an exit code based on their comparison.

	Examples:
	semvertool compare 1.0.0 2.0.0
	2

	semvertool compare 2.0.0 2.0.0
	1

	semvertool compare 2.0.0 1.0.0
	0
	`,
	Run: runCompare,
}

func init() {
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) {
	// Flags can only be bound once, so it needs to be done in the Run function
	// They also need to be done one at a time, so we can't use BindPFlags
	// See https://github.com/spf13/viper/issues/375#issuecomment-794668149
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(flag.Name, flag)
	})
	if len(args) != 2 {
		fmt.Println("Two arguments are required")
		_ = cmd.Help()
		os.Exit(1)
	}
	v1, err := semver.NewVersion(args[0])
	if err != nil {
		fmt.Println("Invalid first version:", err)
		os.Exit(1)
	}
	v2, err := semver.NewVersion(args[1])
	if err != nil {
		fmt.Println("Invalid second version:", err)
		os.Exit(1)
	}
	result := v1.Compare(v2)
	switch {
	case result < 0:
		os.Exit(0)
	case result == 0:
		os.Exit(1)
	case result > 0:
		os.Exit(2)
	}
}
