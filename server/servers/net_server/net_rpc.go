package net_server

import (
	// "gameserver-997/server/base/logger"

	"gameserver-997/server/base/clusterserver"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/utils"
	"gameserver-997/server/domain/entity"
)

type NetRpc struct {
}

func (this *NetRpc) PushMessage(request *iface.CommonRequest) {
	uids := (request.RpcData[0]).([]string)
	route := (request.RpcData[1]).(string)
	data := (request.RpcData[2]).([]byte)
	msg := utils.BuildMsgFromData(utils.MSG_PUSH, 0, route, data)
	sessionService := utils.GlobalObject.TcpServer.GetSessionService()
	sessionService.PushMsgByUids(uids, msg)
}

func (this *NetRpc) PlayerEnterGame(request *iface.CommonRequest) map[string]interface{} {
	resp := make(map[string]interface{})
	roomId := (request.RpcData[0]).(string)
	serverId := (request.RpcData[1]).(string)
	sessionId := (request.RpcData[2]).(string)
	//change session
	sessionService := utils.GlobalObject.TcpServer.GetSessionService()
	session := sessionService.Get(sessionId)
	player := session.Get("player").(*entity.Player)
	player.BeServerId = serverId
	//enter rpc
	clusterserver.GlobalClusterServer.RpcCallServerId(nil, serverId, "EnterRoom", player, roomId)

	//uid := (request.RpcData[0]).([]string)
	// data := (request.RpcData[2]).([]byte)
	// msg := utils.BuildMsgFromData(utils.MSG_PUSH, 0, route, data)
	// sessionService := utils.GlobalObject.TcpServer.GetSessionService()
	// sessionService.PushMsgByUids(uids, msg)
	return resp
}

// func (this *NetRpc) KickUser(request *iface.CommonRequest) {
// 	fmt.Printf("%+v", request)
// 	uid := (request.RpcData[0]).(string)
// 	reason := (request.RpcData[1]).(string)
// 	sessionService := utils.GlobalObject.TcpServer.GetSessionService()
// 	sessionService.Kick(uid, reason)
// 	logger.Info("net_rpc KickUser: ", uid, reason)
// }
