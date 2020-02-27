package main

import (
	"myGo/gameserver/leaf"
	lconf "myGo/gameserver/leaf/conf"
	"myGo/gameserver/conf"
	"myGo/gameserver/game"
	"myGo/gameserver/gate"
	"myGo/gameserver/login"
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
