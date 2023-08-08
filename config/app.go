package config

type AppConfig struct {
	Addr     string     `yaml:"addr"`
	TokenKey []byte     `yaml:"token_key"`
	TLS      *TLSConfig `yaml:"tls"`
}

type TLSConfig struct {
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}
