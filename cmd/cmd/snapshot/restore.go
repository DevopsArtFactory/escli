package snapshot

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/escli/cmd/cmd/builder"
	"github.com/DevopsArtFactory/escli/internal/executor"
)

func NewRestoreCommand() *cobra.Command {
	return builder.NewCmd("restore").
		WithDescription("Restore snapshot from S3 glacier to standard").
		SetFlags().
		SetExample().
		ExactArgs(3, funcRestoreSnapshot)
}

func funcRestoreSnapshot(ctx context.Context, out io.Writer, args []string) error {
	return executor.RunExecutor(ctx, func(executor executor.Executor) error {
		return executor.Runner.RestoreSnapshot(out, args)
	})
}
