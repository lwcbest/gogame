package main

import (
	"github.com/lwcbest/gogame/gameserver/conf"
	"github.com/lwcbest/gogame/gameserver/game"
	"github.com/lwcbest/gogame/gameserver/gate"
	"github.com/lwcbest/gogame/gameserver/leaf"
	lconf "github.com/lwcbest/gogame/gameserver/leaf/conf"
	"github.com/lwcbest/gogame/gameserver/login"
)

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
