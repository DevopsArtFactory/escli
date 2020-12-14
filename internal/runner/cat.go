package runner

import (
	"fmt"
	"io"
	"reflect"

	"github.com/fatih/color"
)

func (r Runner) CatHealth(out io.Writer) error {
	healthMetadata, err := r.Client.CatHealth()
	if err != nil {
		return err
	}

	printHealthMetadata(&healthMetadata[0])

	return nil
}

func printHealthMetadata(metadata interface{}) {
	e := reflect.ValueOf(metadata).Elem()
	filedNum := e.NumField()

	for i := 0; i < filedNum; i++ {
		v := e.Field(i)
		t := e.Type().Field(i)

		switch t.Name {
		case "Status":
			printHealthWithColor(fmt.Sprintf("%s", v.Interface()))
		case "ActiveShardsPercent":
			printActiveShardsPercentWithColor(fmt.Sprintf("%s", v.Interface()))
		default:
			fmt.Printf("%s : %s\n", t.Name, v.Interface())
		}
	}
}

func printHealthWithColor(health string) {
	fmt.Print("Status : ")
	switch health {
	case "green":
		color.Green(health)
	case "yellow":
		color.Yellow(health)
	case "red":
		color.Red(health)
	}
}

func printActiveShardsPercentWithColor(percent string) {
	fmt.Print("ActiveShardsPercent : ")
	switch percent {
	case "100.0%":
		color.Green(percent)
	default:
		color.Red(percent)
	}
}
