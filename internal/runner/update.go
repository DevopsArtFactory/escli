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
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"

	"github.com/DevopsArtFactory/escli/internal/util"
	"github.com/DevopsArtFactory/escli/internal/version"
)

func (r Runner) Update(out io.Writer) error {
	v := semver.MustParse(version.Get().Version)

	latest, found, err := selfupdate.DetectLatest("DevopsArtFactory/escli")
	if err != nil {
		return err
	}

	if !found || latest.Version.LTE(v) {
		fmt.Fprintf(out, "Current binary is the latest version")
		return nil
	}

	if err := util.AskContinue("Do you want to update to " + latest.Version.String()); err != nil {
		return errors.New("task has been canceled")
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		return err
	}

	fmt.Fprintf(out, "Successfully updated to version "+latest.Version.String())

	return nil
}
