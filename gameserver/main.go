package main

import (
	"gameserver/conf"

	"gameserver/game"

	"gameserver/gate"

	"gameserver/leaf"

	lconf "gameserver/leaf/conf"

	"gameserver/login"
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
