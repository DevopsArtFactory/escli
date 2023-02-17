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

package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/DevopsArtFactory/escli/internal/constants"
)

var (
	cfgFile string
	profile string
)

func NewRootCommand(out, stderr io.Writer) *cobra.Command {
	cobra.OnInitialize(initConfig)
	rootCmd := &cobra.Command{
		Use:           "escli",
		Short:         "manage elasticsearch cluster",
		Long:          "manage elasticsearch cluster",
		SilenceErrors: true,
		SilenceUsage:  false,
	}

	rootCmd.AddCommand(NewSnapshotCommand())
	rootCmd.AddCommand(NewCatCommand())
	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(NewDiagCommand())
	rootCmd.AddCommand(NewFixCommand())
	rootCmd.AddCommand(NewIndexCommand())
	rootCmd.AddCommand(NewClusterCommand())
	rootCmd.AddCommand(NewUpdateCommand())
	rootCmd.AddCommand(NewCompletionCommand())
	rootCmd.AddCommand(NewProfilesCommand())
	rootCmd.AddCommand(NewStatsCommand())

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "profile")

	return rootCmd
}

func initConfig() {
	if cfgFile != "" {
		viper.Set("cfgFile", cfgFile)
	} else {
		viper.Set("cfgFile", constants.BaseFilePath)
	}

	viper.Set("profile", profile)

	viper.AutomaticEnv()
}
