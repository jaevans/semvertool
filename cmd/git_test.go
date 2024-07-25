package cmd

import (
	"fmt"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// func TestValidVersionAndRepository(t *testing.T) {
// 	// Mock the necessary dependencies
// 	repo := &goget.Repository{}
// 	version := "1.2.3"
// 	expectedVersion := semver.MustParse("2.0.0")

// 	repo = goget.
// 	result, err := gitBump(version, repo)

// 	// Assert the result
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedVersion, result)
// }

func setupRepo() (*git.Repository, error) {
	fs := memfs.New()
	repo, err := git.Init(memory.NewStorage(), fs)

	return repo, err
}

func commitFile(filename string, repo *git.Repository) (plumbing.Hash, error) {
	w, err := repo.Worktree()
	if err != nil {
		return plumbing.Hash{}, err
	}

	fs := w.Filesystem

	_, err = fs.Create(filename)
	if err != nil {
		return plumbing.Hash{}, err
	}

	_, err = w.Add(filename)
	if err != nil {
		return plumbing.Hash{}, err
	}

	return w.Commit(fmt.Sprintf("Commit of %s", filename), &git.CommitOptions{})
}
func TestValidVersionAndRepostory(t *testing.T) {
	repo, err := setupRepo()
	assert.NoError(t, err)

	commit, err := commitFile("file1.txt", repo)
	assert.NoError(t, err)

	_, err = repo.CreateTag("v1.0.0", commit, nil)
	assert.NoError(t, err)

	commit, err = commitFile("file2.txt", repo)
	assert.NoError(t, err)

	_, err = repo.CreateTag("1.0.1", commit, nil)
	assert.NoError(t, err)

	expected := []string{"v1.0.0", "1.0.1"}
	result, err := getTags(repo)
	resultStrings := make([]string, len(result))
	for i, r := range result {
		resultStrings[i] = r.Original()
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, resultStrings)
}

func TestNoTags(t *testing.T) {
	repo, err := setupRepo()
	assert.NoError(t, err)

	expected := []string{}
	result, err := getTags(repo)
	resultStrings := make([]string, len(result))
	for i, r := range result {
		resultStrings[i] = r.Original()
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, resultStrings)
}

func TestInvalidVersionAndRepository(t *testing.T) {
	repo, err := setupRepo()
	assert.NoError(t, err)

	commit, err := commitFile("file1.txt", repo)
	assert.NoError(t, err)

	_, err = repo.CreateTag("v1.0.0", commit, nil)
	assert.NoError(t, err)

	commit, err = commitFile("file2.txt", repo)
	assert.NoError(t, err)

	_, err = repo.CreateTag("foobarbaz", commit, nil)
	assert.NoError(t, err)

	expected := []string{"v1.0.0"}
	result, err := getTags(repo)
	resultStrings := make([]string, len(result))
	for i, r := range result {
		resultStrings[i] = r.Original()
	}

	assert.NoError(t, err)
	fmt.Println(result)
	assert.Equal(t, expected, resultStrings)
}

func TestGitBump_NoTags(t *testing.T) {
	repo, err := setupRepo()
	assert.NoError(t, err)

	_, err = gitBump(repo)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNoSemverTags)
}

func TestGitBump_OneTag(t *testing.T) {
	repo, err := setupRepo()
	assert.NoError(t, err)

	commit, err := commitFile("file1.txt", repo)
	assert.NoError(t, err)

	_, err = repo.CreateTag("v1.0.0", commit, nil)
	assert.NoError(t, err)

	result, err := gitBump(repo)
	expected := "v1.0.1"
	assert.NoError(t, err)
	assert.Equal(t, expected, result.Original())
}

func TestGitBump_MultipleTags(t *testing.T) {
	repo, err := setupRepo()
	assert.NoError(t, err)

	commit, err := commitFile("file1.txt", repo)
	assert.NoError(t, err)

	_, err = repo.CreateTag("v1.0.0", commit, nil)
	assert.NoError(t, err)

	commit, err = commitFile("file2.txt", repo)
	assert.NoError(t, err)

	_, err = repo.CreateTag("v1.1.0", commit, nil)
	assert.NoError(t, err)

	result, err := gitBump(repo)
	expected := "v1.1.1"
	assert.NoError(t, err)
	assert.Equal(t, expected, result.Original())
}

func TestGitBump_WithHash(t *testing.T) {
	repo, err := setupRepo()
	assert.NoError(t, err)

	commitHash, err := commitFile("file1.txt", repo)
	assert.NoError(t, err)

	_, err = repo.CreateTag("v1.0.0", commitHash, nil)
	assert.NoError(t, err)

	shortHash := commitHash.String()[:7]

	viper.Set("hash", true)
	result, err := gitBump(repo)
	expected := "v1.0.1+" + shortHash
	assert.NoError(t, err)
	assert.Equal(t, expected, result.Original())
}
