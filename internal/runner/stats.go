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
	"time"

	"github.com/DevopsArtFactory/escli/internal/constants"
	indexSchema "github.com/DevopsArtFactory/escli/internal/schema/index"
	"github.com/DevopsArtFactory/escli/internal/util"
)

func (r Runner) Stats(out io.Writer, args []string) error {
	var statsMetadata *indexSchema.Stats
	var err error

	interval, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%-10s\t%20s\t%20s\t%20s\t"+
		"%20s\t%20s\t"+
		"%20s\t%20s\t"+
		"%20s\t%20s\n",
		"time", "total shards", "successful shards", "failed shards",
		"indexing rate", "indexing latency (ms)",
		"query rate", "query latency (ms)",
		"fetch rate", "fetch latency (ms)")

	var prevStatsMetadata *indexSchema.Stats

	for {
		if r.Config.Product == constants.OpenSearch {
			statsMetadata, err = r.Client.OSStats()
		} else {
			statsMetadata, err = r.Client.Stats()
		}

		if err != nil {
			break
		}

		if prevStatsMetadata != nil {
			now := time.Now().Local()
			timestamp := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())

			fmt.Fprintf(out, "%-10s\t%20d\t%20d\t%20d\t"+
				"%20.0f\t%20.2f\t"+
				"%20.0f\t%20.2f\t"+
				"%20.0f\t%20.2f\n",
				timestamp, statsMetadata.Shards.Total, statsMetadata.Shards.Successful, statsMetadata.Shards.Failed,
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
