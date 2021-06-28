package sys_rpc

import (
	"fmt"
	// "gameserver-997/server/base/cluster"
	"gameserver-997/server/base/clusterserver"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/base/utils"
	"os"
	"time"
)

type ChildRpc struct {
}

/*
master 通知父节点上线, 收到通知的子节点需要链接对应父节点
*/
// func (this *ChildRpc) RootTakeProxy(request *cluster.RpcRequest) {
// rname := request.Rpcdata.Args[0].(string)
func (this *ChildRpc) RootTakeProxy(request *iface.CommonRequest) {
	rname := request.RpcData[0].(string)
	logger.Info(fmt.Sprintf("root node %s online. connecting...", rname))
	clusterserver.GlobalClusterServer.ConnectToRemote(rname)
}

/*
关闭节点信号
*/
// func (this *ChildRpc) CloseServer(request *cluster.RpcRequest){
// 	delay := request.Rpcdata.Args[0].(int)
func (this *ChildRpc) CloseServer(request *iface.CommonRequest) {
	delay := request.RpcData[0].(int)
	logger.Error("server close kickdown.", delay, "second...")
	time.Sleep(time.Duration(delay) * time.Second)
	utils.GlobalObject.ProcessSignalChan <- os.Kill
}

/*
重新加载配置文件
*/
// func (this *ChildRpc) ReloadConfig(request *cluster.RpcRequest){
// 	delay := request.Rpcdata.Args[0].(int)
func (this *ChildRpc) ReloadConfig(request *iface.CommonRequest) {
	delay := request.RpcData[0].(int)
	logger.Error("server ReloadConfig kickdown.", delay, "second...")
	time.Sleep(time.Duration(delay) * time.Second)
	clusterserver.GlobalClusterServer.Cconf.Reload()
	utils.GlobalObject.Reload()
	logger.Info("reload config.")
}

/*
检查节点是否下线
*/
// func (this *ChildRpc) CheckAlive(request *cluster.RpcRequest)(response map[string]interface{}){
func (this *ChildRpc) CheckAlive(request *iface.CommonRequest) (response map[string]interface{}) {
	response = make(map[string]interface{})
	conns := utils.GlobalObject.TcpServer.GetConnectionMgr().Len()

	response["ext"] = fmt.Sprintf("%d", conns)
	response["name"] = clusterserver.GlobalClusterServer.Name
	logger.Debug("CheckAlive!%+v", response)
	return
}

/*
通知节点掉线（父节点或子节点）
*/
// func (this *ChildRpc)NodeDownNtf(request *cluster.RpcRequest) {
// 	isChild := request.Rpcdata.Args[0].(bool)
// 	nodeName := request.Rpcdata.Args[1].(string)
func (this *ChildRpc) NodeDownNtf(request *iface.CommonRequest) {
	isChild := request.RpcData[0].(bool)
	nodeName := request.RpcData[1].(string)
	logger.Debug(fmt.Sprintf("node %s down ntf.", nodeName))
	if isChild {
		clusterserver.GlobalClusterServer.RemoveChild(nodeName)
	} else {
		clusterserver.GlobalClusterServer.RemoveRemote(nodeName)
	}
}
