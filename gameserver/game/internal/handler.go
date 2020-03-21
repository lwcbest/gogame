package internal

import (
	"github.com/lwcbest/gogame/gameserver/leaf/log"
	"github.com/lwcbest/gogame/gameserver/leaf/network"
)

func init() {
	handler("ReqGateGetConnector", handleHello)
}

func handler(router string, h interface{}) {
	skeleton.RegisterChanRPC(router, h)
}

func handleHello(args []interface{}) {
	//m := args[0].(*msg.Hello)
	session := args[1].(network.Session)
	log.Debug("hello %v", session)
	//session.WriteMsg(&msg.Hello{
	//	Name:proto.String("client"),
	//})
}
