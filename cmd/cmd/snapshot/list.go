package snapshot

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/escli/cmd/cmd/builder"
	"github.com/DevopsArtFactory/escli/internal/executor"
)

func NewListCommand() *cobra.Command {
	return builder.NewCmd("list").
		WithDescription("Listing snapshots").
		NoArgs(funcListSnapshot)
}

func funcListSnapshot(ctx context.Context, out io.Writer) error {
	return executor.RunExecutor(ctx, func(executor executor.Executor) error {
		return executor.Runner.ListSnapshot(out)
	})
}
