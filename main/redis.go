package main

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var pool *redis.Pool

// 创建redis.pool
func initRedis(addr string,idleConn int,maxConn int,idleTimeout time.Duration){
	pool = &redis.Pool{
		MaxIdle: idleConn,
		MaxActive: maxConn,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",addr)
		},
	}
	return
}



