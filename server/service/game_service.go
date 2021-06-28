package service

import (
	"gameserver-997/server/base/utils"
	"gameserver-997/server/domain/room"
)

func StartGameService(servername string) {
	utils.GlobalObject.CustomData["GameService"] = &GameService{}
	utils.GlobalObject.CustomData["GameService"].(*GameService).Init(servername)
}

type GameService struct {
	RoomPool *room.RoomPool
}

func (this *GameService) Init(servername string) {
	this.RoomPool = &room.RoomPool{PoolId: servername}
}
