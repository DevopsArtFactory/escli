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
