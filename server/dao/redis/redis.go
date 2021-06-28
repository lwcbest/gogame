package redis

import (
	"fmt"
	"gameserver-997/server/util"

	"github.com/go-redis/redis"
)

var pool *redis.Client

func InitPool() error {
	pool = redis.NewClient(&redis.Options{
		Addr:     util.RedisConf.Address,
		Password: util.RedisConf.Password,
		DB:       util.RedisConf.Db,
	})
	pong, err := pool.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong, "++++++++++++++++++++++++++++++++++++++++++")
	return nil
}

func getPool() *redis.Client {
	return pool
}

func TestRedis() {
	// key := "strKey"
	// beforeSetStr, _ := pool.Get(key).Result()
	// fmt.Println("before : ", beforeSetStr)
	// _, err := pool.Set(key, "strValue", 100000*time.Millisecond).Result()
	// fmt.Println("set err: ", err)
	// after, v := pool.Get(key).Result()
	// fmt.Println("after : ", after, v)

	// key := buildPlayerGamingKey("user0")
	// pool.Del(key)
	// m, err := pool.Exists(key).Result()
	// t, err := pool.Type(key).Result()
	// fmt.Println(t, err)
	// m, err := pool.HMSet(key, map[string]interface{}{
	// 	"entered":  "0",
	// 	"roomId":   "snake:1",
	// 	"serverId": "game1",
	// }).Result()
	// fmt.Println(m, err)
}
