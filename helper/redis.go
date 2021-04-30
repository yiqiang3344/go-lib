package helper

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"sync"
	"time"
)

var poolMap map[string]*redis.Pool
var genPoolOnceMap map[string]*sync.Once

func DefaultRedis() redis.Conn {
	return InitRedisPool("redis").Get()
}

func GenRedisKey(args ...string) string {
	return GetCfgString("project") + ":" + strings.Join(args, ":")
}

func InitRedisPool(name string) *redis.Pool {
	if len(genPoolOnceMap) == 0 {
		genPoolOnceMap = make(map[string]*sync.Once)
	}
	if len(poolMap) == 0 {
		poolMap = make(map[string]*redis.Pool)
	}
	if _, ok := genPoolOnceMap[name]; !ok {
		genPoolOnceMap[name] = new(sync.Once)
	}
	genPoolOnceMap[name].Do(func() {
		//DebugLog("redis pool create:"+name, "")
		poolMap[name] = &redis.Pool{
			MaxIdle:     10, //空闲数
			IdleTimeout: 300 * time.Second,
			MaxActive:   20, //最大数
			Dial: func() (redis.Conn, error) {
				cfgMap := GetCfgStringMap(name)
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
				//DebugLog("redis connect:"+name, "")
				if cfgMap["password"] != "" {
					if _, err := c.Do("AUTH", cfgMap["password"]); err != nil {
						//DebugLog("redis close:"+name, "")
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				//DebugLog("redis ping:"+name, "")
				_, err := c.Do("PING")
				return err
			},
		}
	})
	return poolMap[name]
}
