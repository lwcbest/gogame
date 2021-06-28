package game_server

import (
	"fmt"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/base/utils"
	"gameserver-997/server/domain/entity"
	"gameserver-997/server/service"
)

type GameRpc struct{}

func (this *GameRpc) CreateRoom(request *iface.CommonRequest) map[string]interface{} {
	resp := make(map[string]interface{})
	roomType := request.RpcData[0].(int)
	roomMode := request.RpcData[1].(string)

	room, err := getGameService().RoomPool.CreateRoom(roomType, roomMode)
	if err != nil {
		logger.Error("[RPC] create room error~:", err)
		resp["err"] = err.Error()
	} else {
		roomId := room.GetRoomId()
		resp["roomId"] = roomId
	}

	return resp
}

func (this *GameRpc) EnterRoom(request *iface.CommonRequest) map[string]interface{} {
	resp := make(map[string]interface{})
	player := request.RpcData[0].(*entity.Player)
	roomId := request.RpcData[1].(string)
	logger.Info("[EnterRoom]:player,roomId:%v %v", player, roomId)
	room := getGameService().RoomPool.GetRoom(roomId)
	if room == nil {
		resp["err"] = "no target room~"
		resp["ok"] = false
		return resp
	}
	logger.Info("[EnterRoom]:player,roomId:%v %v", player, roomId)
	err := room.PlayerEnter(player)
	if err != nil {
		resp["err"] = err.Error()
		resp["ok"] = false
		return resp
	}

	resp["ok"] = true
	return resp
}

func (this *GameRpc) PlayerOffline(request *iface.CommonRequest) {
	fmt.Printf("gamerpc PlayerOffline uid: %s sId:%s\n", request.BSession.Uid, request.BSession.FrontendId)
	//room.GetRoomPool().PlayerOffline(request.BSession.Uid, request.BSession.FrontendId)
}

func getGameService() *service.GameService {
	return utils.GlobalObject.CustomData["GameService"].(*service.GameService)
}
