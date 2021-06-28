package redis

func PushMatchQueue(deviceId, val string) error {
	key := buildMatchQueueKey(deviceId)
	pool := getPool()
	_, err := pool.LPush(key, val).Result()
	return err
}

func PopMatchQueue(deviceId string) (string, error) {
	key := buildMatchQueueKey(deviceId)
	pool := getPool()
	val, err := pool.RPop(key).Result()
	return val, err
}

func GetMatchQueueRange(deviceId string, start, stop int64) ([]string, error) {
	key := buildMatchQueueKey(deviceId)
	pool := getPool()
	val, err := pool.LRange(key, start, stop).Result()
	return val, err
}

func GetMatchQueueIndex(deviceId string, index int64) (string, error) {
	key := buildMatchQueueKey(deviceId)
	pool := getPool()
	val, err := pool.LIndex(key, index).Result()
	return val, err
}

func GetTimesQueueIndex(deviceId string, index int64) (string, error) {
	key := buildTimesQueueKey(deviceId)
	pool := getPool()
	val, err := pool.LIndex(key, index).Result()
	return val, err
}
