package room

import (
	"errors"
	"fmt"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/constants"
	"sync"
	"time"
)

type RoomPool struct {
	PoolId    string
	rooms     sync.Map
	uid2room  sync.Map
	uid2timer sync.Map //掉线的定时器
}

//创建房间
func (this *RoomPool) CreateRoom(roomType int, roomMode string) (IRoom, error) {
	logger.Info("[CreateRoom]createRoom %v %v", roomType, roomMode)
	var room IRoom
	var err error
	switch roomType {
	case constants.ROOM_TYPE_STATE:
		logger.Error("[CreateRoom]no state room!!")
	case constants.ROOM_TYPE_FRAME:
		room, err = createFrameRoomV(roomType, roomMode, this)
	default:
		return nil, errors.New("invalid roomType")
	}
	if err != nil {
		return nil, err
	}

	this.rooms.Store(room.GetRoomId(), room)
	return room, nil
}

func (this *RoomPool) getRoomCount() int {
	count := 0
	this.rooms.Range(func(key interface{}, val interface{}) bool {
		count++
		return true
	})
	return count
}

//销毁房间
func (this *RoomPool) RemoveRoom(id string, feServerId string) {
	this.rooms.Delete(id)
	// go this.checkKickOffOwner(id, feServerId)
}

//每次销毁房间30s内不创建房间就让这个大屏触发一次重连
// func (this *RoomPool) checkKickOffOwner(mgUid string, feServerId string) {
// 	time.Sleep(time.Second * 30)
// 	gaming, err := redis.GetPlayerGaming(mgUid)
// 	if err != nil {
// 		return
// 	}
// 	if gaming.ServerId == "" {
// 		rpc.RpcPushServerName(feServerId, "KickUser", mgUid, "AfterRoomDestroy")
// 	}
// }

func (this *RoomPool) RemovePlayer(uid string) {
	logger.Info("delete uid2room------------------", uid)
	this.uid2room.Delete(uid)
}

func (this *RoomPool) PlayerOkToStart(uid string) error {
	room := this.getPlayerRoom(uid)
	if room == nil {
		return errors.New("no room in")
	}
	return room.PlayerOkToStart(uid)
}

func (this *RoomPool) PlayerQuit(uid string) error {
	room := this.getPlayerRoom(uid)
	if room == nil {
		return errors.New("no room in")
	}
	//room.playerQuit(uid)
	return nil
}

//玩家离线 可能重连
func (this *RoomPool) PlayerOffline(uid string, feServerId string) {
	now := time.Now()
	room := this.getPlayerRoom(uid)
	if room != nil {
		room.PlayerOffline(uid, feServerId)
		logger.Info("playerOffline cost time %s %d", uid, time.Now().Sub(now)/1e6)
	} else {
		logger.Info("playerOffline no room uid: %s", uid, uid)
		fmt.Println(" uid: ", uid)
	}

	// if room != nil && room.getMgUid() == uid && room.isMonitor() {
	// 	if _, ok := this.uid2timer.Load(uid); ok { //只有没有timer的情况下才能
	// 		return
	// 	}

	// 	timer := &clock.SimonTimer{}
	// 	timer.Init(time.Minute*3, func() {

	// 	})
	// 	this.uid2timer.Store(uid, timer)
	// }
}

func (this *RoomPool) SyncFCommand(uid string, epFrame, exFrame int32, tarUid string, cType int32, paramList []int32) error {
	room := this.getPlayerRoom(uid)
	if room == nil {
		return errors.New("no room in" + uid)
	}
	room.SyncFCommand(epFrame, exFrame, tarUid, cType, paramList)
	return nil
}

func (this *RoomPool) SyncResult(uid string, data []byte) error {
	room := this.getPlayerRoom(uid)
	if room == nil {
		return errors.New("no room in")
	}
	return room.SyncResult(uid, data)
}

func (this *RoomPool) getPlayerRoomId(uid string) (string, bool) {
	roomId, ok := this.uid2room.Load(uid)
	if ok == false {
		return "", false
	}
	return roomId.(string), ok
}

func (this *RoomPool) setPlayerRoomId(uid, roomId string) {
	logger.Info("setPlayerRoomId----------------", uid, roomId)
	this.uid2room.Store(uid, roomId)
}

func (this *RoomPool) getPlayerRoom(uid string) IRoom {
	roomId, ok := this.uid2room.Load(uid)
	if !ok {
		return nil
	}
	room, ok := this.rooms.Load(roomId)
	if ok == false {
		return nil
	}
	return room.(IRoom)
}

func (this *RoomPool) GetPlayerRoom(uid string) IRoom {
	return this.getPlayerRoom(uid)
}

func (this *RoomPool) GetRoom(roomId string) IRoom {
	logger.Info("[GetRoom]this.rooms:%v", this.rooms)
	room, ok := this.rooms.Load(roomId)
	if ok == false {
		return nil
	}
	return room.(IRoom)
}
