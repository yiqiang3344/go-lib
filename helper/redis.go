package helper

import (
	"github.com/garyburd/redigo/redis"
	"github.com/micro/go-micro/v2/config"
	"strconv"
	"sync"
	"time"
)

var poolMap map[string]*redis.Pool
var genPoolOnceMap map[string]*sync.Once

func DefaultRedis() redis.Conn {
	return InitRedisPool("redis").Get()
}

func InitRedisPool(cfg string) *redis.Pool {
	if len(genPoolOnceMap) == 0 {
		genPoolOnceMap = make(map[string]*sync.Once)
	}
	if len(poolMap) == 0 {
		poolMap = make(map[string]*redis.Pool)
	}
	if _, ok := genPoolOnceMap[cfg]; !ok {
		genPoolOnceMap[cfg] = new(sync.Once)
	}
	genPoolOnceMap[cfg].Do(func() {
		//DebugLog("redis pool create:"+cfg, "")
		poolMap[cfg] = &redis.Pool{
			MaxIdle:     10, //空闲数
			IdleTimeout: 300 * time.Second,
			MaxActive:   20, //最大数
			Dial: func() (redis.Conn, error) {
				cfgMap := config.Get(cfg).StringMap(map[string]string{})
				//从mysql查询biz_type配置
				database, _ := strconv.Atoi(cfgMap["database"])
				c, err := redis.Dial(
					"tcp",
					cfgMap["host"]+":"+cfgMap["port"],
					redis.DialDatabase(database),
				)
				if err != nil {
					return nil, err
				}
				//DebugLog("redis connect:"+cfg, "")
				if cfgMap["password"] != "" {
					if _, err := c.Do("AUTH", cfgMap["password"]); err != nil {
						//DebugLog("redis close:"+cfg, "")
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				//DebugLog("redis ping:"+cfg, "")
				_, err := c.Do("PING")
				return err
			},
		}
	})
	return poolMap[cfg]
}
