package constants

import "os"

const (
	DefaultRegion = "ap-northeast-2"
	EmptyString   = ""
)

var (
	ConfigDirectoryPath = HomeDir() + "/.escli"
	BaseFilePath        = ConfigDirectoryPath + "/config.yaml"
)

// Get Home Directory
func HomeDir() string {
	if h := os.Getenv("HOME"); h != EmptyString {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
