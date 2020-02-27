package msg

import (
	"github.com/lwcbest/gogame/gameserver/leaf/network"
)

var Processor = network.NewHybridProcessor()

func init() {
	Processor.Register(&Hello{})
}
