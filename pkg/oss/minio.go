package oss

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var _ OSS = (*MinIO)(nil)

type MinIO struct {
	accessKeyId, secretAccessKey, token string
	tls                                 bool
	bucket, region                      string
	rootCtx                             context.Context
	baseTime                            time.Duration
	client                              *minio.Client
}

func NewMinIO(endpoint, bucket, accessKeyId, secretAccessKey, token string, tls bool) (*MinIO, error) {
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, token),
		Secure: tls,
	}
	client, err := minio.New(endpoint, opts)
	if err != nil {
		return nil, err
	}

	c := &MinIO{
		accessKeyId:     accessKeyId,
		secretAccessKey: secretAccessKey,
		bucket:          bucket,
		token:           token,
		tls:             tls,
		client:          client,
		rootCtx:         context.TODO(),
		baseTime:        time.Second * 3,
	}
	err = c.initBucket(bucket, fmt.Sprintf(readonlyBucketPolicy, bucket))
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *MinIO) PresignPutURL(objectName string, expires time.Duration) (string, string, error) {
	ctx, cancel := context.WithTimeout(c.rootCtx, c.baseTime)
	defer cancel()

	url, err := c.client.PresignedPutObject(ctx, c.bucket, objectName, expires)
	if err != nil {
		return "", "", err
	}

	rawUrl := url.String()
	url.RawQuery = ""

	return rawUrl, url.String(), nil
}

func (c *MinIO) Close() error {
	return nil
}

// initBucket 初始化某个桶
func (c *MinIO) initBucket(name, policy string) error {
	ctx, cancel := context.WithTimeout(c.rootCtx, c.baseTime)
	defer cancel()

	isExist, err := c.client.BucketExists(ctx, name)
	if err != nil {
		return err
	}
	if !isExist {
		ctx1, cancel1 := context.WithTimeout(c.rootCtx, c.baseTime)
		defer cancel1()

		opts := minio.MakeBucketOptions{
			Region:        c.region,
			ObjectLocking: false,
		}
		if err := c.client.MakeBucket(ctx1, name, opts); err != nil {
			return err
		}
	}

	ctx2, cancel2 := context.WithTimeout(c.rootCtx, c.baseTime)
	defer cancel2()

	return c.client.SetBucketPolicy(ctx2, name, policy)
}
