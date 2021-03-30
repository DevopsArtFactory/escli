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

package config

import (
	"io/ioutil"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/DevopsArtFactory/escli/internal/schema"
)

func GetConfig() (*schema.Config, error) {
	var configs []schema.Config

	yamlFile, err := ioutil.ReadFile(viper.GetString("cfgFile"))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		return nil, err
	}

	profile := viper.Get("profile")

	for _, config := range configs {
		if config.Profile == profile {
			return &config, nil
		}
	}

	return &configs[0], nil
}

func GetDefaultConfig() (*schema.Config, error) {
	profile := &schema.Config{}

	return profile, nil
}

func SetInitConfig(profile string, elasticsearchURL string, awsRegion string) []schema.Config {
	config := schema.Config{
		Profile:          profile,
		ElasticSearchURL: elasticsearchURL,
		AWSRegion:        awsRegion,
	}

	return []schema.Config{config}
}
