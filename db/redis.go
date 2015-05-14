package db

import (
    "fmt"
    "time"

    "github.com/nicholaskh/redigo/redis"
)

var (
    redisConn redis.Conn
)

func RedisConn(addr string, connTimeout, readTimeout, writeTimeout time.Duration) redis.Conn {
    var err error
    if redisConn == nil {
        redisConn, err = redis.DialTimeout("tcp", addr, connTimeout,
            readTimeout, writeTimeout)
        if err != nil {
            panic(fmt.Sprintf("Connect to redis error: %s", err.Error()))
        }
    }
    return redisConn
}

