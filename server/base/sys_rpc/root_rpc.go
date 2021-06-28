package sys_rpc

import (
	// "gameserver-997/server/base/cluster"
	"gameserver-997/server/base/clusterserver"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/base/utils"
)

type RootRpc struct {
}

/*
子节点连上来的通知
*/
// func (this *RootRpc) TakeProxy(request *cluster.RpcRequest) {
// 	name := request.Rpcdata.Args[0].(string)
func (this *RootRpc) TakeProxy(request *iface.CommonRequest) {
	name := request.RpcData[0].(string)
	logger.Info("child node " + name + " connected to " + utils.GlobalObject.Name)
	//加到childs并且绑定链接connetion对象
	clusterserver.GlobalClusterServer.AddChild(name, request.Fconn)
}
