package gate

import (
	"myGo/gameserver/game"
	"myGo/gameserver/msg"
)

//register router
func init() {
	msg.Processor.SetRouter(&msg.Hello{}, game.ChanRPC)
}
