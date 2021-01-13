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

package cats

type Health struct {
	Cluster             string `json:"cluster"`
	Status              string `json:"status"`
	NodeTotal           string `json:"node.total"`
	NodeData            string `json:"node.data"`
	Shards              string `json:"shards"`
	InitShards          string `json:"unassign"`
	ActiveShardsPercent string `json:"active_shards_percent"`
}

type Index struct {
	Health        string `json:"health"`
	Status        string `json:"status"`
	Index         string `json:"Index"`
	UUID          string `json:"uuid"`
	PrimaryShards string `json:"pri"`
	ReplicaShards string `json:"rep"`
	StoreSize     string `json:"store.size"`
}

type Node struct {
	IP              string `json:"ip"`
	NodeRole        string `json:"node.role"`
	Name            string `json:"name"`
	DiskUsedPercent string `json:"disk.used_percent"`
	Load1M          string `json:"load_1m"`
	Load5M          string `json:"load_5m"`
	Load15M         string `json:"load_15m"`
	Uptime          string `json:"uptime"`
}

type Shard struct {
	Index            string `json:"index"`
	Shard            string `json:"shard"`
	PriRep           string `json:"prirep"`
	State            string `json:"state"`
	Docs             string `json:"docs"`
	Store            string `json:"store"`
	IP               string `json:"ip"`
	Node             string `json:"node"`
	UnassignedReason string `json:"unassigned.reason"`
}

type Repository struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
