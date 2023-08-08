package config

type MongoConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}
