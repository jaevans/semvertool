/*
Copyright © 2025 James Evans
*/
package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		commit         string
		date           string
		expectedOutput []string // strings that should be in output
		notExpected    []string // strings that should NOT be in output
	}{
		{
			name:           "dev build",
			version:        "dev",
			commit:         "none",
			date:           "unknown",
			expectedOutput: []string{"dev", "snapshot", "commit:", "built:"},
			notExpected:    []string{},
		},
		{
			name:           "tagged release",
			version:        "v1.2.3",
			commit:         "",
			date:           "",
			expectedOutput: []string{"v1.2.3"},
			notExpected:    []string{"snapshot", "commit:", "built:"},
		},
		{
			name:           "snapshot with commit",
			version:        "v0.0.0-next",
			commit:         "abc123def",
			date:           "2025-01-01T10:00:00Z",
			expectedOutput: []string{"v0.0.0-next", "snapshot", "commit:", "abc123def", "built:", "2025-01-01T10:00:00Z"},
			notExpected:    []string{},
		},
		{
			name:           "tagged release with commit set to none",
			version:        "v2.0.0",
			commit:         "none",
			date:           "2025-01-01T10:00:00Z",
			expectedOutput: []string{"v2.0.0"},
			notExpected:    []string{"snapshot", "commit:", "built:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			origVersion := version
			origCommit := commit
			origDate := date

			// Set test values
			version = tt.version
			commit = tt.commit
			date = tt.date

			// Restore original values after test
			defer func() {
				version = origVersion
				commit = origCommit
				date = origDate
			}()

			// Capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the command
			runVersion(versionCmd, []string{})

			// Restore stdout
			w.Close()
			os.Stdout = old

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Check expected strings are present
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, got: %s", expected, output)
				}
			}

			// Check unexpected strings are NOT present
			for _, notExpected := range tt.notExpected {
				if strings.Contains(output, notExpected) {
					t.Errorf("Expected output NOT to contain %q, got: %s", notExpected, output)
				}
			}
		})
	}
}
