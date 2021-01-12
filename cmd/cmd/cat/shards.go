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

package cat

import (
	"context"
	"github.com/DevopsArtFactory/escli/internal/executor"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/escli/cmd/cmd/builder"
)

func NewCatShardsCommand() *cobra.Command {
	return builder.NewCmd("shards").
		WithDescription("_cat/shards").
		NoArgs(funcCatShards)
}

func funcCatShards(ctx context.Context, out io.Writer) error {
	return executor.RunExecutor(ctx, func(executor executor.Executor) error {
		return executor.Runner.CatShards(out)
	})
}