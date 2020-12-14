package cmd

import (
	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/escli/cmd/cmd/cat"
)

func NewCatCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cat",
		Short: "Cat api for elasticsearch cluster",
		Long:  "Cat api for elasticsearch cluster",
	}

	cmd.AddCommand(cat.NewHealthCommand())

	return cmd
}
