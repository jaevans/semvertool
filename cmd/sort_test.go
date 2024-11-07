package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortAscending(t *testing.T) {
	cmd := exec.Command("../semvertool", "sort", "1.0.0", "2.0.0", "0.1.0")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	assert.NoError(t, err)
	assert.Equal(t, "0.1.0 1.0.0 2.0.0\n", out.String())
}

func TestSortDescending(t *testing.T) {
	cmd := exec.Command("../semvertool", "sort", "--desc", "1.0.0", "2.0.0", "0.1.0")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	assert.NoError(t, err)
	assert.Equal(t, "2.0.0 1.0.0 0.1.0\n", out.String())
}

func TestSortWithSeparator(t *testing.T) {
	cmd := exec.Command("../semvertool", "sort", "--sep", ",", "1.0.0", "2.0.0", "0.1.0")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	assert.NoError(t, err)
	assert.Equal(t, "0.1.0,1.0.0,2.0.0\n", out.String())
}

func TestSortWithNewlineSeparator(t *testing.T) {
	cmd := exec.Command("../semvertool", "sort", "--sep", "\n", "1.0.0", "2.0.0", "0.1.0")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	assert.NoError(t, err)
	assert.Equal(t, "0.1.0\n1.0.0\n2.0.0\n", out.String())
}

func TestSortInvalidSemver(t *testing.T) {
	cmd := exec.Command("../semvertool", "sort", "1.0.0", "invalid", "0.1.0")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	assert.Error(t, err)
	assert.Contains(t, out.String(), "Invalid semver string: invalid")
}

func TestSortWithPrereleaseVersions(t *testing.T) {
	cmd := exec.Command("../semvertool", "sort", "1.0.0-alpha.1", "1.0.0-alpha.2", "1.0.0-beta.1")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0-alpha.1 1.0.0-alpha.2 1.0.0-beta.1\n", out.String())
}

func TestSortUnsortedInput(t *testing.T) {
	cmd := exec.Command("../semvertool", "sort", "2.0.0", "0.1.0", "1.0.0")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	assert.NoError(t, err)
	assert.Equal(t, "0.1.0 1.0.0 2.0.0\n", out.String())
}
