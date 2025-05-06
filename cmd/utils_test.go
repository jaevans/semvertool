package cmd

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Returns the correct bump type when given a commit message containing a bump tag with a valid bump type.
func TestExtractBumpTypeValidBumpType(t *testing.T) {
	commitMessage := "[bump major]"
	expected := MajorBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractBumpTypeValidMinorBumpType(t *testing.T) {
	commitMessage := "[bump minor]"
	expected := MinorBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractBumpTypeValidPatchBumpType(t *testing.T) {
	commitMessage := "[bump patch]"
	expected := PatchBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractBumpTypeValidPrereleaseBumpType(t *testing.T) {
	commitMessage := "[bump prerelease]"
	expected := PrereleaseBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractBumpTypeValidCaseInsensitive(t *testing.T) {
	commitMessage := "[bUmP preREleasE]"
	expected := PrereleaseBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

// Returns an empty string when given an empty string.
func TestExtractBumpTypeEmptyString(t *testing.T) {
	commitMessage := ""
	expected := NoBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

// Returns an empty string when given a commit message that does not contain a bump tag.
func TestExtractBumpTypeNoBumpTag(t *testing.T) {
	commitMessage := "Fix bug"
	expected := NoBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

func TestExtractBumpTypeBumpNoType(t *testing.T) {
	commitMessage := "Fix bug [bump ]"
	expected := NoBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

// Returns an empty string when given a commit message containing a bump tag with an invalid bump type.
func TestExtractBumpTypeInvalidBumpType(t *testing.T) {
	commitMessage := "[bump invalid]"
	expected := NoBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

// Returns the correct bump type when given a commit message containing multiple bump tags with valid bump types.
func TestExtractBumpTypeMultipleValidBumpTypes(t *testing.T) {
	commitMessage := "[bump major] [bump minor] [bump patch]"
	expected := MajorBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

// Returns the correct bump type when given a commit message containing a bump tag with a valid bump type but with additional whitespace.
func TestExtractBumpTypeValidBumpTypeWithWhitespace(t *testing.T) {
	commitMessage := "[bump   minor   ]"
	expected := MinorBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

// Returns an empty string when given a commit message containing a bump tag with a valid bump type but with additional non-whitespace characters.
func TestExtractBumpTypeInvalidBumpTypeWithAdditionalCharactersAfter(t *testing.T) {
	commitMessage := "[bump majorabc]"
	expected := NoBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

// Returns an empty string when given a commit message containing a bump tag with a valid bump type but with additional characters before the tag.
func TestExtractBumpTypeInvalidBumpTypeWithAdditionalCharactersBefore(t *testing.T) {
	commitMessage := "Some additional characters [bump asdf major]"
	expected := NoBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

// **********************
// ExtractTrailingDigits
// **********************
func TestExtractTrailingDigitsWithDigitsAtEnd(t *testing.T) {
	s := "abc123"
	prefix, number, err := extractTrailingDigits(s)
	assert.NoError(t, err)
	assert.Equal(t, "abc", prefix)
	assert.Equal(t, 123, number)
}

func TestExtractTrailingDigitsWithNoTrailingDigits(t *testing.T) {
	s := "abc"
	_, _, err := extractTrailingDigits(s)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrNoTrailingDigits.Error())
}

// Should return an error when the input string is empty
func TestExtractTrailingDigitsWithEmptyString(t *testing.T) {
	s := ""
	_, _, err := extractTrailingDigits(s)
	assert.Error(t, err)
}

// Should return a single digit when the input string is a single digit
func TestExtractTrailingDigitsWithSingleDigit(t *testing.T) {
	s := "1"
	prefix, number, err := extractTrailingDigits(s)
	assert.NoError(t, err)
	assert.Equal(t, "", prefix)
	assert.Equal(t, 1, number)
}

func TestExtractTrailingDigitsWithDotDigit(t *testing.T) {
	s := "alpha.1"
	prefix, number, err := extractTrailingDigits(s)
	assert.NoError(t, err)
	assert.Equal(t, "alpha", prefix)
	assert.Equal(t, 1, number)
}

// **********************
// getBumpType
// **********************

func TestGetBumpTypeMajorBump(t *testing.T) {
	viper.Reset()
	viper.Set("major", true)
	expected := MajorBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBumpTypeMinorBump(t *testing.T) {
	viper.Reset()
	viper.Set("minor", true)
	expected := MinorBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBumpTypePatchBump(t *testing.T) {
	viper.Reset()
	viper.Set("patch", true)
	expected := PatchBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBumpTypePrereleaseBump(t *testing.T) {
	viper.Reset()
	viper.Set("prerelease", true)
	expected := PrereleaseBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBumpTypeNothingSet(t *testing.T) {
	viper.Reset()
	expected := PatchBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBUmpTypeFromMessage(t *testing.T) {
	viper.Reset()
	viper.Set("from-message", "[bump major]")
	expected := MajorBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBUmpTypeFromEmptyMessage(t *testing.T) {
	viper.Reset()
	viper.Set("from-message", "")
	expected := PatchBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBUmpTypeFromEmptyBumpMessage(t *testing.T) {
	viper.Reset()
	viper.Set("from-message", "[bump]")
	expected := NoBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBUmpTypeFromInvalidBumpMessage(t *testing.T) {
	viper.Reset()
	viper.Set("from-message", "[bump invalid]")
	expected := NoBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

// **********************
// doBump
// **********************

func TestDoBumpValidVersion(t *testing.T) {
	version := "1.2.3"
	bumpType := MajorBump
	expected := "2.0.0"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBumpValidVersionMinorBump(t *testing.T) {
	version := "1.2.3"
	bumpType := MinorBump
	expected := "1.3.0"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBumpValidVersionPatchBump(t *testing.T) {
	version := "1.2.3"
	bumpType := PatchBump
	expected := "1.2.4"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBumpValidVersionPrereleaseBump(t *testing.T) {
	version := "1.2.3-alpha.1"
	bumpType := PrereleaseBump
	expected := "1.2.3-alpha.2"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBumpValidVersionPrereleaseBumpNoPrerelease(t *testing.T) {
	viper.Reset()
	viper.Set("prerelease-prefix", "alpha")
	version := "1.2.3"
	bumpType := PrereleaseBump
	expected := "1.2.4-alpha.1"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBumpValidVersionPrereleaseBumpValidPrereleaseNoNumber(t *testing.T) {
	version := "1.2.3-alpha"
	bumpType := PrereleaseBump
	expected := "1.2.3-alpha.0"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBumpValidVersionPrereleaseBumpValidPrereleaseNoDot(t *testing.T) {
	version := "1.2.3-alpha0"
	bumpType := PrereleaseBump
	expected := "1.2.3-alpha.1"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBumpValidVersionPrereleaseBumpValidPrereleaseMultipleDots(t *testing.T) {
	version := "1.2.3-alpha.0.9"
	bumpType := PrereleaseBump
	expected := "1.2.3-alpha.0.10"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBumpInvalidVersion(t *testing.T) {
	version := "1.x.3-alpha"
	bumpType := MajorBump
	_, err := doBump(version, bumpType)
	assert.Error(t, err)
}

func TestDoBumpPrereleaseClearsBuild(t *testing.T) {
	version := "1.2.3-alpha.0+build"
	bumpType := PrereleaseBump
	expected := "1.2.3-alpha.1"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}
func TestDoBumpPatchClearsPrerelease(t *testing.T) {
	viper.Reset()
	version := "1.2.3-alpha.0"
	bumpType := PatchBump
	expected := "1.2.3"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}
func TestDoBumpMinorClearsPrerelease(t *testing.T) {
	version := "1.2.3-alpha.0"
	bumpType := MinorBump
	expected := "1.3.0"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBumpMajorClearsPrerelease(t *testing.T) {
	version := "1.2.3-alpha.0"
	bumpType := MajorBump
	expected := "2.0.0"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBumpAddsCorrectPrereleasePrefix(t *testing.T) {
	viper.Reset()
	viper.Set("prerelease-prefix", "snapshot")
	version := "1.2.3"
	bumpType := PrereleaseBump
	expected := "1.2.4-snapshot.1"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestFilterPrereleasesNoPrerelease(t *testing.T) {
	entries := []*semver.Version{
		semver.MustParse("1.0.0"),
		semver.MustParse("1.0.1"),
	}
	expected := []*semver.Version{
		semver.MustParse("1.0.0"),
		semver.MustParse("1.0.1"),
	}

	result := FilterPrerelease(entries)
	assert.Equal(t, expected, result)
}

func TestFilterPrereleasesWithPrerelease(t *testing.T) {
	entries := []*semver.Version{
		semver.MustParse("1.0.0"),
		semver.MustParse("1.0.1-alpha"),
		semver.MustParse("1.0.2-beta"),
	}
	expected := []*semver.Version{
		semver.MustParse("1.0.0"),
	}

	result := FilterPrerelease(entries)
	assert.Equal(t, expected, result)
}

func TestFilterPrereleasesEmpty(t *testing.T) {
	entries := []*semver.Version{}
	expected := []*semver.Version{}

	result := FilterPrerelease(entries)
	assert.Equal(t, expected, result)
}

func TestFilterPrereleasesAllPrerelease(t *testing.T) {
	entries := []*semver.Version{
		semver.MustParse("1.0.0-alpha"),
		semver.MustParse("1.0.1-beta"),
	}
	expected := []*semver.Version{}

	result := FilterPrerelease(entries)
	assert.Equal(t, expected, result)
}

func TestVersionsToStrings(t *testing.T) {
	entries := []*semver.Version{
		semver.MustParse("1.0.0"),
		semver.MustParse("1.0.1"),
	}
	expected := []string{
		"1.0.0",
		"1.0.1",
	}

	result := VersionsToStrings(entries)
	assert.Equal(t, expected, result)
}
