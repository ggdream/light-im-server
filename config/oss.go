package config

type OSSConfig struct {
	Type  string       `yaml:"type"`
	MinIO *MinIOConfig `yaml:"minio"`
}

type MinIOConfig struct {
	Endpoint        string `yaml:"endpoint"`
	Bucket          string `yaml:"bucket"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	Token           string `yaml:"token"`
	TLS             bool   `yaml:"tls"`
}
