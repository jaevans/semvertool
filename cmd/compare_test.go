package cmd

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		version1 string
		version2 string
		expected int
	}{
		{"FirstSmaller", "1.0.0", "2.0.0", 0},
		{"Equal", "2.0.0", "2.0.0", 1},
		{"SecondSmaller", "2.0.0", "1.0.0", 2},
		{"PreReleaseFirstSmaller", "1.0.0-alpha", "1.0.0-beta", 0},
		{"PreReleaseEqual", "1.0.0-alpha", "1.0.0-alpha", 1},
		{"PreReleaseSecondSmaller", "1.0.0-beta", "1.0.0-alpha", 2},
		{"PreReleaseAlpha1Smaller", "1.0.0-alpha.1", "1.0.0-alpha.8", 0},
		{"PreReleaseAlphaEqual", "1.0.0-alpha.1", "1.0.0-alpha.1", 1},
		{"PreReleaseAlpha8Smaller", "1.0.0-alpha.8", "1.0.0-alpha.1", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", "../main.go", "compare", tt.version1, tt.version2)
			err := cmd.Run()
			exitError, ok := err.(*exec.ExitError)
			if err == nil {
				assert.Equal(t, tt.expected, 1)
			} else if !ok {
				t.Fatalf("expected exit error, got %v", err)
			} else {
				assert.Equal(t, tt.expected, exitError.ExitCode())
			}
		})
	}
}
