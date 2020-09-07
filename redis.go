package main

import (
	"time"
	"log"
	"github.com/garyburd/redigo/redis"
)

var (
	redisPool *redis.Pool
	redisAddr string
)

func initRedis() {
	redisPool = &redis.Pool {
		MaxIdle:      10,
		MaxActive:    1000,
		IdleTimeout:  300 * time.Second,
		Dial:       func() (redis.Conn, error) {
					conn, err := redis.Dial(
						"tcp",
						redisAddr,
						redis.DialConnectTimeout(time.Duration(5000 * time.Millisecond)),
						redis.DialReadTimeout(time.Duration(180000 * time.Millisecond)),
						redis.DialWriteTimeout(time.Duration(3000 * time.Millisecond)),
					)
		
					if err != nil {
						return nil, err
					}
		
					_, err = conn.Do("SELECT", "0")
					if err != nil {
						conn.Close()
						return nil, err
					}
		
					return conn, nil
				},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
						_, err := conn.Do("PING")
						if err != nil {
							log.Printf("无效连接_%s", err.Error())
						}
						return err
					},
		Wait: true,
	}
}

// 执行redis命令
func execRedisCommand(command string, args ...interface{}) (interface{}, error) {
	redis := redisPool.Get()
	defer redis.Close()

	return redis.Do(command, args...)
}