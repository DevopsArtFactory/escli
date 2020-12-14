package version

import (
	"fmt"
	"runtime"
)

var version, gitCommit, gitTreeState, buildDate string
var platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

type Controller struct{}

type Info struct {
	Version      string
	BuildDate    string
	GitCommit    string
	GitTreeState string
	Platform     string
}

func Get() Info {
	return Info{
		Version:      version,
		BuildDate:    buildDate,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		Platform:     platform,
	}
}

func (v Controller) Print(info Info) error {
	_, err := fmt.Println(info.Version)
	return err
}
