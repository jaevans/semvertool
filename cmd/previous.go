/*
Copyright © 2025 James Evans
*/
package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/Masterminds/semver/v3"
	goget "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

var (
	releasedOnly bool
	repoPath     string
)

// previousCmd represents the previous command
var previousCmd = &cobra.Command{
	Use:   "previous",
	Short: "Get the previous semver tag from git",
	Long: `Get the previous semver tag from git history.

If the current commit is tagged with a semver version, returns the semver
version that came before it. If the current commit is not tagged, returns
the most recent semver tag in the commit's history.

The --released flag will only consider released versions (no prerelease
component or metadata).`,
	Run: runPrevious,
}

func init() {
	previousCmd.Flags().BoolVar(&releasedOnly, "released", false, "Only consider released versions (no prerelease or metadata)")
	previousCmd.Flags().StringVarP(&repoPath, "repository", "r", ".", "Path to the git repository (defaults to current directory)")
	rootCmd.AddCommand(previousCmd)
}

func getPreviousTag(repo *goget.Repository, onlyReleased bool) (string, error) {
	// Get the HEAD reference
	headRef, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("error getting HEAD: %w", err)
	}

	// Get all tags
	tagRefs, err := repo.Tags()
	if err != nil {
		return "", fmt.Errorf("error getting tags: %w", err)
	}
	defer tagRefs.Close()

	// Map to store tag name -> commit hash
	tagMap := make(map[string]plumbing.Hash)
	// Slice to store valid semver tags
	var semverTags []*semver.Version
	// Separate slice to store only released versions if needed
	var releasedVersions []*semver.Version

	// Read all tags and filter for valid semver
	err = tagRefs.ForEach(func(ref *plumbing.Reference) error {
		tagName := ref.Name().Short()

		// Try to parse as semver
		v, err := semver.NewVersion(tagName)
		if err != nil {
			// Skip if not a valid semver
			return nil
		}

		// Store the tag name and its target commit
		tagMap[v.Original()] = ref.Hash()
		semverTags = append(semverTags, v)

		// If it's a released version, add to separate slice
		if v.Prerelease() == "" && v.Metadata() == "" {
			releasedVersions = append(releasedVersions, v)
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error iterating tags: %w", err)
	}

	if len(semverTags) == 0 {
		return "", fmt.Errorf("no semver tags found")
	}

	// Choose which collection to use based on the onlyReleased flag
	tagsToUse := semverTags
	if onlyReleased {
		if len(releasedVersions) == 0 {
			return "", fmt.Errorf("no released versions found")
		}
		tagsToUse = releasedVersions
	}

	// Early check for single tag repository
	if len(tagsToUse) == 1 {
		var errorMsg string
		if onlyReleased {
			errorMsg = "no previous tag available - only one released tag exists"
		} else {
			errorMsg = "no previous tag available - only one tag exists"
		}
		return "", fmt.Errorf(errorMsg)
	}

	// Sort tags by semver (newest first)
	sort.Sort(sort.Reverse(semver.Collection(tagsToUse)))

	// Map the versions to their original strings for better error reporting
	tagsOriginal := make([]string, len(tagsToUse))
	for i, v := range tagsToUse {
		tagsOriginal[i] = v.Original()
	}

	// Check if HEAD is tagged with a valid semver from our collection
	headCommit := headRef.Hash()
	var headTagVersion *semver.Version
	var headIndex int = -1

	// Find if HEAD is tagged with a semver tag
	for i, v := range tagsToUse {
		tagCommit := tagMap[v.Original()]
		if tagCommit == headCommit {
			headTagVersion = v
			headIndex = i
			break
		}
	}

	// If HEAD is tagged with a semver version, return the previous version
	if headTagVersion != nil {
		// If HEAD is at the oldest tag, there is no previous version
		if headIndex == len(tagsToUse)-1 {
			return "", fmt.Errorf("no previous tag available - already at oldest tag")
		}
		return tagsToUse[headIndex+1].Original(), nil
	}

	// HEAD is not tagged with a semver version
	// Get the commit history to find the most recent tag
	commitIter, err := repo.Log(&goget.LogOptions{
		From: headCommit,
	})
	if err != nil {
		return "", fmt.Errorf("error getting commit history: %w", err)
	}
	defer commitIter.Close()

	// Find the most recent tag in the commit history
	var mostRecentTag *semver.Version

	err = commitIter.ForEach(func(commit *object.Commit) error {
		for _, v := range tagsToUse {
			tagCommit := tagMap[v.Original()]
			if commit.Hash == tagCommit {
				mostRecentTag = v
				return fmt.Errorf("stop") // Use an error to break out of the loop
			}
		}
		return nil
	})

	// Check if we found a tag in the history
	if mostRecentTag == nil {
		return "", fmt.Errorf("no semver tags found in commit history")
	}

	// Return the most recent tag that we found
	// This change handles the test case where we have a commit after v0.9.0
	return mostRecentTag.Original(), nil
}

func runPrevious(cmd *cobra.Command, args []string) {
	if len(args) != 0 {
		fmt.Printf("Unexpected arguments: %v\n", args)
		_ = cmd.Help()
		os.Exit(1)
	}

	// Open the git repository
	repo, err := goget.PlainOpenWithOptions(repoPath, &goget.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open git repository: %s\n", err)
		os.Exit(1)
	}

	// Get the previous semver tag
	prevTag, err := getPreviousTag(repo, releasedOnly)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding previous tag: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(prevTag)
}
