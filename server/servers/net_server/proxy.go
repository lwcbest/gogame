package net_server

import (
	"fmt"
	"gameserver-997/pb/gopb"
	"gameserver-997/server/base/clusterserver"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/constants"
	"time"
)

//RPC
func (this *NetHan) Proxy(serverType string, target string, session iface.ISession, data []byte) (rdata []byte) {
	msg := &gopb.ResError{}
	now1 := time.Now()
	result, err := clusterserver.GlobalClusterServer.RpcSystemCallServerType(session, serverType, target, data)
	if err == nil {
		rdata = result["data"].([]byte)
	} else {
		msg.Code = constants.MSG_CODE.FAIL
		msg.Msg = fmt.Sprintf("rpc error: %s", err.Error())
		rdata, _ = proto.Marshal(msg)
		logger.Info("proxy %s err: %+v cost: %d", target, err, time.Now().Unix()-now1.Unix())
	}
	return
}