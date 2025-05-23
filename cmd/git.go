/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	goget "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var ErrNoSemverTags = fmt.Errorf("no semver tags found")

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Bump a version based on git tags",
	Long: `
	Bump a version based on the latest semver tag in the git repository.

	Examples:
	git tag v0.1.0
	semvertool git 
	v0.1.1

	git tag v1.0.0
	semvertool git --minor
	v1.1.0

	git tag v1.0.0-alpha.0
	semvertool git --prerelease
	v1.0.0-alpha.1

	git tag v1.0.0-alpha.1
	semvertool git --prerelease --hash
	v1.0.0-alpha.2+3f6d1270

	git tag v1.0.0-alpha.2+3f6d1270
	semvertool git --minor
	v1.1.0

	git tag v1.0.0-alpha.3+3f6d1270
	semvertool git
	v1.0.1

	git tag v1.1.0
	semvertool git --prerelease
	v1.1.1-alpha.1
	`,
	Run: runGit,
}

var deprecatedGitCmd = &cobra.Command{
	Use:        gitCmd.Use,
	Short:      gitCmd.Short,
	Long:       gitCmd.Long,
	Run:        gitCmd.Run,
	Deprecated: "and will be removed in a future release. Use 'bump git' instead",
}

func init() {
	cf := getCommonBumpFlags()
	gitCmd.Flags().AddFlagSet(cf)
	gitCmd.Flags().BoolP("hash", "s", false, "Append the short hash (sha) to the version as metadata information.")
	gitCmd.Flags().BoolP("from-commit", "c", false, "Extract the bump type from a commit message")
	gitCmd.MarkFlagsMutuallyExclusive("major", "minor", "patch", "prerelease", "from-message", "from-commit")

	deprecatedGitCmd.Flags().AddFlagSet(cf)
	deprecatedGitCmd.Flags().BoolP("hash", "s", false, "Append the short hash (sha) to the version as metadata information.")
	deprecatedGitCmd.Flags().BoolP("from-commit", "c", false, "Extract the bump type from a commit message")
	deprecatedGitCmd.MarkFlagsMutuallyExclusive("major", "minor", "patch", "prerelease", "from-message", "from-commit")

}

func gitBump(repo *goget.Repository) (*semver.Version, error) {
	bumpType := getBumpType()

	semverTags, err := getTags(repo)
	if err != nil {
		fmt.Println("Could not get tags", err)
		return nil, err
	}

	if len(semverTags) == 0 {
		fmt.Println("No semver tags found")
		return nil, ErrNoSemverTags
	}
	latestSemverTag := semverTags[len(semverTags)-1]

	newVersion, err := doBump(latestSemverTag.Original(), bumpType)
	if err != nil {
		fmt.Println("Could not bump version", err)
		return nil, err
	}
	if viper.GetBool("hash") {

		hash, err := repo.ResolveRevision(plumbing.Revision("HEAD"))
		if err != nil {
			fmt.Println("Could not get revision hash of HEAD", err)
			return nil, err
		}
		newV, err := newVersion.SetMetadata(hash.String()[:7])
		if err != nil {
			fmt.Println("Could not add hash metadata", err)
			return nil, err
		}
		newVersion = &newV
	}
	return newVersion, nil
}

func runGit(cmd *cobra.Command, args []string) {
	// Flags can only be bound once, so it needs to be done in the Run function
	// The also need to be done one at a time, so we can't use BindPFlags
	// See https://github.com/spf13/viper/issues/375#issuecomment-794668149
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(flag.Name, flag)
	})
	if len(args) != 0 {
		fmt.Printf("Unexpected arguments: %v\n", args)
		return
	}

	repo, err := goget.PlainOpenWithOptions(".", &goget.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		fmt.Println("Could not open git repository in current directory", err)
		return
	}

	newVersion, err := gitBump(repo)
	if err != nil {
		fmt.Println("Could not bump version", err)
		return
	}

	fmt.Println(newVersion.String())

}
