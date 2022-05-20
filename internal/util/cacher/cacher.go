package cacher

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"go-clean-architecture/internal/util/env"
)

type ICacher interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) *string
	Delete(key string) error
	Ping() error
}

type Cacher struct {
	client      *redis.Client
	environment env.IEnvironment
	timeout     time.Duration
}

// New
// Returns a new Cacher.
func New(environment env.IEnvironment) ICacher {
	return &Cacher{
		environment: environment,
		client: redis.NewClient(&redis.Options{
			Addr:     environment.Get(env.RedisAddress),
			Password: environment.Get(env.RedisPassword),
			DB:       0,
		}),
		timeout: time.Second * 5,
	}
}

func (c *Cacher) Set(key string, value interface{}, duration time.Duration) {
	bytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	err = c.client.Set(ctx, key, bytes, duration).Err()
	if err != nil {
		panic(err)
	}
}

func (c *Cacher) Get(key string) *string {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	stringCmd := c.client.Get(ctx, key)

	if stringCmd.Val() == "" {
		return nil
	}

	value, err := stringCmd.Result()
	if err != nil {
		panic(err)
	}

	return &value
}

func (c *Cacher) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	intCmd := c.client.Del(ctx, key)
	err := intCmd.Err()
	if err != nil {
		return err
	}

	return nil
}
func (c *Cacher) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	statusCmd := c.client.Ping(ctx)
	err := statusCmd.Err()
	if err != nil {
		return err
	}

	return nil
}
