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

package builder

import (
	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/escli/internal/util"
)

type Example struct {
	Description string
	DefinedOn   []string
}

var ExampleList = []Example{
	{
		Description: "  escli restore snapshot application-log snapshot-2020.01.01 application-log-2020.01.01",
		DefinedOn:   []string{"snapshot restore"},
	},
	{
		Description: "  escli snapshot delete application-log snapshot-2020.01.01",
		DefinedOn:   []string{"snapshot delete"},
	},
	{
		Description: "  escli snapshot create application-log snapshot-2020.01.01 application-log-2020.01.01,access-log-2020.01.01",
		DefinedOn:   []string{"snapshot create [repositoryID] [snapshotID] [indices]"},
	},
	{
		Description: "  escli snapshot archive application-log snapshot-2020.01.01",
		DefinedOn:   []string{"snapshot archive [repositoryID] [snapshotID]"},
	},
}

func SetCommandExample(cmd *cobra.Command) {
	for i := range ExampleList {
		ex := &ExampleList[i]

		if util.IsStringInArray(util.GetFullCommandUse(cmd), ex.DefinedOn) {
			cmd.Example = ex.Description
		}
	}
}
