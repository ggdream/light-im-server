package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var (
	c *Config
)

type Config struct {
	App   *AppConfig   `yaml:"app"`
	OSS   *OSSConfig   `yaml:"oss"`
	Mongo *MongoConfig `yaml:"mongo"`
	Redis *RedisConfig `yaml:"redis"`
}

func Init() error {
	file, err := os.Open("./config.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	return yaml.NewDecoder(file).Decode(&c)
}

func SetFromManual(cfg *Config) { c = cfg }

func GetApp() *AppConfig { return c.App }

func GetMongo() *MongoConfig { return c.Mongo }

func GetRedis() *RedisConfig { return c.Redis }

func GetOSS() *OSSConfig { return c.OSS }
