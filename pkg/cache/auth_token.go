package cache

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthToken struct{}

func NewAuthToken() *AuthToken { return new(AuthToken) }

func (c *AuthToken) Set(userId string, token string, expireAt int64) error {
	ctx, cancel := withTimeout()
	defer cancel()

	return client.SetArgs(ctx, c.joinName(userId), token, redis.SetArgs{
		ExpireAt: time.UnixMilli(expireAt),
	}).Err()
}

func (c *AuthToken) Get(userId string) (string, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	return client.Get(ctx, c.joinName(userId)).Result()
}

func (c *AuthToken) Verify(userId string, token string) (bool, error) {
	ctx, cancel := withTimeout()
	defer cancel()

	tk, err := client.Get(ctx, c.joinName(userId)).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, err
	}

	return tk == token, nil
}

func (c *AuthToken) Del(userId string) error {
	ctx, cancel := withTimeout()
	defer cancel()

	return client.Unlink(ctx, c.joinName(userId)).Err()
}

func (c *AuthToken) joinName(userId string) string {
	return fmt.Sprintf("lim:auth:token:%s", userId)
}
