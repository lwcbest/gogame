package internal

import (
	"gameserver/leaf/network"
)

func init() {
	skeleton.RegisterChanRPC("NewSession", rpcNewSession)
	skeleton.RegisterChanRPC("CloseSession", rpcCloseSession)
}

var sessions = make(map[*network.Session]struct{})

func rpcNewSession(args []interface{}) {
	session := args[0].(*network.Session)
	sessions[session] = struct{}{}
}

func rpcCloseSession(args []interface{}) {
	session := args[0].(*network.Session)
	delete(sessions, session)
}
