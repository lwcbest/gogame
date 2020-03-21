package gate

import (
	"github.com/lwcbest/gogame/gameserver/game"
	SimonProto "github.com/lwcbest/gogame/gameserver/goproto"
	"github.com/lwcbest/gogame/gameserver/msg"
)

//register router
func init() {
	msg.Instance.SetRouter(&SimonProto.ReqGateGetConnector{}, game.ChanRPC)
}
