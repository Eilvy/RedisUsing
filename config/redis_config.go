package config

import (
	"github.com/redis/go-redis/v9"
	"time"
)

func RedisConnect() {
	client := redis.NewClient(&redis.Options{
		Addr:        "redis-14520.c299.asia-northeast1-1.gce.cloud.redislabs.com:14520",
		Password:    "rPYdtUeiD5CeJSqcGZdoyHDd6Ou2uApa",
		DB:          0,
		DialTimeout: time.Second * 2,
	}) //配置redis云端数据库
	DB = client //将redis连接状态赋予DB
}
