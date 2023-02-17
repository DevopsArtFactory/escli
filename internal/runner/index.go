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
	"strconv"
	"time"

	"github.com/DevopsArtFactory/escli/internal/constants"
	indexSchema "github.com/DevopsArtFactory/escli/internal/schema/index"
	"github.com/DevopsArtFactory/escli/internal/util"
)

func (r Runner) IndexSettings(out io.Writer, args []string) error {
	var resp string
	var err error

	switch len(args) {
	case constants.GetIndexSetting:
		if r.Config.Product == constants.OpenSearch {
			resp, err = r.Client.OSGetIndexSetting(args[0], "")
		} else {
			resp, err = r.Client.GetIndexSetting(args[0], "")
		}
	case constants.GetIndexSettingWithName:
		resp, err = r.Client.GetIndexSetting(args[0], args[1])
	case constants.PutIndexSetting:
		requestBody, _ := util.JSONtoPrettyString(composeIndexRequestBody(args))
		fmt.Fprintf(out, "%s\n", util.YellowString(requestBody))

		if !r.Flag.Force {
			if err := util.AskContinue("Are you sure to update settings of index"); err != nil {
				return errors.New("task has benn canceled")
			}
		}

		if r.Config.Product == constants.OpenSearch {
			resp, err = r.Client.OSPutIndexSetting(args[0], requestBody)
		} else {
			resp, err = r.Client.PutIndexSetting(args[0], requestBody)
		}
	default:
		return errors.New("arguments must be 1 or 2 or 3")
	}

	fmt.Fprintf(out, "%s\n", resp)
	return err
}

func (r Runner) CreateIndex(out io.Writer, args []string) error {
	var resp string
	var err error

	if r.Config.Product == constants.OpenSearch {
		resp, err = r.Client.OSCreateIndex(args[0])
	} else {
		resp, err = r.Client.CreateIndex(args[0])
	}

	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%s\n", resp)
	return err
}

func (r Runner) DeleteIndex(out io.Writer, args []string) error {
	var resp string
	var err error

	if !r.Flag.Force {
		if err := util.AskContinue("Are you sure to delete index"); err != nil {
			return errors.New("task has benn canceled")
		}
	}

	if r.Config.Product == constants.OpenSearch {
		resp, err = r.Client.OSDeleteIndex(args)
	} else {
		resp, err = r.Client.DeleteIndex(args)
	}

	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%s\n", resp)
	return err
}

func (r Runner) StatsIndex(out io.Writer, args []string) error {
	var statsMetadata *indexSchema.Stats
	var err error

	interval, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%-10s\t%-20s\t%20s\t%20s\t%20s\t"+
		"%20s\t%20s\t"+
		"%20s\t%20s\t"+
		"%20s\t%20s\n",
		"time", "index", "total shards", "successful shards", "failed shards",
		"indexing rate", "indexing latency (ms)",
		"query rate", "query latency (ms)",
		"fetch rate", "fetch latency (ms)")

	var prevStatsMetadata *indexSchema.Stats

	for {
		if r.Config.Product == constants.OpenSearch {
			statsMetadata, err = r.Client.OSIndexStats(args[0])
		} else {
			statsMetadata, err = r.Client.IndexStats(args[0])
		}

		if err != nil {
			break
		}

		if prevStatsMetadata != nil {
			now := time.Now().Local()
			timestamp := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())

			fmt.Fprintf(out, "%-10s\t%-20s\t%20d\t%20d\t%20d\t"+
				"%20.0f\t%20.2f\t"+
				"%20.0f\t%20.2f\t"+
				"%20.0f\t%20.2f\n",
				timestamp, args[0], statsMetadata.Shards.Total, statsMetadata.Shards.Successful, statsMetadata.Shards.Failed,
				util.Divide(statsMetadata.All.Total.Indexing.IndexTotal-prevStatsMetadata.All.Total.Indexing.IndexTotal, interval), util.Divide(statsMetadata.All.Total.Indexing.IndexTimeInMillis-prevStatsMetadata.All.Total.Indexing.IndexTimeInMillis, statsMetadata.All.Total.Indexing.IndexTotal-prevStatsMetadata.All.Total.Indexing.IndexTotal),
				util.Divide(statsMetadata.All.Total.Search.QueryTotal-prevStatsMetadata.All.Total.Search.QueryTotal, interval), util.Divide(statsMetadata.All.Total.Search.QueryTimeInMillis-prevStatsMetadata.All.Total.Search.QueryTimeInMillis, statsMetadata.All.Total.Search.QueryTotal-prevStatsMetadata.All.Total.Search.QueryTotal),
				util.Divide(statsMetadata.All.Total.Search.FetchTotal-prevStatsMetadata.All.Total.Search.FetchTotal, interval), util.Divide(statsMetadata.All.Total.Search.FetchTimeInMillis-prevStatsMetadata.All.Total.Search.FetchTimeInMillis, statsMetadata.All.Total.Search.FetchTotal-prevStatsMetadata.All.Total.Search.FetchTotal),
			)
		}

		prevStatsMetadata = statsMetadata

		time.Sleep(time.Second * time.Duration(interval))
	}

	return err
}

func composeIndexRequestBody(args []string) indexSchema.RequestBody {
	return indexSchema.RequestBody{
		Index: map[string]string{
			args[1]: args[2],
		},
	}
}
