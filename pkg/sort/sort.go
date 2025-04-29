package sort

import (
	"sort"

	"github.com/Masterminds/semver/v3"
)

func SortVersions(versions semver.Collection, ascending bool) {
	if ascending {
		sort.Sort(versions)
	} else {
		sort.Sort(sort.Reverse(versions))
	}
}
