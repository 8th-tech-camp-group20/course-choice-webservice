package database

import (
	"course-choice-webservice/config"
	"github.com/garyburd/redigo/redis"
	"time"
)

var RedisClient *redis.Pool

func InitRedis(rconf *config.RedisConfig) {
	RedisClient = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(rconf.Type, rconf.Address)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}
