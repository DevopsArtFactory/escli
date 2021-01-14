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

package runner

import (
	"errors"
	"fmt"
	"io"

	"github.com/DevopsArtFactory/escli/internal/constants"
	clusterSchema "github.com/DevopsArtFactory/escli/internal/schema/cluster"
	"github.com/DevopsArtFactory/escli/internal/util"
)

func (r Runner) ClusterReroute(out io.Writer) error {
	resp, err := r.Client.ClusterReroute()

	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%s\n", resp)
	return err
}

func (r Runner) ClusterSettings(out io.Writer, args []string) error {
	var resp string
	var err error

	switch len(args) {
	case constants.GetClusterSetting:
		resp, err = r.Client.GetClusterSetting()
	case constants.PutClusterSetting:
		requestBody, _ := util.JSONtoPrettyString(composeClusterRequestBody(args))
		fmt.Fprintf(out, "%s\n", util.YellowString(requestBody))

		if !r.Flag.Force {
			if err := util.AskContinue("Are you sure to update settings of cluster"); err != nil {
				return errors.New("task has benn canceled")
			}
		}

		resp, err = r.Client.PutClusterSetting(requestBody)
	default:
		return errors.New("arguments must be 0 or 3")
	}

	fmt.Fprintf(out, "%s\n", resp)
	return err
}

func composeClusterRequestBody(args []string) clusterSchema.RequestBody {
	switch args[0] {
	case "persistent":
		return clusterSchema.RequestBody{
			Persistent: map[string]string{
				args[1]: args[2],
			},
		}
	case "transient":
		return clusterSchema.RequestBody{
			Transient: map[string]string{
				args[1]: args[2],
			},
		}
	default:
		return clusterSchema.RequestBody{}
	}
}
