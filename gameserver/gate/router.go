package gate

import (
	"gameserver/game"
	SimonProto "gameserver/goproto"
	"gameserver/msg"
)

//register router
func init() {
	msg.Instance.SetRouter(&SimonProto.ReqGateGetConnector{}, game.ChanRPC)
	msg.Instance.SetRouter(&SimonProto.ReqUserRegister{}, game.ChanRPC)
}
