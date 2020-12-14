package cmd

import (
	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/escli/cmd/cmd/snapshot"
)

func NewSnapshotCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage snapshot in elasticsearch cluster",
		Long:  "Manage snapshot in elasticsearch cluster",
	}

	cmd.AddCommand(snapshot.NewListCommand())
	cmd.AddCommand(snapshot.NewArchiveCommand())
	cmd.AddCommand(snapshot.NewRestoreCommand())

	return cmd
}
