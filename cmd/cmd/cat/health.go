package cat

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/escli/cmd/cmd/builder"
	"github.com/DevopsArtFactory/escli/internal/executor"
)

func NewHealthCommand() *cobra.Command {
	return builder.NewCmd("health").
		WithDescription("_cat/health").
		SetFlags().
		NoArgs(funcCatHealth)
}

func funcCatHealth(ctx context.Context, out io.Writer) error {
	return executor.RunExecutor(ctx, func(executor executor.Executor) error {
		return executor.Runner.CatHealth(out)
	})
}
