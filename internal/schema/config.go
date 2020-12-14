package schema

type Config struct {
	ElasticSearchURL string `yaml:"elasticsearchurl"`
	AWSRegion        string `yaml:"awsregion"`
}
