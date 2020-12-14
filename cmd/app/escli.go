package app

import (
	"context"
	"io"

	"github.com/DevopsArtFactory/escli/cmd/cmd"
)

func Run(out, stderr io.Writer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := cmd.NewRootCommand(out, stderr)
	return c.ExecuteContext(ctx)
}
