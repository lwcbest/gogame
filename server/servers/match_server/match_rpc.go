package match_server

import (
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/base/utils"
	"gameserver-997/server/domain/entity"
	"gameserver-997/server/service"
)

type MatchRpc struct{}

func (this *MatchRpc) JoinMatchQueue(request *iface.CommonRequest) map[string]interface{} {
	resp := make(map[string]interface{})
	//queueId := request.RpcData[0].(int32)

	player := request.BSession.Setting["player"].(*entity.Player)
	matchService := utils.GlobalObject.CustomData["MatchService"].(*service.MatchService)
	curLen := matchService.MatchPool.JoinPlayer(*player)

	logger.Info("[JoinMatchQueue]Current Length:%v", curLen)
	resp["CurLen"] = curLen
	return resp
}
