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
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"

	"github.com/DevopsArtFactory/escli/internal/schema"
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

func (c Client) ListSnapshot() error {
	response, err := c.ESClient.Snapshot.GetRepository()
	if err != nil {
		return err
	}

	var repository map[string]interface{}
	util.ConvertToMap(response.Body, &repository)

	var snapshots []schema.RepositorySnapshot

	for k := range repository {
		response, _ = c.ESClient.Cat.Snapshots(
			c.ESClient.Cat.Snapshots.WithRepository(k),
			c.ESClient.Cat.Snapshots.WithFormat("json"))
		util.ConvertToRepositorySnapshot(response.Body, &snapshots)
		for i := range snapshots {
			fmt.Println(snapshots[i])
		}
	}

	return nil
}

func (c Client) GetRepositoryMetadata(repositoryName string) (string, string) {
	var repository map[string]interface{}

	resp, _ := c.ESClient.Snapshot.GetRepository(
		c.ESClient.Snapshot.GetRepository.WithRepository(repositoryName))
	util.ConvertToMap(resp.Body, &repository)

	repositoryMetadata := repository[repositoryName].(map[string]interface{})
	settings := repositoryMetadata["settings"].(map[string]interface{})

	return settings["bucket"].(string), settings["base_path"].(string)
}

func (c Client) GetSnapshotMetadata(repositoryName string, snapshotName string) io.ReadCloser {
	req := esapi.SnapshotGetRequest{
		Repository: repositoryName,
		Snapshot:   []string{snapshotName},
	}
	resp, _ := req.Do(context.Background(), c.ESClient)

	return resp.Body
}

func (c Client) RestoreSnapshot(requestBody string, repositoryName string, snapshotName string) (*esapi.Response, error) {
	resp, err := c.ESClient.Snapshot.Restore(repositoryName, snapshotName,
		c.ESClient.Snapshot.Restore.WithBody(strings.NewReader(requestBody)))

	return resp, err
}

func (c Client) CatHealth() ([]schema.HealthMetadata, error) {
	var healthMetadata []schema.HealthMetadata

	resp, err := c.ESClient.Cat.Health(
		c.ESClient.Cat.Health.WithFormat("json"))
	if err != nil {
		return nil, err
	}
	util.ConvertToHealthMetadata(resp.Body, &healthMetadata)

	return healthMetadata, nil
}
