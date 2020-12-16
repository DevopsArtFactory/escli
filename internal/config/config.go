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
	"github.com/spf13/viper"

	"github.com/DevopsArtFactory/escli/internal/schema"
)

func GetConfig() (*schema.Config, error) {
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	conf := &schema.Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func GetDefaultConfig() (*schema.Config, error) {
	conf := &schema.Config{}

	return conf, nil
}

func SetInitConfig(elasticsearchURL string, awsRegion string) schema.Config {
	config := schema.Config{
		ElasticSearchURL: elasticsearchURL,
		AWSRegion:        awsRegion,
	}

	return config
}
