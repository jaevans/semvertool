package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Returns the correct bump type when given a commit message containing a bump tag with a valid bump type.
func TestExtractBumpType_ValidBumpType(t *testing.T) {
	commitMessage := "[bump major]"
	expected := MajorBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractBumpType_ValidMinorBumpType(t *testing.T) {
	commitMessage := "[bump minor]"
	expected := MinorBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractBumpType_ValidPatchBumpType(t *testing.T) {
	commitMessage := "[bump patch]"
	expected := PatchBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractBumpType_ValidPrereleaseBumpType(t *testing.T) {
	commitMessage := "[bump prerelease]"
	expected := PrereleaseBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestExtractBumpType_ValidCaseInsensitive(t *testing.T) {
	commitMessage := "[bUmP preREleasE]"
	expected := PrereleaseBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

// Returns an empty string when given an empty string.
func TestExtractBumpType_EmptyString(t *testing.T) {
	commitMessage := ""
	expected := NoBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

// Returns an empty string when given a commit message that does not contain a bump tag.
func TestExtractBumpType_NoBumpTag(t *testing.T) {
	commitMessage := "Fix bug"
	expected := NoBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

func TestExtractBumpType_BumpNoType(t *testing.T) {
	commitMessage := "Fix bug [bump ]"
	expected := NoBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

// Returns an empty string when given a commit message containing a bump tag with an invalid bump type.
func TestExtractBumpType_InvalidBumpType(t *testing.T) {
	commitMessage := "[bump invalid]"
	expected := NoBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

// Returns the correct bump type when given a commit message containing multiple bump tags with valid bump types.
func TestExtractBumpType_MultipleValidBumpTypes(t *testing.T) {
	commitMessage := "[bump major] [bump minor] [bump patch]"
	expected := MajorBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

// Returns the correct bump type when given a commit message containing a bump tag with a valid bump type but with additional whitespace.
func TestExtractBumpType_ValidBumpTypeWithWhitespace(t *testing.T) {
	commitMessage := "[bump   minor   ]"
	expected := MinorBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

// Returns an empty string when given a commit message containing a bump tag with a valid bump type but with additional non-whitespace characters.
func TestExtractBumpType_InvalidBumpTypeWithAdditionalCharactersAfter(t *testing.T) {
	commitMessage := "[bump majorabc]"
	expected := NoBump

	result := extractBumpTypeFromMessage(commitMessage)

	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

// Returns an empty string when given a commit message containing a bump tag with a valid bump type but with additional characters before the tag.
func TestExtractBumpType_InvalidBumpTypeWithAdditionalCharactersBefore(t *testing.T) {
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

func TestGetBumpType_MajorBump(t *testing.T) {
	viper.Reset()
	viper.Set("major", true)
	expected := MajorBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBumpType_MinorBump(t *testing.T) {
	viper.Reset()
	viper.Set("minor", true)
	expected := MinorBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBumpType_PatchBump(t *testing.T) {
	viper.Reset()
	viper.Set("patch", true)
	expected := PatchBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBumpType_PrereleaseBump(t *testing.T) {
	viper.Reset()
	viper.Set("prerelease", true)
	expected := PrereleaseBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBumpType_NothingSet(t *testing.T) {
	viper.Reset()
	expected := PatchBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBUmpType_FromMessage(t *testing.T) {
	viper.Reset()
	viper.Set("from-message", "[bump major]")
	expected := MajorBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBUmpType_FromEmptyMessage(t *testing.T) {
	viper.Reset()
	viper.Set("from-message", "")
	expected := PatchBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBUmpType_FromEmptyBumpMessage(t *testing.T) {
	viper.Reset()
	viper.Set("from-message", "[bump]")
	expected := NoBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

func TestGetBUmpType_FromInvalidBumpMessage(t *testing.T) {
	viper.Reset()
	viper.Set("from-message", "[bump invalid]")
	expected := NoBump
	result := getBumpType()
	assert.Equal(t, expected, result)
}

// **********************
// doBump
// **********************

func TestDoBump_ValidVersion(t *testing.T) {
	version := "1.2.3"
	bumpType := MajorBump
	expected := "2.0.0"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBump_ValidVersionMinorBump(t *testing.T) {
	version := "1.2.3"
	bumpType := MinorBump
	expected := "1.3.0"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBump_ValidVersionPatchBump(t *testing.T) {
	version := "1.2.3"
	bumpType := PatchBump
	expected := "1.2.4"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBump_ValidVersionPrereleaseBump(t *testing.T) {
	version := "1.2.3-alpha.1"
	bumpType := PrereleaseBump
	expected := "1.2.3-alpha.2"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBump_ValidVersionPrereleaseBumpNoPrerelease(t *testing.T) {
	viper.Reset()
	viper.Set("prerelease-prefix", "alpha")
	version := "1.2.3"
	bumpType := PrereleaseBump
	expected := "1.2.4-alpha.1"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBump_ValidVersionPrereleaseBumpValidPrereleaseNoNumber(t *testing.T) {
	version := "1.2.3-alpha"
	bumpType := PrereleaseBump
	expected := "1.2.3-alpha.0"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBump_ValidVersionPrereleaseBumpValidPrereleaseNoDot(t *testing.T) {
	version := "1.2.3-alpha0"
	bumpType := PrereleaseBump
	expected := "1.2.3-alpha.1"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBump_ValidVersionPrereleaseBumpValidPrereleaseMultipleDots(t *testing.T) {
	version := "1.2.3-alpha.0.9"
	bumpType := PrereleaseBump
	expected := "1.2.3-alpha.0.10"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBump_InvalidVersion(t *testing.T) {
	version := "1.x.3-alpha"
	bumpType := MajorBump
	_, err := doBump(version, bumpType)
	assert.Error(t, err)
}

func TestDoBump_PrereleaseClearsBuild(t *testing.T) {
	version := "1.2.3-alpha.0+build"
	bumpType := PrereleaseBump
	expected := "1.2.3-alpha.1"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}
func TestDoBump_PatchClearsPrerelease(t *testing.T) {
	viper.Reset()
	version := "1.2.3-alpha.0"
	bumpType := PatchBump
	expected := "1.2.3"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}
func TestDoBump_MinorClearsPrerelease(t *testing.T) {
	version := "1.2.3-alpha.0"
	bumpType := MinorBump
	expected := "1.3.0"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBump_MajorClearsPrerelease(t *testing.T) {
	version := "1.2.3-alpha.0"
	bumpType := MajorBump
	expected := "2.0.0"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}

func TestDoBump_AddsCorrectPrereleasePrefix(t *testing.T) {
	viper.Reset()
	viper.Set("prerelease-prefix", "snapshot")
	version := "1.2.3"
	bumpType := PrereleaseBump
	expected := "1.2.4-snapshot.1"
	result, err := doBump(version, bumpType)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.String())
}
