package base

import (
	"github.com/lwcbest/gogame/gameserver/conf"
	"github.com/lwcbest/gogame/gameserver/leaf/chanrpc"
	"github.com/lwcbest/gogame/gameserver/leaf/module"
)

func NewSkeleton() *module.Skeleton {
	skeleton := &module.Skeleton{
		GoLen:              conf.GoLen,
		TimerDispatcherLen: conf.TimerDispatcherLen,
		AsynCallLen:        conf.AsynCallLen,
		ChanRPCServer:      chanrpc.NewServer(conf.ChanRPCLen),
	}
	skeleton.Init()
	return skeleton
}
