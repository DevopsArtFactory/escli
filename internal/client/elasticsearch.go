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
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"

	"github.com/DevopsArtFactory/escli/internal/constants"
	catSchema "github.com/DevopsArtFactory/escli/internal/schema/cat"
	snapshotSchema "github.com/DevopsArtFactory/escli/internal/schema/snapshot"
	"github.com/DevopsArtFactory/escli/internal/util"
)

func GetESClientFn(elasticSearchURL string) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticSearchURL,
		},
	}
	es, _ := elasticsearch.NewClient(cfg)

	return es
}

func (c Client) GetRepositories() ([]catSchema.Repository, error) {
	var repositories []catSchema.Repository

	resp, err := c.ESClient.Cat.Repositories(
		c.ESClient.Cat.Repositories.WithFormat("json"))
	if err != nil {
		return nil, err
	}

	util.ConvertJSONtoMetadata(resp.Body, &repositories)

	return repositories, nil
}

func (c Client) GetSnapshots(repositoryID string) (snapshotSchema.Snapshots, error) {
	var snapshots snapshotSchema.Snapshots

	resp, err := c.ESClient.Snapshot.Get(repositoryID, []string{"*"})
	if err != nil {
		return snapshots, err
	}

	util.ConvertJSONtoMetadata(resp.Body, &snapshots)

	return snapshots, nil
}

func (c Client) CreateSnapshot(repositoryID, snapshotID, requestBody string) (int, error) {
	resp, err := c.ESClient.Snapshot.Create(
		repositoryID,
		snapshotID,
		c.ESClient.Snapshot.Create.WithBody(strings.NewReader(requestBody)))
	if err != nil {
		return resp.StatusCode, err
	}

	return resp.StatusCode, util.ReturnErrorFromResponseBody(resp)
}

func (c Client) DeleteSnapshot(repositoryID, snapshotID string) (int, error) {
	resp, err := c.ESClient.Snapshot.Delete(
		repositoryID,
		snapshotID,
	)
	if err != nil {
		return resp.StatusCode, err
	}

	return resp.StatusCode, util.ReturnErrorFromResponseBody(resp)
}

func (c Client) GetRepository(repositoryID string) snapshotSchema.Repository {
	var repositories map[string]snapshotSchema.Repository

	resp, _ := c.ESClient.Snapshot.GetRepository(c.ESClient.Snapshot.GetRepository.WithRepository(repositoryID))

	util.ConvertJSONtoMetadata(resp.Body, &repositories)

	return repositories[repositoryID]
}

func (c Client) GetSnapshot(repositoryID string, snapshotID string) snapshotSchema.Snapshot {
	var snapshots snapshotSchema.Snapshots

	resp, _ := c.ESClient.Snapshot.Get(repositoryID, []string{snapshotID})

	util.ConvertJSONtoMetadata(resp.Body, &snapshots)

	return snapshots.Snapshots[0]
}

func (c Client) RestoreSnapshot(requestBody string, repositoryName string, snapshotName string) (*esapi.Response, error) {
	resp, err := c.ESClient.Snapshot.Restore(repositoryName, snapshotName,
		c.ESClient.Snapshot.Restore.WithBody(strings.NewReader(requestBody)))

	return resp, err
}

func (c Client) CatHealth() ([]catSchema.Health, error) {
	var healthMetadata []catSchema.Health

	resp, err := c.ESClient.Cat.Health(
		c.ESClient.Cat.Health.WithFormat("json"))
	if err != nil {
		return nil, err
	}
	util.ConvertJSONtoMetadata(resp.Body, &healthMetadata)

	return healthMetadata, nil
}

func (c Client) CatIndices(sortKey string) ([]catSchema.Index, error) {
	var indicesMetadata []catSchema.Index

	if sortKey == constants.EmptyString {
		sortKey = "index"
	}

	resp, err := c.ESClient.Cat.Indices(
		c.ESClient.Cat.Indices.WithFormat("json"),
		c.ESClient.Cat.Indices.WithS(sortKey))
	if err != nil {
		return nil, err
	}
	util.ConvertJSONtoMetadata(resp.Body, &indicesMetadata)

	return indicesMetadata, nil
}

func (c Client) CatNodes(sortKey string) ([]catSchema.Node, error) {
	var nodesMetadata []catSchema.Node

	if sortKey == constants.EmptyString {
		sortKey = "id"
	}

	resp, err := c.ESClient.Cat.Nodes(
		c.ESClient.Cat.Nodes.WithFormat("json"),
		c.ESClient.Cat.Nodes.WithH("id,node.role,ip,name,disk.used_percent,load_1m,load_5m,load_15m,uptime"),
		c.ESClient.Cat.Nodes.WithS(sortKey))
	if err != nil {
		return nil, err
	}
	util.ConvertJSONtoMetadata(resp.Body, &nodesMetadata)

	return nodesMetadata, nil
}

func (c Client) CatShards(sortKey string) ([]catSchema.Shard, error) {
	var shardsMetadata []catSchema.Shard

	if sortKey == constants.EmptyString {
		sortKey = "index"
	}

	resp, err := c.ESClient.Cat.Shards(
		c.ESClient.Cat.Shards.WithFormat("json"),
		c.ESClient.Cat.Shards.WithH("index,shard,prirep,state,docs,store,ip,node,unassigned.reason"),
		c.ESClient.Cat.Shards.WithS(sortKey))
	if err != nil {
		return nil, err
	}
	util.ConvertJSONtoMetadata(resp.Body, &shardsMetadata)

	return shardsMetadata, nil
}

func (c Client) GetIndexSetting(indexName, settingName string) (string, error) {
	var resp *esapi.Response
	var err error

	if settingName == "" {
		resp, err = c.ESClient.Indices.GetSettings(
			c.ESClient.Indices.GetSettings.WithIndex(indexName),
			c.ESClient.Indices.GetSettings.WithPretty(),
		)
	} else {
		resp, err = c.ESClient.Indices.GetSettings(
			c.ESClient.Indices.GetSettings.WithIndex(indexName),
			c.ESClient.Indices.GetSettings.WithName("index."+settingName),
			c.ESClient.Indices.GetSettings.WithPretty(),
		)
	}

	if err != nil {
		return constants.EmptyString, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}

func (c Client) PutIndexSetting(indexName, requestBody string) (string, error) {
	resp, err := c.ESClient.Indices.PutSettings(
		strings.NewReader(requestBody),
		c.ESClient.Indices.PutSettings.WithIndex(indexName),
		c.ESClient.Indices.PutSettings.WithPretty(),
	)

	if err != nil {
		return constants.EmptyString, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}

func (c Client) GetClusterSetting() (string, error) {
	var resp *esapi.Response
	var err error

	resp, err = c.ESClient.Cluster.GetSettings(
		c.ESClient.Cluster.GetSettings.WithPretty(),
		c.ESClient.Cluster.GetSettings.WithIncludeDefaults(true))

	if err != nil {
		return constants.EmptyString, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}

func (c Client) PutClusterSetting(requestBody string) (string, error) {
	resp, err := c.ESClient.Cluster.PutSettings(
		strings.NewReader(requestBody),
		c.ESClient.Cluster.PutSettings.WithPretty(),
	)

	if err != nil {
		return constants.EmptyString, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}

func (c Client) ClusterReroute() (string, error) {
	resp, err := c.ESClient.Cluster.Reroute(
		c.ESClient.Cluster.Reroute.WithRetryFailed(true))

	if err != nil {
		return constants.EmptyString, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.String(), nil
}
