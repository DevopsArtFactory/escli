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

	"github.com/DevopsArtFactory/escli/internal/util"
)

func (r Runner) CatHealth(out io.Writer) error {
	healthMetadata, err := r.Client.CatHealth()
	if err != nil {
		return err
	}

	e := reflect.ValueOf(&healthMetadata[0]).Elem()
	filedNum := e.NumField()

	for i := 0; i < filedNum; i++ {
		v := e.Field(i)
		t := e.Type().Field(i)

		fmt.Printf("%-20s : %10s\n", t.Name,
			util.StringWithColor(fmt.Sprintf("%s", v.Interface())),
		)
	}

	return nil
}

func (r Runner) CatIndices(out io.Writer) error {
	indicesMetadata, err := r.Client.CatIndices(r.Flag.SortBy)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%-50s\t%s\t%s\t%s\t%s\t%10s\n",
		"index",
		"health",
		"status",
		"pri",
		"rep",
		"store.size")
	for _, index := range indicesMetadata {
		if r.Flag.TroubledOnly && index.Health == "green" {
			continue
		}
		fmt.Fprintf(out, "%-50s\t%s\t%s\t%s\t%s\t%10s\n",
			index.Index,
			util.StringWithColor(index.Health),
			index.Status,
			index.PrimaryShards,
			index.ReplicaShards,
			index.StoreSize)
	}

	return nil
}

func (r Runner) CatNodes(out io.Writer) error {
	nodesMetadata, err := r.Client.CatNodes(r.Flag.SortBy)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%-50s\t%s\t%4s\t%7s\t%7s\t%8s\t%6s\t%17s\n",
		"name",
		"ip",
		"role",
		"load_1m",
		"load_5m",
		"load_15m",
		"uptime",
		"disk.used_percent")
	for _, node := range nodesMetadata {
		fmt.Fprintf(out, "%-50s\t%s\t%4s\t%7s\t%7s\t%8s\t%6s\t%17s\n",
			node.Name,
			node.IP,
			node.NodeRole,
			node.Load1M,
			node.Load5M,
			node.Load15M,
			node.Uptime,
			node.DiskUsedPercent)
	}

	return nil
}

func (r Runner) CatShards(out io.Writer) error {
	shardsMetadata, err := r.Client.CatShards(r.Flag.SortBy)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%-50s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		"index",
		"shard",
		"prirep",
		"state",
		"docs",
		"store",
		"ip",
		"node",
		"unassigned.reason")
	for _, shard := range shardsMetadata {
		if r.Flag.TroubledOnly && shard.State == "STARTED" {
			continue
		}
		fmt.Fprintf(out, "%-50s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			shard.Index,
			shard.Shard,
			shard.PriRep,
			shard.State,
			shard.Docs,
			shard.Store,
			shard.IP,
			shard.Node,
			shard.UnassignedReason)
	}

	return nil
}
