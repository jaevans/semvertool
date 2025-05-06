package sort_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	semversort "github.com/jaevans/semvertool/pkg/sort"
	"github.com/stretchr/testify/assert"
)

// Returns the correct bump type when given a commit message containing a bump tag with a valid bump type.
func TestSortSimple(t *testing.T) {
	entries := semver.Collection{
		semver.MustParse("1.0.0"),
		semver.MustParse("2.0.0"),
		semver.MustParse("1.0.1"),
	}

	expected := semver.Collection{
		semver.MustParse("1.0.0"),
		semver.MustParse("1.0.1"),
		semver.MustParse("2.0.0"),
	}

	semversort.SortVersions(entries, true)

	assert.Equal(t, expected, entries)
}

func TestSortDescending(t *testing.T) {
	entries := semver.Collection{
		semver.MustParse("1.0.0"),
		semver.MustParse("2.0.0"),
		semver.MustParse("1.0.1"),
	}

	expected := semver.Collection{
		semver.MustParse("2.0.0"),
		semver.MustParse("1.0.1"),
		semver.MustParse("1.0.0"),
	}

	semversort.SortVersions(entries, false)

	assert.Equal(t, expected, entries)
}
