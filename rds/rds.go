package rds

import (
	"context"
	"log"

	redis "github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"
)

var sf = &singleflight.Group{}

type Options func(*redis.Options)

func NewRedis(o ...Options) *redis.Client {
	options := &redis.Options{}
	for _, item := range o {
		item(options)
	}
	cluster, err, _ := sf.Do("redis", func() (interface{}, error) {
		client := redis.NewClient(options)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if _, err := client.Ping(ctx).Result(); err != nil {
			return nil, err
		}
		return client, nil
	})
	if err != nil {
		log.Panic(err)
	}

	return cluster.(*redis.Client)
}
