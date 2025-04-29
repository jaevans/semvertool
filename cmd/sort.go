package cmd

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	goget "github.com/go-git/go-git/v5"
	"github.com/jaevans/semvertool/pkg/sort"
	"github.com/spf13/cobra"
)

var (
	order        string
	noPrerelease bool
	gitTags      bool
)

// SortCmd represents the sort command
var SortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort a list of semver versions.",
	Long: `Sort a list of semver versions provided on the command line or from git tags.

Options include sorting order (ascending or descending) and filtering out prerelease versions.`,
	RunE: RunSort,
}

func init() {
	SortCmd.Flags().StringVar(&order, "order", "ascending", "Sort order: ascending, asc, descending, or dsc")
	SortCmd.Flags().BoolVar(&noPrerelease, "no-prerelease", false, "Exclude prerelease versions from the list")
	SortCmd.Flags().BoolVar(&gitTags, "git", false, "Read versions from git tags instead of command-line arguments")
}

func RunSort(cmd *cobra.Command, args []string) error {
	var versions []*semver.Version

	if gitTags {
		repo, err := goget.PlainOpenWithOptions(".", &goget.PlainOpenOptions{DetectDotGit: true})
		if err != nil {
			return err
		}
		// Fetch tags from git
		tags, err := getTags(repo)
		if err != nil {
			return err
		}
		versions = tags
	} else {
		// Parse versions from command-line arguments
		for _, arg := range args {
			v, err := semver.NewVersion(arg)
			if err != nil {
				return fmt.Errorf("invalid semver version: %s", arg)
			}
			versions = append(versions, v)
		}
	}

	if noPrerelease {
		// Filter out prerelease versions
		filtered := FilterPrerelease(versions)
		versions = filtered
	}

	ascending := true
	// Sort versions
	if order == "descending" || order == "dsc" {
		// Sort in descending order
		ascending = false
	}
	sort.SortVersions(semver.Collection(versions), ascending)

	result := VersionsToStrings(versions)
	// Print sorted versions
	fmt.Println(strings.Join(result, " "))
	return nil
}
