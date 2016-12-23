package db

import (
	"ggs/log"
	"ggsserver/conf"

	"github.com/mediocregopher/radix.v2/pool"
)

var Redis *pool.Pool

func init() {
	redis, err := pool.New("tcp", "localhost:"+conf.Server.RedisPort, 5)
	if err != nil {
		log.Fatal("redis pool creation failed: %v", err)
	}

	Redis = redis
}
