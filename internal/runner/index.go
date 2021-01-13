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
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/DevopsArtFactory/escli/internal/constants"
	indexSchema "github.com/DevopsArtFactory/escli/internal/schema/index"
)

func (r Runner) IndexSettings(out io.Writer, args []string) error {
	var resp string
	var err error

	switch len(args) {
	case constants.GetIndexSetting:
		resp, err = r.Client.GetIndexSetting(args[0], "")
	case constants.GetIndexSettingWithName:
		resp, err = r.Client.GetIndexSetting(args[0], args[1])
	case constants.PutIndexSetting:
		requestBody, _ := json.Marshal(composeIndexRequestBody(args))
		resp, err = r.Client.PutIndexSetting(args[0], string(requestBody))
	default:
		return errors.New("arguments must be 1 or 2 or 3")
	}

	fmt.Fprintf(out, "%s\n", resp)
	return err
}

func composeIndexRequestBody(args []string) indexSchema.RequestBody {
	return indexSchema.RequestBody{
		Index: map[string]string{
			args[1]: args[2],
		},
	}
}
