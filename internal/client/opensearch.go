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

package client

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"

	"github.com/DevopsArtFactory/escli/internal/constants"
	catSchema "github.com/DevopsArtFactory/escli/internal/schema/cat"
	indexSchema "github.com/DevopsArtFactory/escli/internal/schema/index"
	"github.com/DevopsArtFactory/escli/internal/util"
)

func GetOSClientFn(url, httpUsername, httpPassword string) *opensearch.Client {
	cfg := opensearch.Config{
		Addresses: []string{
			url,
		},
		Username: httpUsername,
		Password: httpPassword,
	}
	os, _ := opensearch.NewClient(cfg)

	return os
}

func (c Client) OSCatHealth() ([]catSchema.Health, error) {
	var healthMetadata []catSchema.Health

	resp, err := c.OSClient.Cat.Health(
		c.OSClient.Cat.Health.WithFormat("json"))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, util.ErrorWithStatus(resp.StatusCode)
	}
	util.ConvertJSONtoMetadata(resp.Body, &healthMetadata)

	return healthMetadata, nil
}

func (c Client) OSCatIndices(sortKey string) ([]catSchema.Index, error) {
	var indicesMetadata []catSchema.Index

	if sortKey == constants.EmptyString {
		sortKey = constants.DefaultSortKey
	}

	resp, err := c.OSClient.Cat.Indices(
		c.OSClient.Cat.Indices.WithFormat("json"),
		c.OSClient.Cat.Indices.WithS(sortKey))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, util.ErrorWithStatus(resp.StatusCode)
	}
	util.ConvertJSONtoMetadata(resp.Body, &indicesMetadata)

	return indicesMetadata, nil
}

func (c Client) OSCatNodes(sortKey string) ([]catSchema.Node, error) {
	var nodesMetadata []catSchema.Node

	if sortKey == constants.EmptyString {
		sortKey = "id"
	}

	resp, err := c.OSClient.Cat.Nodes(
		c.OSClient.Cat.Nodes.WithFormat("json"),
		c.OSClient.Cat.Nodes.WithH("id,node.role,ip,name,disk.used_percent,load_1m,load_5m,load_15m,uptime,master,disk.total,disk.used,disk.avail"),
		c.OSClient.Cat.Nodes.WithS(sortKey))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, util.ErrorWithStatus(resp.StatusCode)
	}
	util.ConvertJSONtoMetadata(resp.Body, &nodesMetadata)

	return nodesMetadata, nil
}

func (c Client) OSCatShards(sortKey string) ([]catSchema.Shard, error) {
	var shardsMetadata []catSchema.Shard

	if sortKey == constants.EmptyString {
		sortKey = constants.DefaultSortKey
	}

	resp, err := c.OSClient.Cat.Shards(
		c.OSClient.Cat.Shards.WithFormat("json"),
		c.OSClient.Cat.Shards.WithH("index,shard,prirep,state,docs,store,ip,node,unassigned.reason"),
		c.OSClient.Cat.Shards.WithS(sortKey))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, util.ErrorWithStatus(resp.StatusCode)
	}
	util.ConvertJSONtoMetadata(resp.Body, &shardsMetadata)

	return shardsMetadata, nil
}

func (c Client) OSGetIndexSetting(indexName, settingName string) (string, error) {
	var resp *opensearchapi.Response
	var err error

	if settingName == "" {
		resp, err = c.OSClient.Indices.GetSettings(
			c.OSClient.Indices.GetSettings.WithIndex(indexName),
			c.OSClient.Indices.GetSettings.WithPretty(),
		)
	} else {
		resp, err = c.OSClient.Indices.GetSettings(
			c.OSClient.Indices.GetSettings.WithIndex(indexName),
			c.OSClient.Indices.GetSettings.WithName("index."+settingName),
			c.OSClient.Indices.GetSettings.WithPretty(),
		)
	}

	if err != nil {
		return constants.EmptyString, err
	}
	if resp.StatusCode != http.StatusOK {
		return constants.EmptyString, util.ErrorWithStatus(resp.StatusCode)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}

func (c Client) OSPutIndexSetting(indexName, requestBody string) (string, error) {
	resp, err := c.OSClient.Indices.PutSettings(
		strings.NewReader(requestBody),
		c.OSClient.Indices.PutSettings.WithIndex(indexName),
		c.OSClient.Indices.PutSettings.WithPretty(),
	)

	if err != nil {
		return constants.EmptyString, err
	}
	if resp.StatusCode != http.StatusOK {
		return constants.EmptyString, util.ErrorWithStatus(resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}

func (c Client) OSGetClusterSetting() (string, error) {
	var resp *opensearchapi.Response
	var err error

	resp, err = c.OSClient.Cluster.GetSettings(
		c.OSClient.Cluster.GetSettings.WithPretty(),
		c.OSClient.Cluster.GetSettings.WithIncludeDefaults(true))

	if err != nil {
		return constants.EmptyString, err
	}
	if resp.StatusCode != http.StatusOK {
		return constants.EmptyString, util.ErrorWithStatus(resp.StatusCode)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}

func (c Client) OSPutClusterSetting(requestBody string) (string, error) {
	resp, err := c.OSClient.Cluster.PutSettings(
		strings.NewReader(requestBody),
		c.OSClient.Cluster.PutSettings.WithPretty(),
	)

	if err != nil {
		return constants.EmptyString, err
	}
	if resp.StatusCode != http.StatusOK {
		return constants.EmptyString, util.ErrorWithStatus(resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}

func (c Client) OSClusterReroute() (string, error) {
	resp, err := c.OSClient.Cluster.Reroute(
		c.OSClient.Cluster.Reroute.WithRetryFailed(true))

	if err != nil {
		return constants.EmptyString, err
	}
	if resp.StatusCode != http.StatusOK {
		return constants.EmptyString, util.ErrorWithStatus(resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}

func (c Client) OSCreateIndex(index string) (string, error) {
	resp, err := c.OSClient.Indices.Create(
		index,
	)

	if err != nil {
		return constants.EmptyString, err
	}
	if resp.StatusCode != http.StatusOK {
		return constants.EmptyString, util.ErrorWithStatus(resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}

func (c Client) OSDeleteIndex(indices []string) (string, error) {
	resp, err := c.OSClient.Indices.Delete(
		indices,
	)

	if err != nil {
		return constants.EmptyString, err
	}
	if resp.StatusCode != http.StatusOK {
		return constants.EmptyString, util.ErrorWithStatus(resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}

func (c Client) OSIndexStats(index string) (*indexSchema.Stats, error) {
	var indexStatsMetadata indexSchema.Stats

	resp, err := c.OSClient.Indices.Stats(
		c.OSClient.Indices.Stats.WithIndex(index),
	)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, util.ErrorWithStatus(resp.StatusCode)
	}
	util.ConvertJSONtoMetadata(resp.Body, &indexStatsMetadata)

	return &indexStatsMetadata, nil
}

func (c Client) OSStats() (*indexSchema.Stats, error) {
	var indexStatsMetadata indexSchema.Stats

	resp, err := c.OSClient.Indices.Stats()

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, util.ErrorWithStatus(resp.StatusCode)
	}
	util.ConvertJSONtoMetadata(resp.Body, &indexStatsMetadata)

	return &indexStatsMetadata, nil
}
