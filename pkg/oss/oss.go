package oss

import (
	"time"

	"github.com/pkg/errors"

	"lim/config"
)

var (
	client OSS
)

type OSS interface {
	PresignPutURL(string, time.Duration) (string, string, error)
	Close() error
}

func Init() (err error) {
	ossCfg := config.GetOSS()
	switch ossCfg.Type {
	case "minio":
		cfg := ossCfg.MinIO
		client, err = NewMinIO(cfg.Endpoint, cfg.Bucket, cfg.AccessKeyID, cfg.SecretAccessKey, cfg.Token, cfg.TLS)
	default:
		err = errors.New("invalid: oss type")
	}

	return
}

func Client() OSS { return client }
