package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/DevopsArtFactory/escli/internal/constants"
)

var (
	cfgFile string
)

func NewRootCommand(out, stderr io.Writer) *cobra.Command {
	cobra.OnInitialize(initConfig)
	rootCmd := &cobra.Command{
		Use:           "escli",
		Short:         "manage elasticsearch cluster",
		Long:          "manage elasticsearch cluster",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
		},
	}

	rootCmd.AddCommand(NewSnapshotCommand())
	rootCmd.AddCommand(NewCatCommand())
	rootCmd.AddCommand(NewInitCommand())
	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	return rootCmd
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile(constants.BaseFilePath)
	}

	viper.AutomaticEnv()
}
