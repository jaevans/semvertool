package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type BumpType string

const (
	MajorBump      BumpType = "major"
	MinorBump      BumpType = "minor"
	PatchBump      BumpType = "patch"
	PrereleaseBump BumpType = "prerelease"
	UnknownBump    BumpType = ""
	NoBump         BumpType = ""
)

func extractBumpTypeFromMessage(s string) BumpType {
	re := regexp.MustCompile(`(?i)\[bump\s*(major|minor|patch|prerelease)\s*\]`)

	matches := re.FindStringSubmatch(s)

	if len(matches) < 2 {
		return NoBump
	}
	switch strings.ToLower(matches[1]) {
	case "major":
		return MajorBump
	case "minor":
		return MinorBump
	case "patch":
		return PatchBump
	case "prerelease":
		return PrereleaseBump
	}

	// Can't actually get here, you've either failed the regex, or matched one of the above
	return UnknownBump
}

var ErrNoTrailingDigits = errors.New("no trailing digits found")

func extractTrailingDigits(s string) (string, int, error) {
	// re := regexp.MustCompile(`(?i)([^0-9\.]*)\.?(\d+)$`)
	re := regexp.MustCompile(`(?i)(.*\D)?(\d+)$`)

	matches := re.FindStringSubmatch(s)

	if len(matches) < 3 {
		return "", -1, ErrNoTrailingDigits
	}

	number, err := strconv.Atoi(matches[2])
	return strings.TrimSuffix(matches[1], "."), number, err
}

func getBumpType() BumpType {
	if fromMessage := viper.GetString("from-message"); fromMessage != "" {
		messageBump := extractBumpTypeFromMessage(fromMessage)
		if messageBump == UnknownBump || messageBump == NoBump {
			fmt.Println("No valid bump type found in the commit message")
		} else {
			return messageBump
		}
	}
	if viper.GetBool("major") {
		return MajorBump
	}
	if viper.GetBool("minor") {
		return MinorBump
	}
	if viper.GetBool("patch") {
		return PatchBump
	}
	if viper.GetBool("prerelease") {
		return PrereleaseBump
	}
	return PatchBump
}

func doBump(version string, bumpWhat BumpType) (*semver.Version, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return &semver.Version{}, err
	}

	switch bumpWhat {
	case MajorBump:
		vNew := v.IncMajor()
		v = &vNew
	case MinorBump:
		vNew := v.IncMinor()
		v = &vNew
	case PatchBump:
		vNew := v.IncPatch()
		v = &vNew
	case PrereleaseBump:
		prerelease := v.Prerelease()
		if len(prerelease) == 0 {
			fmt.Println("No prerelease found, bumping patch version")
			vNew := v.IncPatch()
			v = &vNew
			return v, nil
		}
		prefix, number, err := extractTrailingDigits(prerelease)
		if err == ErrNoTrailingDigits && !strings.Contains(prefix, ".") {
			prefix = prerelease
			number = -1 // to get us to zero when we bump
		} else if err != nil {
			return &semver.Version{}, err
		}
		number++
		vNew, err := v.SetPrerelease(fmt.Sprintf("%s.%d", prefix, number))
		if err != nil {
			return &semver.Version{}, err
		}
		v = &vNew

		// clear the build information
		vNew, err = v.SetMetadata("")
		if err != nil {
			return &semver.Version{}, err
		}

	}
	return v, nil
}

func addCommonBumpFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("major", false, "Bump the major version")
	cmd.Flags().Bool("minor", false, "Bump the minor version")
	cmd.Flags().Bool("patch", false, "Bump the patch version")
	cmd.Flags().Bool("prerelease", false, "Bump the prerelease version")
	cmd.Flags().StringP("from-message", "m", "", "Extract the bump type from a commit message")
	cmd.MarkFlagsMutuallyExclusive("major", "minor", "patch", "prerelease", "from-message")
	viper.BindPFlags(cmd.Flags())
}
