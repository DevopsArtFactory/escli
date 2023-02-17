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

	"github.com/DevopsArtFactory/escli/internal/constants"
	catSchema "github.com/DevopsArtFactory/escli/internal/schema/cat"
	"github.com/DevopsArtFactory/escli/internal/util"
)

func (r Runner) DiagCluster(out io.Writer) error {
	var healthMetadata []catSchema.Health
	var err error

	if r.Config.Product == constants.OpenSearch {
		healthMetadata, err = r.Client.OSCatHealth()
	} else {
		healthMetadata, err = r.Client.CatHealth()
	}

	if err != nil {
		return err
	}

	fmt.Fprintf(out, "check cluster status...........................")
	printStatusByHealth(out, healthMetadata[0])

	var indexMetadata []catSchema.Index

	if r.Config.Product == constants.OpenSearch {
		indexMetadata, err = r.Client.OSCatIndices(r.Flag.SortBy)
	} else {
		indexMetadata, err = r.Client.CatIndices(r.Flag.SortBy)
	}

	if err != nil {
		return err
	}

	fmt.Fprintf(out, "check yellow status indices....................")
	printIndicesByHealth(out, indexMetadata, "yellow")

	fmt.Fprintf(out, "check red status indices.......................")
	printIndicesByHealth(out, indexMetadata, "red")

	// Check role of nodes
	fmt.Fprintf(out, "check number of master nodes...................")

	var nodeMetadata []catSchema.Node

	if r.Config.Product == constants.OpenSearch {
		nodeMetadata, err = r.Client.OSCatNodes(r.Flag.SortBy)
	} else {
		nodeMetadata, err = r.Client.CatNodes(r.Flag.SortBy)
	}

	if err != nil {
		return err
	}

	checkNumberOfMasterNodes(out, nodeMetadata)

	// Check disk usage of nodes
	fmt.Fprintf(out, "check minimum disk used percent od data node...")
	checkMinimumDiskUsedPercentOfNodes(out, nodeMetadata)

	fmt.Fprintf(out, "check maximum disk used percent of data node...")
	checkMaximumDiskUsedPercentOfNodes(out, nodeMetadata)

	return nil
}

func printStatusByHealth(out io.Writer, health catSchema.Health) {
	if health.Status != "green" {
		fmt.Fprintf(out, "[%s] %v\n", util.StringWithColor(health.Status), emoji.FaceScreamingInFear)
	} else {
		fmt.Fprintf(out, "[%s] %v\n", util.StringWithColor(health.Status), emoji.SmilingFaceWithSunglasses)
	}
}

func checkMinimumDiskUsedPercentOfNodes(out io.Writer, nodeMetadata []catSchema.Node) {
	minimumDiskUsedPercent := 100.0
	for _, v := range nodeMetadata {
		if strings.Contains(v.NodeRole, "d") {
			diskUsedPercent, _ := strconv.ParseFloat(v.DiskUsedPercent, 64)
			if diskUsedPercent <= minimumDiskUsedPercent {
				minimumDiskUsedPercent = diskUsedPercent
			}
		}
	}
	fmt.Fprintf(out, "[%s]\n", util.FloatWithColor(minimumDiskUsedPercent))
}

func checkMaximumDiskUsedPercentOfNodes(out io.Writer, nodeMetadata []catSchema.Node) {
	maximumDiskUsedPercent := 0.0
	for _, v := range nodeMetadata {
		if strings.Contains(v.NodeRole, "d") {
			diskUsedPercent, _ := strconv.ParseFloat(v.DiskUsedPercent, 64)
			if diskUsedPercent >= maximumDiskUsedPercent {
				maximumDiskUsedPercent = diskUsedPercent
			}
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

	if util.IsEvenNumber(numberOfMasterNodes) || numberOfMasterNodes < 2 {
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
