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
	"io"

	"github.com/fatih/color"

	"github.com/DevopsArtFactory/escli/internal/config"
)

// Old Version: "0.0.3"
// Current Version: "0.0.4"
func (r Runner) Fix(out io.Writer) error {
	p, err := config.GetOldConfig()
	if err != nil || p == nil {
		return err
	}

	if err := config.ConvertToNewConfig(p); err != nil {
		return err
	}

	color.Green("Successfully convert configuration file")

	return nil
}
