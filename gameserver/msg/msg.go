package msg

import (
	"myGo/gameserver/leaf/network"
)

var Processor = network.NewHybridProcessor()

func init() {
	Processor.Register(&Hello{})
}
