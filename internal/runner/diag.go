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
	"strconv"
	"strings"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"

	catSchema "github.com/DevopsArtFactory/escli/internal/schema/cat"
	"github.com/DevopsArtFactory/escli/internal/util"
)

func (r Runner) DiagCluster(out io.Writer) error {
	healthMetadata, err := r.Client.CatHealth()
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "check cluster status...........................")
	printStatusByHealth(out, healthMetadata[0])

	indexMetadata, err := r.Client.CatIndices(r.Flag.SortBy)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "check yellow status indices....................")
	printIndicesByHealth(out, indexMetadata, "yellow")

	fmt.Fprintf(out, "check red status indices.......................")
	printIndicesByHealth(out, indexMetadata, "red")

	// Check role of nodes
	fmt.Fprintf(out, "check number of master nodes...................")

	nodeMetadata, err := r.Client.CatNodes(r.Flag.SortBy)
	if err != nil {
		return err
	}

	checkNumberOfMasterNodes(out, nodeMetadata)

	// Check role of nodes
	fmt.Fprintf(out, "check maximum disk used percent of nodes.......")

	checkDiskUsedPercentOfNodes(out, nodeMetadata)

	return nil
}

func printStatusByHealth(out io.Writer, health catSchema.Health) {
	if health.Status != "green" {
		fmt.Fprintf(out, "[%s] %v\n", util.StringWithColor(health.Status), emoji.FaceScreamingInFear)
	} else {
		fmt.Fprintf(out, "[%s] %v\n", util.StringWithColor(health.Status), emoji.SmilingFaceWithSunglasses)
	}
}

func checkDiskUsedPercentOfNodes(out io.Writer, nodeMetadata []catSchema.Node) {
	maximumDiskUsedPercent := 0.0
	for _, v := range nodeMetadata {
		diskUsedPercent, _ := strconv.ParseFloat(v.DiskUsedPercent, 64)
		if diskUsedPercent >= maximumDiskUsedPercent {
			maximumDiskUsedPercent = diskUsedPercent
		}
	}
	fmt.Fprintf(out, "[%s]\n", util.FloatWithColor(maximumDiskUsedPercent))
}

func checkNumberOfMasterNodes(out io.Writer, nodeMetadata []catSchema.Node) {
	numberOfMasterNodes := 0
	for _, v := range nodeMetadata {
		if strings.Contains(v.NodeRole, "m") {
			numberOfMasterNodes++
		}
	}

	if util.IsEvenNumber(numberOfMasterNodes) {
		fmt.Fprintf(out, "[%s] %v\n", util.IntWithColor(numberOfMasterNodes, "red"), emoji.FaceScreamingInFear)
		fmt.Fprintf(out, "%v check more information by %s\n", emoji.ExclamationMark, util.StringWithColor("escli cat master"))
		color.Red("It will be caused split brain\n")
		color.Red("Number of master nodes will be odd number and minimum_master_nodes will be half number of master nodes + 1\n")
	} else {
		fmt.Fprintf(out, "[%s]\n", util.IntWithColor(numberOfMasterNodes, "green"))
	}
}

func printIndicesByHealth(out io.Writer, indexMetadata []catSchema.Index, health string) {
	numberOfTroubledIndices := 0
	for _, v := range indexMetadata {
		if v.Health == health {
			numberOfTroubledIndices++
		}
	}
	if numberOfTroubledIndices > 0 {
		fmt.Fprintf(out, "[%s] %v\n", util.IntWithColor(numberOfTroubledIndices, "red"), emoji.FaceScreamingInFear)
		fmt.Fprintf(out, "%v check more information by %s\n", emoji.ExclamationMark, util.StringWithColor("escli cat indices"))
	} else {
		fmt.Fprintf(out, "[%s] %v\n", util.IntWithColor(numberOfTroubledIndices, "green"), emoji.SmilingFaceWithSunglasses)
	}
}
