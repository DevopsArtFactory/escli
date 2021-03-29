package cmd

import (
	"github.com/DevopsArtFactory/escli/cmd/cmd/profiles"
	"github.com/spf13/cobra"
)

func NewProfilesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "profiles",
		Short: "manage profiles for escli",
		Long: "manage profiles for escli",
	}

	profilesListCommand := profiles.NewProfileListCommand()

	cmd.AddCommand(profilesListCommand)

	return cmd
}