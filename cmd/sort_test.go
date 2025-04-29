package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestSort_RunSortSimple(t *testing.T) {
	args := []string{"1.0.0", "2.0.0", "1.0.1"}
	expected := "1.0.0 1.0.1 2.0.0"

	// Capture stdout
	r, w, _ := os.Pipe()
	originalStdout := os.Stdout
	os.Stdout = w

	err := RunSort(&cobra.Command{}, args)
	assert.NoError(t, err)

	// Close the writer and restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Read the captured output
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	assert.NoError(t, err)

	assert.Equal(t, expected, strings.TrimSuffix(buf.String(), "\n"))
}

func TestSort_RunSortDescending(t *testing.T) {
	args := []string{"1.0.0", "2.0.0", "1.0.1"}
	expected := "2.0.0 1.0.1 1.0.0"

	// Capture stdout
	r, w, _ := os.Pipe()
	originalStdout := os.Stdout
	os.Stdout = w

	oldOrder := order
	order = "descending"
	err := RunSort(&cobra.Command{}, args)
	assert.NoError(t, err)
	order = oldOrder // Restore the original order

	// Close the writer and restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Read the captured output
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	assert.NoError(t, err)

	assert.Equal(t, expected, strings.TrimSuffix(buf.String(), "\n"))
}

func TestSort_RunSortNoPrerelease(t *testing.T) {
	args := []string{"1.0.0", "2.0.0", "1.0.1-alpha.1"}
	expected := "1.0.0 2.0.0"

	// Capture stdout
	r, w, _ := os.Pipe()
	originalStdout := os.Stdout
	os.Stdout = w

	oldNoPrerelease := noPrerelease
	noPrerelease = true
	err := RunSort(&cobra.Command{}, args)
	assert.NoError(t, err)
	noPrerelease = oldNoPrerelease // Restore the original noPrerelease

	// Close the writer and restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Read the captured output
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	assert.NoError(t, err)

	assert.Equal(t, expected, strings.TrimSuffix(buf.String(), "\n"))
}

func TestSort_RunSortInvalidTag(t *testing.T) {
	args := []string{"1.0.0", "2.0.0", "1.0.1-alpha.1", "invalid-tag"}
	expected := "invalid semver version: invalid-tag"

	// Capture stdout
	r, w, _ := os.Pipe()
	originalStdout := os.Stdout
	os.Stdout = w

	// Capture stdout
	re, we, _ := os.Pipe()
	originalStderr := os.Stderr
	os.Stderr = we

	err := RunSort(&cobra.Command{}, args)
	assert.NoError(t, err)

	// Close the writer and restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Close the writer and restore stderr
	we.Close()
	os.Stderr = originalStderr

	// Read the captured output
	var buf bytes.Buffer
	_, err = buf.ReadFrom(re)
	assert.NoError(t, err)

	assert.Equal(t, expected, strings.TrimSuffix(buf.String(), "\n"))

	_, err = buf.ReadFrom(r)
	assert.NoError(t, err)
}
