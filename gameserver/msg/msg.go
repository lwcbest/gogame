package msg

import (
	SimonProto "gameserver/goproto"

	"gameserver/leaf/network"
)

var Instance = network.NewHybridProcessor()

func init() {
	Instance.Register(&SimonProto.ReqGateGetConnector{})
	Instance.Register(&SimonProto.ReqUserRegister{})
}
