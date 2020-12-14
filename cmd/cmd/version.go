package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/escli/cmd/cmd/builder"
	"github.com/DevopsArtFactory/escli/internal/version"
)

func NewVersionCommand() *cobra.Command {
	return builder.NewCmd("version").
		WithDescription("Print the version information").
		NoArgs(funcVersion)
}

// funcVersion
func funcVersion(_ context.Context, _ io.Writer) error {
	return version.Controller{}.Print(version.Get())
}
