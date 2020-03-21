package msg

import (
	SimonProto "github.com/lwcbest/gogame/gameserver/goproto"
	"github.com/lwcbest/gogame/gameserver/leaf/network"
)

var Instance = network.NewHybridProcessor()

func init() {
	Instance.Register(&SimonProto.ReqGateGetConnector{})
}
