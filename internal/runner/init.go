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
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"

	"github.com/DevopsArtFactory/escli/internal/config"
	"github.com/DevopsArtFactory/escli/internal/constants"
	"github.com/DevopsArtFactory/escli/internal/util"
)

func (r Runner) InitConfiguration() error {
	if util.FileExists(constants.BaseFilePath) {
		return fmt.Errorf("you already had configuration file: %s", constants.BaseFilePath)
	}

	// check base AWS directory
	if !util.FileExists(constants.ConfigDirectoryPath) {
		if err := os.MkdirAll(constants.ConfigDirectoryPath, 0755); err != nil {
			return err
		}
	}

	// Ask base account name which should be a company mail
	elasticsearchURL, err := AskElasticSearchURL()
	if err != nil {
		return err
	}

	awsRegion, err := AskAWSRegion()
	if err != nil {
		return err
	}

	c := config.SetInitConfig(elasticsearchURL, awsRegion)
	y, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// show configuration file for double-check
	if err := generateConfigFile(y); err != nil {
		return err
	}

	// ask to continue
	if err := util.AskContinue("Are you sure to generate configuration file? "); err != nil {
		return errors.New("initialization has been canceled")
	}

	if err := util.CreateFile(constants.BaseFilePath, string(y)); err != nil {
		return err
	}

	color.Blue("New configuration file is successfully generated in %s", constants.BaseFilePath)

	return nil
}

// AskBaseAccountName asks user's base account
func AskElasticSearchURL() (string, error) {
	var elasticsearchURL string
	prompt := &survey.Input{
		Message: "Your ElasticSearch URL : ",
	}
	survey.AskOne(prompt, &elasticsearchURL)

	if len(elasticsearchURL) == 0 {
		return elasticsearchURL, errors.New("input elasticsearch url has been canceled")
	}

	return elasticsearchURL, nil
}

func AskAWSRegion() (string, error) {
	var region string
	prompt := &survey.Input{
		Message: "Your AWS Default Region : ",
	}
	survey.AskOne(prompt, &region)

	if len(region) == 0 {
		return region, errors.New("input aws default region has been canceled")
	}

	return region, nil
}

func generateConfigFile(b []byte) error {
	_, err := fmt.Println(string(b))
	return err
}
