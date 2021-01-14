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

package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func ConvertJSONtoMetadata(r io.Reader, d interface{}) {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(d)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func ParseInt(numericString string) int {
	num, _ := strconv.Atoi(numericString)

	return num
}

func IsStringInArray(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}

	return false
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// AskContinue provides interactive terminal for users to answer if they continue process or not
func AskContinue(msg string) error {
	var answer string
	prompt := &survey.Input{
		Message: msg,
	}
	survey.AskOne(prompt, &answer)

	if IsStringInArray(strings.ToLower(answer), []string{"yes", "y"}) {
		return nil
	}

	return errors.New("stop process")
}

//CreateFile creates/overrides files
func CreateFile(filePath string, writeData string) error {
	if err := ioutil.WriteFile(filePath, []byte(writeData), 0644); err != nil {
		return err
	}
	return nil
}

func GreenString(content string) string {
	return color.GreenString(content)
}

func RedString(content string) string {
	return color.RedString(content)
}

func YellowString(content string) string {
	return color.YellowString(content)
}

func StringWithColor(content string) string {
	switch content {
	case "green":
		return color.GreenString(content)
	case "yellow":
		return color.YellowString(content)
	case "red":
		return color.RedString(content)
	case "FAILED":
		return color.RedString(content)
	case "IN_PROGRESS":
		return color.YellowString(content)
	case "PARTIAL":
		return color.RedString(content)
	case "SUCCESS":
		return color.GreenString(content)
	case "100.0%":
		return color.GreenString(content)
	}

	if strings.Contains(content, "%") {
		return color.RedString(content)
	}

	return color.BlueString(content)
}

func IntWithColor(number int, status string) string {
	switch status {
	case "green":
		return color.GreenString("%d", number)
	case "yellow":
		return color.YellowString("%d", number)
	case "red":
		return color.RedString("%d", number)
	}

	return color.BlueString("%d", number)
}

func FloatWithColor(number float64) string {
	switch {
	case number > 90:
		return color.RedString("%.0f", number)
	case number > 80:
		return color.YellowString("%.0f", number)
	case number > 70:
		return color.BlueString("%.0f", number)
	default:
		return color.GreenString("%.0f", number)
	}
}

func IsEvenNumber(number int) bool {
	return number%2 == 0
}

func GetFullCommandUse(cmd *cobra.Command) string {
	if cmd.Parent() != nil {
		return fmt.Sprintf("%s %s", GetFullCommandUse(cmd.Parent()), cmd.Use)
	}
	return cmd.Use
}

func ReturnErrorFromResponseBody(response *esapi.Response) error {
	switch response.StatusCode {
	case 200:
		return nil
	default:
		return errors.New(responseBodyToString(response.Body))
	}
}

func responseBodyToString(closer io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(closer)
	return buf.String()
}

func JSONtoPrettyString(v interface{}) (string, error) {
	jsonPrettyString, err := json.MarshalIndent(v, "", "\t")
	return string(jsonPrettyString), err
}
