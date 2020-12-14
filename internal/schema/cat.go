package schema

type HealthMetadata struct {
	Cluster             string `json:"cluster"`
	Status              string `json:"status"`
	NodeTotal           string `json:"node.total"`
	NodeData            string `json:"node.data"`
	Shards              string `json:"shards"`
	InitShards          string `json:"unassign"`
	ActiveShardsPercent string `json:"active_shards_percent"`
}
