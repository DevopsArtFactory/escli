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

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/escli/cmd/cmd/builder"
	"github.com/DevopsArtFactory/escli/cmd/cmd/index"
)

func NewIndexCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "get or set settings for index",
		Long:  "get or set settings for index",
	}

	indexSettingsCommand := index.NewIndexSettingsCommand()
	indexDeleteCommand := index.NewIndexDeleteCommand()

	cmd.AddCommand(indexSettingsCommand)
	cmd.AddCommand(indexDeleteCommand)

	builder.SetCommandFlags(indexSettingsCommand)
	builder.SetCommandFlags(indexDeleteCommand)

	return cmd
}
