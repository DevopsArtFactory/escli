package executor

import (
	"context"

	"github.com/DevopsArtFactory/escli/internal/builder"
	"github.com/DevopsArtFactory/escli/internal/config"
	"github.com/DevopsArtFactory/escli/internal/runner"
	"github.com/DevopsArtFactory/escli/internal/schema"
)

type Executor struct {
	Runner  runner.Runner
	Context context.Context
}

var NewExecutor = createNewExecutor

func RunExecutor(ctx context.Context, action func(Executor) error) error {
	c, err := config.GetConfig()
	if err != nil {
		return err
	}

	executor, _ := createNewExecutor(c)

	err = action(*executor)

	return alwaysSucceedWhenCancelled(ctx, err)
}

func RunExecutorWithoutCheckingConfig(ctx context.Context, action func(Executor) error) error {
	c, err := config.GetDefaultConfig()
	if err != nil {
		return err
	}

	executor, _ := createNewExecutor(c)

	err = action(*executor)

	return alwaysSucceedWhenCancelled(ctx, err)
}

func createNewExecutor(config *schema.Config) (*Executor, error) {
	flags, err := builder.ParseFlags()

	if err != nil {
		return nil, err
	}

	executor := Executor{
		Context: context.Background(),
		Runner:  runner.New(flags, config),
	}
	return &executor, nil
}

func alwaysSucceedWhenCancelled(ctx context.Context, err error) error {
	// if the context was cancelled act as if all is well
	if err != nil && ctx.Err() == context.Canceled {
		return nil
	}
	return err
}
