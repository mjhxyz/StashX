package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool      *redis.Pool
	redisHost = "192.168.60.100:6379"
	redisPass = "000000"
)

// newRedisPool : 创建redis连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,  // 最大空闲连接数
		MaxActive:   30,  // 最大活跃连接数
		IdleTimeout: 300, // 最大空闲时间
		Dial: func() (redis.Conn, error) {
			// 1. 打开连接
			conn, err := redis.Dial("tcp", redisHost)
			if err != nil {
				return nil, err
			}
			// 2. 访问认证
			if _, err := conn.Do("AUTH", redisPass); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			// 检查连接可用性
			if time.Since(t) < time.Minute {
				return nil
			}
			if _, err := c.Do("PING"); err != nil {
				return err
			}
			return nil
		},
	}
}

// init : 初始化连接池
func init() {
	pool = newRedisPool()
}

// RedisPool : 返回连接池实例
func RedisPool() *redis.Pool {
	return pool
}
