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

package executor

import (
	"fmt"
	"strings"

	"github.com/DevopsArtFactory/escli/internal/builder"
	"github.com/DevopsArtFactory/escli/internal/constants"
	"github.com/DevopsArtFactory/escli/internal/util"
)

// FlagValidation check flag value validation
func FlagValidation(flags builder.Flags) error {
	if len(flags.RestoreTier) > 0 && !util.IsStringInArray(flags.RestoreTier, constants.ValidRestoreTier) {
		return fmt.Errorf("restore-tier should be one of [%s]", strings.Join(constants.ValidRestoreTier, ","))
	}

	if flags.MaxConcurrentJob > 0 && flags.MaxConcurrentJob > constants.HardLimitMaxConcurrentJob || flags.MaxConcurrentJob < 0 {
		return fmt.Errorf("max-concurrent-job cannot over %d or less than 1", constants.HardLimitMaxConcurrentJob)
	}

	return nil
}
