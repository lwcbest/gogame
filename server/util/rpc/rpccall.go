package rpc

import (
	"gameserver-997/server/base/clusterserver"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/utils"
)

func RPCRandomCall(session iface.ISession, serverType string, target string, args ...interface{}) (map[string]interface{}, error) {
	return clusterserver.GlobalClusterServer.RpcRandomCallServerType(session, serverType, target, args)
}

func RPCCall(session iface.ISession, serverId string, target string, args ...interface{}) (map[string]interface{}, error) {
	return clusterserver.GlobalClusterServer.RpcCallServerId(session, serverId, target, args)
}

func PushMsgByUids(route string, msg interface{}, sessionArray []iface.ISession) {
	utils.GlobalObject.TcpServer.GetChannelService().PushMsgByUids(route, msg, sessionArray)
}
