/*
copyright 2020 the Escli authors

licensed under the apache license, version 2.0 (the "license");
you may not use this file except in compliance with the license.
you may obtain a copy of the license at

    http://www.apache.org/licenses/license-2.0

unless required by applicable law or agreed to in writing, software
distributed under the license is distributed on an "as is" basis,
without warranties or conditions of any kind, either express or implied.
see the license for the specific language governing permissions and
limitations under the license.
*/

package builder

import (
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/DevopsArtFactory/escli/internal/constants"
	"github.com/DevopsArtFactory/escli/internal/util"
)

type Flag struct {
	Name          string
	Shorthand     string
	Usage         string
	Value         interface{}
	DefValue      interface{}
	FlagAddMethod string
	DefinedOn     []string
	Hidden        bool

	pflag *pflag.Flag
}

var FlagRegistry = []Flag{
	{
		Name:          "force",
		Shorthand:     "f",
		Usage:         "force transit storage class",
		Value:         aws.Bool(false),
		DefValue:      false,
		FlagAddMethod: "BoolVar",
		DefinedOn: []string{
			"snapshot restore [repositoryID] [snapshotID] [indexName]",
			"snapshot create [repositoryID] [snapshotID] [indices]",
			"snapshot delete [repositoryID] [snapshotID]",
			"snapshot archive [repositoryID] [snapshotID]",
			"index settings",
			"cluster settings",
			"index delete",
		},
	},
	{
		Name:          "troubled-only",
		Usage:         "show troubled only",
		Value:         aws.Bool(false),
		DefValue:      false,
		FlagAddMethod: "BoolVar",
		DefinedOn:     []string{"cat indices", "cat shards"},
	},
	{
		Name:          "sort-by",
		Shorthand:     "s",
		Usage:         "sort by field",
		Value:         aws.String(constants.EmptyString),
		DefValue:      constants.EmptyString,
		FlagAddMethod: "StringVar",
		DefinedOn:     []string{"cat indices", "cat nodes", "cat shards"},
	},
	{
		Name:          "repo-only",
		Usage:         "shows information of repository only",
		Value:         aws.Bool(false),
		DefValue:      false,
		FlagAddMethod: "BoolVar",
		DefinedOn:     []string{"snapshot list"},
	},
	{
		Name:          "with-repo",
		Usage:         "shows snapshots only repo",
		Value:         aws.String(constants.EmptyString),
		DefValue:      constants.EmptyString,
		FlagAddMethod: "StringVar",
		DefinedOn:     []string{"snapshot list"},
	},
	{
		Name:          "region",
		Usage:         "specify AWS region",
		Value:         aws.String(constants.EmptyString),
		DefValue:      constants.EmptyString,
		FlagAddMethod: "StringVar",
		DefinedOn: []string{
			"snapshot restore [repositoryID] [snapshotID] [indexName]",
			"snapshot archive [repositoryID] [snapshotID]"},
	},
	{
		Name:          "restore-tier",
		Usage:         "Tier of restore job",
		Value:         aws.String("Standard"),
		DefValue:      "Standard",
		FlagAddMethod: "StringVar",
		DefinedOn: []string{
			"snapshot restore [repositoryID] [snapshotID] [indexName]",
		},
	},
	{
		Name:          "max-concurrent-job",
		Usage:         "Maximum number of concurrent jobs for restoring snapshot",
		Value:         aws.Int64(constants.DefaultMaxConcurrentJob),
		DefValue:      int64(constants.DefaultMaxConcurrentJob),
		FlagAddMethod: "Int64Var",
		DefinedOn: []string{
			"snapshot restore [repositoryID] [snapshotID] [indexName]",
		},
	},
}

func (fl *Flag) flag() *pflag.Flag {
	if fl.pflag != nil {
		return fl.pflag
	}

	inputs := []interface{}{fl.Value, fl.Name}
	if fl.FlagAddMethod != "Var" {
		inputs = append(inputs, fl.DefValue)
	}
	inputs = append(inputs, fl.Usage)

	fs := pflag.NewFlagSet(fl.Name, pflag.ContinueOnError)
	reflect.ValueOf(fs).MethodByName(fl.FlagAddMethod).Call(reflectValueOf(inputs))
	f := fs.Lookup(fl.Name)
	f.Shorthand = fl.Shorthand
	f.Hidden = fl.Hidden

	fl.pflag = f
	return f
}

func reflectValueOf(values []interface{}) []reflect.Value {
	var results []reflect.Value
	for _, v := range values {
		results = append(results, reflect.ValueOf(v))
	}
	return results
}

//Add command flags
func SetCommandFlags(cmd *cobra.Command) {
	var flagsForCommand []*Flag
	for i := range FlagRegistry {
		fl := &FlagRegistry[i]

		if util.IsStringInArray(util.GetFullCommandUse(cmd), fl.DefinedOn) {
			cmd.Flags().AddFlag(fl.flag())
			flagsForCommand = append(flagsForCommand, fl)
		}
	}

	// Apply command-specific default values to flags.
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Update default values.
		for _, fl := range flagsForCommand {
			viper.BindPFlag(fl.Name, cmd.Flags().Lookup(fl.Name))
		}

		// Since PersistentPreRunE replaces the parent's PersistentPreRunE,
		// make sure we call it, if it is set.
		if parent := cmd.Parent(); parent != nil {
			if preRun := parent.PersistentPreRunE; preRun != nil {
				if err := preRun(cmd, args); err != nil {
					return err
				}
			} else if preRun := parent.PersistentPreRun; preRun != nil {
				preRun(cmd, args)
			}
		}

		return nil
	}
}
