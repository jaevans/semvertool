/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	goget "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

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
	`,
	Run: runGit,
}

func init() {
	rootCmd.AddCommand(gitCmd)

	addCommonBumpFlags(gitCmd)
	gitCmd.Flags().BoolP("hash", "s", false, "Append the short hash (sha) to the version as metadata information.")
	gitCmd.Flags().BoolP("from-commit", "c", false, "Extract the bump type from a commit message")
	viper.BindPFlags(gitCmd.Flags())
}

func getTags(repo *goget.Repository) ([]string, error) {
	iter, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	semverTags := make([]string, 0)
	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		shortTag := ref.Name().Short()
		v, err := semver.NewVersion(shortTag)
		if err != nil {
			fmt.Printf("Could not parse tag %s as semver: %s\n", shortTag, err)
			return nil
		}
		semverTags = append(semverTags, v.String())
		return nil
	}); err != nil {
		return nil, err
	}
	return semverTags, nil
}

func gitBump(version string, repo *goget.Repository) (*semver.Version, error) {
	bumpType := getBumpType()

	newVersion, err := doBump(version, bumpType)
	if err != nil {
		fmt.Println("Could not bump version", err)
		return nil, err
	}
	if viper.GetBool("hash") {

		hash, err := repo.ResolveRevision(plumbing.Revision("HEAD"))
		if err != nil {
			fmt.Println("Could not get hash", err)
			return nil, err
		}
		newV, err := newVersion.SetMetadata("sha." + hash.String()[:7])
		if err != nil {
			fmt.Println("Could not add hash metadata", err)
			return nil, err
		}
		newVersion = &newV
	}
	return newVersion, nil
}

func runGit(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Println("No arguments allowed")
		return
	}

	repo, err := goget.PlainOpenWithOptions(".", &goget.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		fmt.Println("Could not open git repository in current directory", err)
		return
	}

	semverTags, err := getTags(repo)
	if err != nil {
		fmt.Println("Could not get tags", err)
		return
	}

	if len(semverTags) == 0 {
		fmt.Println("No semver tags found")
		return
	}
	latestSemverTag := semverTags[len(semverTags)-1]

	newVersion, err := gitBump(latestSemverTag, repo)
	if err != nil {
		fmt.Println("Could not bump version", err)
		return
	}

	fmt.Println(newVersion.String())

}
