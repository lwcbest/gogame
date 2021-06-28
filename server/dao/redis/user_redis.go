package redis

import (
	"errors"
	"gameserver-997/server/base/logger"
	"time"
)

type PlayerGaming struct {
	ServerId, RoundId, GameType, DeviceId string
	Time                                  int64
	Ip                                    string
}

func SetPlayerGaming(userId string, serverId, deviceId, gameType string, roundId int64, ip string, expire time.Duration) error {
	key := buildPlayerGamingKey(userId)
	pool := getPool()
	m := make(map[string]interface{})
	m["serverId"] = serverId
	m["deviceId"] = deviceId
	m["gameType"] = gameType
	m["roundId"] = roundId
	m["ip"] = ip

	_, err := pool.HMSet(key, m).Result()
	if err != nil {
		logger.Info("SetPlayerGaming failed %s %v err: %+v", key, m, err)
		return err
	}
	if expire > 0 {
		err = pool.Expire(key, expire).Err()
	}
	logger.Info("SetPlayerGaming---------------%s %s %s %+v", userId, serverId, deviceId, err)
	return err
}

func GetPlayerGaming(uid string) (*PlayerGaming, error) {
	key := buildPlayerGamingKey(uid)
	pool = getPool()
	m, err := pool.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	gaming := &PlayerGaming{
		ServerId: m["serverId"],
		RoundId:  m["roundId"],
		GameType: m["gameType"],
		DeviceId: m["deviceId"],
		Ip:       m["ip"],
	}
	return gaming, err
}

func DelPlayerGaming(userId string) error {
	key := buildPlayerGamingKey(userId)
	pool := getPool()
	_, err := pool.Del(key).Result()
	return err
}

func DelDeviceGaming(deviceId string) {
	key := buildPlayerGamingKey(deviceId)
	pool := getPool()
	m := make(map[string]interface{})
	m["serverId"] = ""
	m["deviceId"] = ""
	m["gameType"] = ""
	m["ip"] = ""
	m["roundId"] = time.Now().Unix()

	err := pool.HMSet(key, m).Err()
	if err != nil {
		logger.Error("DelDeviceGaming failed! mac %s err: %v", deviceId, err)
	}
}

func DelUserMatching(uid string, deviceId string) error {
	pool = getPool()
	key := buildPlayerMatchingKey(uid)
	// val, err := pool.Get(key).Result()
	// if val == "" || err != nil {
	// 	logger.Info("DelUserMatching get failed uid:%s err: %+v", uid, err)
	// return errors.New("DelUserMatching failed for no matchKey")
	// val = deviceId
	// }
	queueKey := buildMatchQueueKey(deviceId)
	n, err := pool.LRem(queueKey, 0, uid).Result()
	if err != nil {
		logger.Info("DelUserMatching get failed uid:%s err: %+v", uid, err)
		return err
	}
	if n == 0 {
		return errors.New("remove uid from queue return 0")
	}
	//如果删不掉说明已经退出了这个key也应该删除 且应该进游戏
	pool.Del(key).Err()
	return nil
}

func DelUserTimesMatching(uid string, deviceId string) error {
	pool = getPool()
	key := buildPlayerMatchingKey(uid)
	// val, err := pool.Get(key).Result()
	// if val == "" || err != nil {
	// 	logger.Info("DelUserMatching get failed uid:%s err: %+v", uid, err)
	// 	val = deviceId
	// }
	queueKey := buildTimesQueueKey(deviceId)
	n, err := pool.LRem(queueKey, 0, uid).Result()
	if err != nil {
		logger.Info("DelUserMatching get failed uid:%s err: %+v", uid, err)
		return err
	}

	if n == 0 {
		return errors.New("remove uid from queue return 0")
	}
	//如果删不掉说明已经退出了这个key也应该删除
	pool.Del(key).Err()
	return nil
}
