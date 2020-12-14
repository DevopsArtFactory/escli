package config

import (
	"github.com/spf13/viper"

	"github.com/DevopsArtFactory/escli/internal/schema"
)

func GetConfig() (*schema.Config, error) {
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	conf := &schema.Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func GetDefaultConfig() (*schema.Config, error) {
	conf := &schema.Config{}

	return conf, nil
}

func SetInitConfig(elasticsearchURL string, awsRegion string) schema.Config {
	config := schema.Config{
		ElasticSearchURL: elasticsearchURL,
		AWSRegion:        awsRegion,
	}

	return config
}
