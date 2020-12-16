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
	"github.com/DevopsArtFactory/escli/internal/builder"
	"github.com/DevopsArtFactory/escli/internal/client"
	"github.com/DevopsArtFactory/escli/internal/constants"
	"github.com/DevopsArtFactory/escli/internal/schema"
)

type Runner struct {
	Client client.Client
	Flag   *builder.Flags
	Config *schema.Config
}

func New(flags *builder.Flags, config *schema.Config) Runner {
	region := flags.Region
	if len(region) == 0 {
		region = constants.DefaultRegion
	}

	return Runner{
		Client: client.NewClient(client.GetAwsSession(region), nil, config),
		Flag:   flags,
		Config: config,
	}
}
