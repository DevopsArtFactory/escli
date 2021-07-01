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
	"fmt"
	"io/ioutil"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/DevopsArtFactory/escli/internal/schema"
	"github.com/DevopsArtFactory/escli/internal/util"
)

func GetConfig() (*schema.Config, error) {
	var configs []schema.Config

	filePath := viper.GetString("cfgFile")
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		if old, err := detectOldConfig(); err == nil && old {
			return nil, fmt.Errorf("you are using old version configuration. please run `escli fix` to migrate configuration")
		}
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

func detectConfig() (bool, error) {
	var config []schema.Config
	filePath := viper.GetString("cfgFile")
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return false, err
	}

	return true, err
}

func detectOldConfig() (bool, error) {
	var config schema.OldConfig
	filePath := viper.GetString("cfgFile")
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return false, err
	}

	return true, err
}

func GetOldConfig() (*schema.OldConfig, error) {
	var config schema.OldConfig

	filePath := viper.GetString("cfgFile")
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		if new, err := detectConfig(); err == nil && new {
			color.Green("your configuration is already updated")
			return nil, nil
		}
		return nil, err
	}

	return &config, nil
}

func ConvertToNewConfig(oc *schema.OldConfig) error {
	cfgFile := viper.GetString("cfgFile")

	profile, err := util.AskProfile()
	if err != nil {
		return err
	}

	c := SetInitConfig(profile, oc.ElasticSearchURL, oc.AWSRegion)
	y, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	if err := util.CreateFile(cfgFile, string(y)); err != nil {
		return err
	}

	return nil
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
