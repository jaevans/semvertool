/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
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

func getTags(repo *goget.Repository) ([]*semver.Version, error) {
	iter, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	semverTags := make([]*semver.Version, 0)
	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		shortTag := ref.Name().Short()
		t, err := semver.NewVersion(shortTag)
		if err != nil {
			fmt.Printf("Could not parse tag %s as semver: %s\n", shortTag, err)
			return nil
		}
		semverTags = append(semverTags, t)
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Sort(semver.Collection(semverTags))
	return semverTags, nil
}

func getTagsStrings(repo *goget.Repository) ([]string, error) {
	tags, err := getTags(repo)
	if err != nil {
		return []string{}, err
	}
	tagStrings := make([]string, len(tags))
	for i, t := range tags {
		tagStrings[i] = t.Original()
	}
	return tagStrings, nil
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
	if len(args) != 0 {
		fmt.Println("No arguments allowed")
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
