package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// sortCmd represents the sort command
var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort a list of semver strings",
	Long: `Sort a list of semver strings in ascending or descending order.

Examples:
semvertool sort 1.0.0 2.0.0 0.1.0
0.1.0 1.0.0 2.0.0

semvertool sort --descending 1.0.0 2.0.0 0.1.0
2.0.0 1.0.0 0.1.0

semvertool sort --separator "," 1.0.0 2.0.0 0.1.0
0.1.0,1.0.0,2.0.0
`,
	Run: runSort,
}

func init() {
	rootCmd.AddCommand(sortCmd)

	sortCmd.Flags().BoolP("descending", "d", false, "Sort in descending order")
	sortCmd.Flags().StringP("separator", "s", " ", "Separator for the output")
}

func runSort(cmd *cobra.Command, args []string) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(flag.Name, flag)
	})

	if len(args) < 1 {
		_ = cmd.Help()
		return
	}

	versions := make([]*semver.Version, len(args))
	for i, arg := range args {
		v, err := semver.NewVersion(arg)
		if err != nil {
			fmt.Printf("Invalid semver string: %s\n", arg)
			return
		}
		versions[i] = v
	}

	sort.Slice(versions, func(i, j int) bool {
		if viper.GetBool("descending") {
			return versions[i].LessThan(versions[j])
		}
		return versions[j].LessThan(versions[i])
	})

	sep := viper.GetString("separator")
	versionStrings := make([]string, len(versions))
	for i, v := range versions {
		versionStrings[i] = v.String()
	}

	fmt.Println(strings.Join(versionStrings, sep))
}
