package gate

import (
	"github.com/lwcbest/gogame/gameserver/game"
	"github.com/lwcbest/gogame/gameserver/msg"
)

//register router
func init() {
	msg.Processor.SetRouter(&msg.Hello{}, game.ChanRPC)
}
