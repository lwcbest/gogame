package room

import "gameserver-997/server/domain/entity"

type IRoom interface {
	GetRoomId() string
	GetPlayerCount() int
	Destroy()

	Start()
	SyncFCommand(epFrame, exFrame int32, uid string, ctype int32, paramList []int32)
	SyncResult(string, []byte) error
	PlayerEnter(player *entity.Player) error
	PlayerReady(uid string) error
	PlayerOkToStart(uid string) error
	PlayerOffline(string, string)
	GetState() int32
}
