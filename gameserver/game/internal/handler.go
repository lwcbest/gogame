package internal

import (
	SimonProto "gameserver/goproto"
	"gameserver/leaf/log"
	"gameserver/leaf/network"

	"github.com/golang/protobuf/proto"
)

func init() {
	handler("ReqGateGetConnector", handlerGateGetConnector)
	handler("ReqUserRegister", handlerUserRegister)
}

func handler(router string, h interface{}) {
	skeleton.RegisterChanRPC(router, h)
}

func handlerGateGetConnector(args []interface{}) {
	req := &SimonProto.ReqGateGetConnector{}
	m := args[0].(*network.Message)
	proto.Unmarshal(m.Data, req)

	//core
	session := args[1].(*network.Session)
	res := &SimonProto.ResGateGetConnector{
		Code: 1,
		Host: "127.0.0.1",
		Port: 3014,
	}
	proto.Marshal(res)

	session.WriteRes(m.Route, m.Id, res)
	log.Debug("msg: %v %v", req, session)
}

func handlerUserRegister(args []interface{}) {
	// req := &SimonProto.ReqGateGetConnector{}
	// m := args[0].(*network.Message)
	// proto.Unmarshal(m.Data, req)

	// log.Debug("hello！！！！！！ %v %v", &args[0], &args[1])
	// m := args[0].(*network.Message)
	// session := args[1].(*network.Session)

	// log.Debug("msg: %v %v", m, session)
}
