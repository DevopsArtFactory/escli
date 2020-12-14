package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/DevopsArtFactory/escli/internal/schema"
)

func ConvertToMap(r io.Reader, d *map[string]interface{}) {
	decoder := json.NewDecoder(r)
	for {
		if err := decoder.Decode(d); err == io.EOF {
			break
		}
	}
}

func ConvertToHealthMetadata(r io.Reader, d *[]schema.HealthMetadata) {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(d)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func ConvertToRepositorySnapshot(r io.Reader, d *[]schema.RepositorySnapshot) {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(d)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func ConvertToSnapshotMetadata(r io.Reader, d *schema.SnapshotMetadata) {
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
