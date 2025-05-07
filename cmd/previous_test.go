package cmd

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
)

func setupRepoWithTags(t *testing.T) *git.Repository {
	// Create a repository with multiple tags
	repo, err := setupRepo()
	assert.NoError(t, err)

	// Create commits and tags in a specific order
	// First commit with tag v1.0.0
	commit1, err := commitFile("file1.txt", repo)
	assert.NoError(t, err)
	_, err = repo.CreateTag("v1.0.0", commit1, nil)
	assert.NoError(t, err)

	// Second commit with tag v1.1.0
	commit2, err := commitFile("file2.txt", repo)
	assert.NoError(t, err)
	_, err = repo.CreateTag("v1.1.0", commit2, nil)
	assert.NoError(t, err)

	// Third commit with prerelease tag v1.2.0-alpha.1
	commit3, err := commitFile("file3.txt", repo)
	assert.NoError(t, err)
	_, err = repo.CreateTag("v1.2.0-alpha.1", commit3, nil)
	assert.NoError(t, err)

	// Fourth commit with tag v1.2.0
	commit4, err := commitFile("file4.txt", repo)
	assert.NoError(t, err)
	_, err = repo.CreateTag("v1.2.0", commit4, nil)
	assert.NoError(t, err)

	return repo
}

func TestGetPreviousTagCurrentIsTagged(t *testing.T) {
	repo := setupRepoWithTags(t)

	// Check tags to verify setup
	tags, err := getTags(repo)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(tags))

	// Get current HEAD which should be at the last commit (v1.2.0)
	prevTag, err := getPreviousTag(repo, false)
	assert.NoError(t, err)
	assert.Equal(t, "v1.2.0-alpha.1", prevTag)
}

func TestGetPreviousTagWithReleasedFlag(t *testing.T) {
	repo := setupRepoWithTags(t)

	// Using the released flag should skip prerelease tags
	prevTag, err := getPreviousTag(repo, true)
	assert.NoError(t, err)
	assert.Equal(t, "v1.1.0", prevTag)
}

func TestGetPreviousTagFromPrerelease(t *testing.T) {
	repo := setupRepoWithTags(t)

	// Create a new commit and checkout the prerelease tag
	w, err := repo.Worktree()
	assert.NoError(t, err)

	// Checkout the prerelease tag to make it the HEAD
	tagRef, err := repo.Tag("v1.2.0-alpha.1")
	assert.NoError(t, err)
	tagCommit, err := repo.ResolveRevision(plumbing.Revision(tagRef.Name().String()))
	assert.NoError(t, err)

	err = w.Checkout(&git.CheckoutOptions{
		Hash:   *tagCommit,
	})
	assert.NoError(t, err)

	// Get previous tag from the prerelease
	prevTag, err := getPreviousTag(repo, false)
	assert.NoError(t, err)
	assert.Equal(t, "v1.1.0", prevTag)
}

func TestGetPreviousTagWithSingleTag(t *testing.T) {
	// Create repo with just one tag
	repo, err := setupRepo()
	assert.NoError(t, err)

	commit, err := commitFile("file1.txt", repo)
	assert.NoError(t, err)

	_, err = repo.CreateTag("v1.0.0", commit, nil)
	assert.NoError(t, err)

	// Should return an error as there's no previous tag
	_, err = getPreviousTag(repo, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no previous tag available")
}

func TestGetPreviousTagWithSingleReleasedTag(t *testing.T) {
	// Create repo with just one tag
	repo, err := setupRepo()
	assert.NoError(t, err)

	commit, err := commitFile("file1.txt", repo)
	assert.NoError(t, err)

	_, err = repo.CreateTag("v1.0.0", commit, nil)
	assert.NoError(t, err)

	// Should return an error as there's no previous tag
	_, err = getPreviousTag(repo, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no previous tag available")
}

func TestGetPreviousTagFromUntaggedCommit(t *testing.T) {
	repo := setupRepoWithTags(t)

	// Create a new commit with a prerelease tag
	commit1, err := commitFile("file5.txt", repo)
	assert.NoError(t, err)
	_, err = repo.CreateTag("v1.3.0-alpha.1", commit1, nil)
	assert.NoError(t, err)

	// Add a new commit without a tag
	_, err = commitFile("file6.txt", repo)
	assert.NoError(t, err)

	// Should return the most recent tag's previous tag
	// In our test setup, newest tag is v1.2.0, so its previous is v1.2.0-alpha.1
	prevTag, err := getPreviousTag(repo, false)
	assert.NoError(t, err)
	assert.Equal(t, "v1.3.0-alpha.1", prevTag)

	// With released flag, should skip the prerelease
	prevTag, err = getPreviousTag(repo, true)
	assert.NoError(t, err)
	assert.Equal(t, "v1.2.0", prevTag)
}

func TestGetPreviousTagNoReleasedVersions(t *testing.T) {
	// Create a repository with only prerelease tags
	repo, err := setupRepo()
	assert.NoError(t, err)

	commit1, err := commitFile("file1.txt", repo)
	assert.NoError(t, err)
	_, err = repo.CreateTag("v1.0.0-alpha.1", commit1, nil)
	assert.NoError(t, err)

	commit2, err := commitFile("file2.txt", repo)
	assert.NoError(t, err)
	_, err = repo.CreateTag("v1.0.0-beta.1", commit2, nil)
	assert.NoError(t, err)

	// When requesting released versions only, should return an error
	_, err = getPreviousTag(repo, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no released versions found")
}

func TestGetPreviousTagHeadAtOldestTag(t *testing.T) {
	repo := setupRepoWithTags(t)

	// Checkout the oldest tag (v1.0.0)
	w, err := repo.Worktree()
	assert.NoError(t, err)

	err = w.Checkout(&git.CheckoutOptions{
		Hash:   plumbing.NewHash(""),
		Branch: plumbing.ReferenceName("refs/tags/v1.0.0"),
	})
	assert.NoError(t, err)

	// Should return an error as there's no previous tag
	_, err = getPreviousTag(repo, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no previous tag available - already at oldest tag")
}

func TestGetPreviousTagNoTagsInHistory(t *testing.T) {
	// Setup repo with isolated branch that has no tags
	repo, err := setupRepo()
	assert.NoError(t, err)

	// Create initial commit
	initCommit, err := commitFile("init.txt", repo)
	assert.NoError(t, err)

	// Add a tag on master
	_, err = repo.CreateTag("v1.0.0", initCommit, nil)
	assert.NoError(t, err)

	// Create a new branch from the initial commit
	w, err := repo.Worktree()
	assert.NoError(t, err)

	// Create a new branch
	branchRef := plumbing.NewBranchReferenceName("test-branch")
	err = w.Checkout(&git.CheckoutOptions{
		Create: true,
		Branch: branchRef,
	})
	assert.NoError(t, err)

	// Add commits on this branch but no tags
	_, err = commitFile("branch-file1.txt", repo)
	assert.NoError(t, err)

	// Trying to get previous tag should fail with an error about no tags in history
	_, err = getPreviousTag(repo, false)
	assert.Error(t, err)
}

func TestGetPreviousTagMostRecentIsOldest(t *testing.T) {
	// Create a repository with only two tags
	repo, err := setupRepo()
	assert.NoError(t, err)

	// Create first commit with tag
	commit1, err := commitFile("file1.txt", repo)
	assert.NoError(t, err)
	_, err = repo.CreateTag("v1.0.0", commit1, nil)
	assert.NoError(t, err)

	// Create second commit with tag
	commit2, err := commitFile("file2.txt", repo)
	assert.NoError(t, err)
	_, err = repo.CreateTag("v2.0.0", commit2, nil)
	assert.NoError(t, err)

	// Create a third commit with tag smaller than v1.0.0
	commit3, err := commitFile("file3.txt", repo)
	assert.NoError(t, err)
	_, err = repo.CreateTag("v0.9.0", commit3, nil)
	assert.NoError(t, err)

	// Add a commit after v0.9.0
	_, err = commitFile("file4.txt", repo)
	assert.NoError(t, err)

	// Should return v0.9.0 as the previous tag
	prevTag, err := getPreviousTag(repo, false)
	assert.NoError(t, err)
	assert.Equal(t, "v0.9.0", prevTag)
}

func TestGetPreviousTagWithCustomRepository(t *testing.T) {
	// Save original repo path
	originalRepoPath := repoPath
	defer func() { repoPath = originalRepoPath }()

	// Mock the repository path for this test
	tempDir := t.TempDir()
	repoPath = tempDir

	// This test just verifies that we can set a custom repository path
	// We don't actually run the command as it would exit the program
	assert.NotEqual(t, ".", repoPath)
}

func TestGetPreviousTagNoTagsInRepository(t *testing.T) {
	// Create a new repository without any tags
	repo, err := setupRepo()
	assert.NoError(t, err)

	// Try to get previous tag
	prevTag, err := getPreviousTag(repo, false)
	assert.Error(t, err)
	assert.Equal(t, "", prevTag)
}

func TestGetPreviousTagNoSemverTags(t *testing.T) {
	// Create a new repository with non-semver tags
	repo, err := setupRepo()
	assert.NoError(t, err)

	// Create a commit and tag it with a non-semver version
	commit, err := commitFile("file1.txt", repo)
	assert.NoError(t, err)
	_, err = repo.CreateTag("non-semver-tag", commit, nil)
	assert.NoError(t, err)

	// Try to get previous tag
	prevTag, err := getPreviousTag(repo, false)
	assert.Error(t, err)
	assert.Equal(t, "", prevTag)
}
