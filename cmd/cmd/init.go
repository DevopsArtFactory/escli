package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/escli/cmd/cmd/builder"
	"github.com/DevopsArtFactory/escli/internal/executor"
)

func NewInitCommand() *cobra.Command {
	return builder.NewCmd("init").
		WithDescription("Initialize escli configuration").
		NoArgs(funcInitCommand)
}

func funcInitCommand(ctx context.Context, out io.Writer) error {
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		return executor.Runner.InitConfiguration()
	})
}
