package rdclient

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient :
type RedisClient struct {
	Client *redis.Client
	Env    string
	Ctx    context.Context
}

// NewRedisClient :
func NewRedisClient(
	ctx context.Context,
	host string,
	port string,
	password string,
	db int,
	env string,
) *RedisClient {
	return &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
			DB:       db,
		}),
		Env: env,
		Ctx: ctx,
	}
}

// GetUserLastUpdateCacheKey :
func (c *RedisClient) GetUserLastUpdateCacheKey(userID uint) string {
	return c.Env + "-lastupdate-" + fmt.Sprint(userID)
}

// GetUserLastUpdate :
func (c *RedisClient) GetUserLastUpdate(userID uint) (string, error) {
	lastUpdate, err := c.Client.Get(
		c.Ctx,
		c.GetUserLastUpdateCacheKey(userID),
	).Result()
	if err == redis.Nil {
		return c.SetUserLastUpdate(userID)
	}
	if err != nil {
		return c.GetTimeNow(), err
	}
	return lastUpdate, nil
}

// SetUserLastUpdate :
func (c *RedisClient) SetUserLastUpdate(userID uint) (string, error) {
	val := c.GetTimeNow()
	key := c.GetUserLastUpdateCacheKey(userID)
	err := c.Client.Set(c.Ctx, key, val, 0).Err()
	return val, err
}

// GetTimeNow :
func (c *RedisClient) GetTimeNow() string {
	return time.Now().Format("2006-01-02T15:04:05.999999")
}
