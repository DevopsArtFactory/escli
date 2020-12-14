package snapshot

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/escli/cmd/cmd/builder"
	"github.com/DevopsArtFactory/escli/internal/executor"
)

func NewArchiveCommand() *cobra.Command {
	return builder.NewCmd("archive").
		WithDescription("Archive snapshot to S3 glacier").
		SetFlags().
		ExactArgs(2, funcArchiveSnapshot)
}

func funcArchiveSnapshot(ctx context.Context, out io.Writer, args []string) error {
	return executor.RunExecutor(ctx, func(executor executor.Executor) error {
		return executor.Runner.ArchiveSnapshot(out, args)
	})
}
