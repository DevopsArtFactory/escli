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
	"github.com/DevopsArtFactory/escli/internal/constants"
	"github.com/DevopsArtFactory/escli/internal/schema"
	"io"
)

func (r Runner) IndexSettings(out io.Writer, args []string) error {
	switch len(args) {
	case constants.GetSetting:
		return r.Client.GetIndexSetting(args[0], "")
	case constants.GetSettingWithName:
		return r.Client.GetIndexSetting(args[0], args[1])
	case constants.PutSetting:
		requestBody, _ := json.Marshal(
			schema.IndexIndexSettingsRequestBody{
				Index: map[string]string{
					args[1]: args[2],
				},
			},
		)
		return r.Client.PutIndexSetting(args[0], string(requestBody))
	default:
		errors.New("arguments must be 2 or 3")
	}

	return nil
}