package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for CompareVersions function
func TestCompareVersionsLessThan(t *testing.T) {
	result, err := CompareVersions("1.0.0", "2.0.0")
	assert.NoError(t, err)
	assert.Equal(t, 12, result) // v2 is newer
}

func TestCompareVersionsEqual(t *testing.T) {
	result, err := CompareVersions("1.0.0", "1.0.0")
	assert.NoError(t, err)
	assert.Equal(t, 0, result) // equal versions
}

func TestCompareVersionsGreaterThan(t *testing.T) {
	result, err := CompareVersions("2.0.0", "1.0.0")
	assert.NoError(t, err)
	assert.Equal(t, 11, result) // v1 is newer
}

func TestCompareVersionsFirstInvalid(t *testing.T) {
	_, err := CompareVersions("invalid", "1.0.0")
	assert.Error(t, err)
}

func TestCompareVersionsSecondInvalid(t *testing.T) {
	_, err := CompareVersions("1.0.0", "invalid")
	assert.Error(t, err)
}

func TestCompareVersionsPatchVersions(t *testing.T) {
	result, err := CompareVersions("1.0.1", "1.0.2")
	assert.NoError(t, err)
	assert.Equal(t, 12, result) // v2 is newer
}

func TestCompareVersionsPrereleaseVsRelease(t *testing.T) {
	result, err := CompareVersions("1.0.0-alpha", "1.0.0")
	assert.NoError(t, err)
	assert.Equal(t, 12, result) // v2 is newer
}

// Tests for IsReleased function
func TestIsReleasedSimpleReleaseVersion(t *testing.T) {
	result, err := IsReleased("1.0.0")
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestIsReleasedPrereleaseVersion(t *testing.T) {
	result, err := IsReleased("1.0.0-alpha.1")
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestIsReleasedVersionWithMetadata(t *testing.T) {
	result, err := IsReleased("1.0.0+20130313144700")
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestIsReleasedPrereleaseWithMetadata(t *testing.T) {
	result, err := IsReleased("1.0.0-beta.1+exp.sha.5114f85")
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestIsReleasedInvalidVersion(t *testing.T) {
	_, err := IsReleased("invalid")
	assert.Error(t, err)
}

func TestIsReleasedComplexVersion(t *testing.T) {
	result, err := IsReleased("2.1.0-rc.2")
	assert.NoError(t, err)
	assert.False(t, result)
}
