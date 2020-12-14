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
