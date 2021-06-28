package redis

import "fmt"

//玩家游戏状态保存
func buildPlayerGamingKey(uid string) string {
	str := fmt.Sprintf("gaming:%s", uid)
	return str
}

//创建房间中状态
func buildRoomCreatingKey(deviceId string) string {
	str := fmt.Sprintf("rooming:%s", deviceId)
	return str
}

// 游戏房间和服务器的队列
func buildRoomKey(serverid string, roomId string) string {
	return fmt.Sprintf("%s&%s&queue", roomId, serverid)
}

// 游戏房间和服务器锁
func buildRoomLockKey(serverid string, roomId string) string {
	return fmt.Sprintf("%s&%s&lock", roomId, serverid)
}

//玩家是否在匹配中
func buildPlayerMatchingKey(uid string) string {
	return fmt.Sprintf("matching:%s", uid)
}

//匹配队列对应的Key
func buildMatchQueueKey(deviceId string) string {
	str := fmt.Sprintf("MatchQueue:%s", deviceId)
	return str
}

//匹配队列对应的Key
func buildTimesQueueKey(deviceId string) string {
	str := fmt.Sprintf("MatchQueue:Times:%s", deviceId)
	return str
}
