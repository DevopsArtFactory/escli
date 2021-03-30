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
	"io/ioutil"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/DevopsArtFactory/escli/internal/config"
	"github.com/DevopsArtFactory/escli/internal/schema"
	"github.com/DevopsArtFactory/escli/internal/util"
)

func (r Runner) ListProfile(out io.Writer) error {
	var configs []schema.Config

	yamlFile, err := ioutil.ReadFile(viper.GetString("cfgFile"))
	if err != nil {
		return err
	}

	_ = yaml.Unmarshal(yamlFile, &configs)

	for _, config := range configs {
		fmt.Fprintf(out, "Profile           : %s\n", util.StringWithColor(config.Profile))
		fmt.Fprintf(out, "ElasticSearch URL : %s\n", util.StringWithColor(config.ElasticSearchURL))
		fmt.Fprintf(out, "AWS Region        : %s\n\n", util.StringWithColor(config.AWSRegion))
	}

	return nil
}

func (r Runner) RemoveProfile(out io.Writer, args []string) error {
	var configs []schema.Config

	profileName := args[0]
	cfgFile := viper.GetString("cfgFile")

	yamlFile, err := ioutil.ReadFile(viper.GetString("cfgFile"))
	if err != nil {
		return err
	}

	_ = yaml.Unmarshal(yamlFile, &configs)

	for k, config := range configs {
		if config.Profile == profileName {
			if err := util.AskContinue("Are you sure to remove profile " + util.RedString(profileName) + " from configuration file? "); err != nil {
				return errors.New("removing profile has been canceled")
			}
			configs = util.RemoveSlice(configs, k)
		}
	}

	y, err := yaml.Marshal(configs)
	if err != nil {
		return err
	}

	if err := util.CreateFile(cfgFile, string(y)); err != nil {
		return err
	}

	if len(configs) == 0 {
		os.Create(cfgFile)
	}

	fmt.Fprintf(out, "Remove profile %s from configuration file is successfully in %s", util.BlueString(profileName), util.BlueString(cfgFile))

	return nil
}

func (r Runner) AddProfile(out io.Writer) error {
	cfgFile := viper.GetString("cfgFile")

	if !util.FileExists(cfgFile) {
		_, err := os.Create(cfgFile)
		if err != nil {
			return err
		}
	}

	profile, err := AskProfile()
	if err != nil {
		return err
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

	c := config.SetInitConfig(profile, elasticsearchURL, awsRegion)
	y, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// show configuration file for double-check
	if err := generateConfigFile(y); err != nil {
		return err
	}

	// ask to continue
	if err := util.AskContinue("Are you sure to add profile to configuration file? "); err != nil {
		return errors.New("adding profile has been canceled")
	}

	if err := util.AppendFile(cfgFile, string(y)); err != nil {
		return err
	}

	fmt.Fprintf(out, "Adding profile to configuration file is successfully in %s", util.BlueString(cfgFile))

	return nil
}

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

func AskProfile() (string, error) {
	var profile string
	prompt := &survey.Input{
		Message: "Your Profile Name : ",
	}
	survey.AskOne(prompt, &profile)

	if len(profile) == 0 {
		return profile, errors.New("input profile name has been canceled")
	}

	return profile, nil
}

func AskAWSRegion() (string, error) {
	var region string
	prompt := &survey.Input{
		Message: "Your AWS Default Region (If you don't use AWS, type blank) : ",
	}
	survey.AskOne(prompt, &region)

	return region, nil
}

func generateConfigFile(b []byte) error {
	_, err := fmt.Println(string(b))
	return err
}
