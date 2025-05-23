package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	goget "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type BumpType string

const (
	MajorBump      BumpType = "major"
	MinorBump      BumpType = "minor"
	PatchBump      BumpType = "patch"
	PrereleaseBump BumpType = "prerelease"
	UnknownBump    BumpType = "unknown"
	NoBump         BumpType = "none"
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
		}
		return messageBump

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
			vNew := v.IncPatch()
			vNew, err := vNew.SetPrerelease(viper.GetString("prerelease-prefix") + ".1")
			return &vNew, err
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

func getCommonBumpFlags() *pflag.FlagSet {
	commonFlags := pflag.NewFlagSet("common", pflag.ExitOnError)
	commonFlags.Bool("major", false, "Bump the major version")
	commonFlags.Bool("minor", false, "Bump the minor version")
	commonFlags.Bool("patch", false, "Bump the patch version")
	commonFlags.Bool("prerelease", false, "Bump the prerelease version")
	commonFlags.StringP("from-message", "m", "", "Extract the bump type from a commit message")
	commonFlags.StringP("prerelease-prefix", "p", "prerelease", "Set the prefix for the prerelease version if there is no existing prefix.")
	return commonFlags
}

func getTags(repo *goget.Repository) ([]*semver.Version, error) {
	iter, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	semverTags := make([]*semver.Version, 0)
	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		shortTag := ref.Name().Short()
		t, err := semver.NewVersion(shortTag)
		if err != nil {
			fmt.Printf("Could not parse tag %s as semver: %s\n", shortTag, err)
			return nil
		}
		semverTags = append(semverTags, t)
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Sort(semver.Collection(semverTags))
	return semverTags, nil
}

func getTagsStrings(repo *goget.Repository) ([]string, error) {
	tags, err := getTags(repo)
	if err != nil {
		return []string{}, err
	}
	tagStrings := make([]string, len(tags))
	for i, t := range tags {
		tagStrings[i] = t.Original()
	}
	return tagStrings, nil
}

func FilterPrerelease(versions []*semver.Version) []*semver.Version {
	// Filter out prerelease versions
	filtered := []*semver.Version{}
	for _, v := range versions {
		if len(v.Prerelease()) == 0 {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func VersionsToStrings(versions []*semver.Version) []string {
	result := make([]string, len(versions))
	for i, v := range versions {
		result[i] = v.Original()
	}
	return result
}
