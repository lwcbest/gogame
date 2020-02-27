package internal

import (
	"github.com/lwcbest/gogame/gameserver/leaf/log"
	"github.com/lwcbest/gogame/gameserver/leaf/network"
)

func init() {
	skeleton.RegisterChanRPC("NewSession", rpcNewSession)
	skeleton.RegisterChanRPC("CloseSession", rpcCloseSession)
}

var sessions = make(map[*network.Session]struct{})

func rpcNewSession(args []interface{}) {
	session := args[0].(*network.Session)
	log.Debug("new session:", session)
	sessions[session] = struct{}{}
}

func rpcCloseSession(args []interface{}) {
	session := args[0].(*network.Session)
	log.Debug("del session:", session)
	delete(sessions, session)
}
