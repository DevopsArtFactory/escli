/*
copyright 2020 the Escli authors

licensed under the apache license, version 2.0 (the "license");
you may not use this file except in compliance with the license.
you may obtain a copy of the license at

    http://www.apache.org/licenses/license-2.0

unless required by applicable law or agreed to in writing, software
distributed under the license is distributed on an "as is" basis,
without warranties or conditions of any kind, either express or implied.
see the license for the specific language governing permissions and
limitations under the license.
*/

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
