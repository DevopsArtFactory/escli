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
		DefinedOn:   []string{"restore snapshot"},
	},
}

func SetCommandExample(cmd *cobra.Command) {
	for i := range ExampleList {
		ex := &ExampleList[i]

		if util.IsStringInArray(cmd.Short, ex.DefinedOn) {
			cmd.Example = ex.Description
		}
	}
}
