package redis

import (
	"gameserver-997/server/base/logger"
	"time"
)

//检查玩家是否在房间列表里
func CheckUserInRoom(uid, serverId, roomId string, playerNum int) (bool, error) {
	return true, nil
	key := buildRoomKey(serverId, roomId)
	pool := getPool()
	uids, err := pool.LRange(key, 0, int64(playerNum)).Result()
	if err != nil {
		return false, nil
	}
	for _, u := range uids {
		if u == uid {
			return true, nil
		}
	}
	return false, nil
}

//检查玩家是否在候补席里面
func CheckUserInWaiting(uid, serverId, roomId string, playerNum int) (bool, error) {
	return true, nil
	key := buildRoomKey(serverId, roomId)
	pool := getPool()
	uids, err := pool.LRange(key, int64(playerNum), int64(playerNum*2)).Result()
	if err != nil {
		return false, nil
	}
	for _, u := range uids {
		if u == uid {
			return true, nil
		}
	}
	return false, nil
}

// 设置房间状态为为创建中
func SetRoomCreating(deviceId string, expire time.Duration) (bool, error) {
	key := buildRoomCreatingKey(deviceId)
	pool := getPool()

	done, err := pool.HSetNX(key, deviceId, deviceId).Result()
	if err != nil {
		logger.Info("SetRoomCreating failed %s %v err: %+v", key, deviceId, err)
		return false, err
	}
	if expire > 0 && done == true {
		err = pool.Expire(key, expire).Err()
		if err != nil {
			// 过期设置失败，要删除缓存
			_, err := pool.Del(key).Result()
			return false, err
		}
	}
	logger.Info("SetRoomCreating---------------%s", deviceId, err)
	return done, err
}

// 解除创建中状态
func DelRoomCreating(deviceId string) error {
	key := buildRoomCreatingKey(deviceId)
	pool := getPool()
	_, err := pool.Del(key).Result()
	return err
}
