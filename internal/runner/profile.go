package runner

import (
	"fmt"
	"github.com/DevopsArtFactory/escli/internal/schema"
	"github.com/DevopsArtFactory/escli/internal/util"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
)

func (r Runner) ListProfile(out io.Writer) error {
	var configs []schema.Config

	yamlFile, err := ioutil.ReadFile(viper.GetString("cfgFile"))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &configs)

	for _, config := range configs {
		fmt.Fprintf(out, "Profile           : %s\n", util.StringWithColor(config.Profile))
		fmt.Fprintf(out, "ElasticSearch URL : %s\n", util.StringWithColor(config.ElasticSearchURL))
		fmt.Fprintf(out, "AWS Region        : %s\n\n", util.StringWithColor(config.AWSRegion))
	}

	return nil
}
