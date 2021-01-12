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
	"github.com/DevopsArtFactory/escli/cmd/cmd/builder"
	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/escli/cmd/cmd/snapshot"
)

func NewSnapshotCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "manage snapshot in elasticsearch cluster",
		Long:  "manage snapshot in elasticsearch cluster",
	}

	snapshotListCommand := snapshot.NewSnapshotListCommand()
	snapshotArchiveCommand := snapshot.NewSnapshotArchiveCommand()
	snapshotRestoreCommand := snapshot.NewSnapshotRestoreCommand()
	snapshotCreateCommand := snapshot.NewSnapshotCreateCommand()
	snapshotDeleteCommand := snapshot.NewSnapshotDeleteCommand()

	cmd.AddCommand(snapshotListCommand)
	builder.SetCommandFlags(snapshotListCommand)

	cmd.AddCommand(snapshotArchiveCommand)
	builder.SetCommandExample(snapshotArchiveCommand)
	builder.SetCommandFlags(snapshotArchiveCommand)

	cmd.AddCommand(snapshotRestoreCommand)
	builder.SetCommandFlags(snapshotRestoreCommand)

	cmd.AddCommand(snapshotCreateCommand)
	builder.SetCommandExample(snapshotCreateCommand)

	cmd.AddCommand(snapshotDeleteCommand)
	builder.SetCommandExample(snapshotRestoreCommand)

	return cmd
}
