package profiles

import (
	"context"
	"github.com/DevopsArtFactory/escli/cmd/cmd/builder"
	"github.com/DevopsArtFactory/escli/internal/executor"
	"github.com/spf13/cobra"
	"io"
)

func NewProfileListCommand() *cobra.Command {
	return builder.NewCmd("list").
		WithDescription("list profiles").
		NoArgs(funcProfileList)
}

func funcProfileList(ctx context.Context, out io.Writer) error {
	return executor.RunExecutor(ctx, func(executor executor.Executor) error {
		return executor.Runner.ListProfile(out)
	})
}
