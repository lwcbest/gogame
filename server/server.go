package main

import (
	"encoding/gob"
	xingo "gameserver-997/server/base"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/base/utils"
	"gameserver-997/server/dao/redis"
	"gameserver-997/server/domain/entity"
	"gameserver-997/server/filter"
	"gameserver-997/server/servers/dc_server"
	"gameserver-997/server/servers/game_server"
	"gameserver-997/server/servers/match_server"
	"gameserver-997/server/servers/net_server"
	"gameserver-997/server/service"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"fmt"
)

func main() {
	args := os.Args // exename servername
	dir, err := filepath.Abs(filepath.Dir("."))
	fmt.Println("check file path", dir, filepath.Dir("."))
	if err != nil {
		panic("resolve filepath failed")
	}
	logger.Info("args: ", args, args[1] == "master", os.Getenv("name"))
	name := os.Getenv("name")
	if name == "" {
		name = args[1]
	}
	if name == "master" {
		s := xingo.NewXingoMaster("clusterconf.json")
		s.StartMaster()
	} else {
		s := xingo.NewXingoCluterServer(name, "clusterconf.json")

		//不同服务器加载不同模块(api,http,rpc)
		if strings.HasPrefix(name, "net") {
			s.AddModule("net", &net_server.NetHan{}, &net_server.NetHttp{}, &net_server.NetRpc{})
		}
		if strings.HasPrefix(name, "game") {
			s.AddModule("game", nil, nil, &game_server.GameHan{})
			s.AddModule("game", nil, nil, &game_server.GameRpc{})
			service.StartGameService(name)
		}
		if strings.HasPrefix(name, "dc") {
			s.AddModule("dc", nil, nil, &dc_server.DcHan{})
			s.AddModule("dc", nil, nil, &dc_server.DcRpc{})
			service.StartDataService()
		}

		if strings.HasPrefix(name, "match") {
			s.AddModule("match", nil, nil, &match_server.MatchHan{})
			s.AddModule("match", nil, nil, &match_server.MatchRpc{})
			service.StartMatchService()
		}

		s.StartClusterServer()
	}
}

func init() {
	redis.InitPool()
	redis.TestRedis()
	utils.GlobalObject.OnClosed = func(conn iface.Iconnection) {
		sessionId := fmt.Sprintf("%s:%d", utils.GlobalObject.Name, conn.GetSessionId())
		SessionService := utils.GlobalObject.TcpServer.GetSessionService()

		if Session := SessionService.Get(sessionId); Session != nil {
			logger.Info("GlobalObject.OnClosed %v sid:%s addr: %v", sessionId, Session.GetId(), conn.RemoteAddr())
			SessionService.GetRootCluster().RpcPushServerType(Session, "game", "PlayerOffline")
			// rpc.RpcCallServerType(Session, "game", "PlayerOffline")
			Session.Closed("lost connection")
		} else {
			logger.Info("GlobalObject.OnClosed %v without session addr: %v", sessionId, conn.RemoteAddr())
		}
	}

	//[Filter]register filter!!!!
	utils.GlobalObject.MsgFiltersObj.AddFilter(&filter.LogFilter{})
	utils.GlobalObject.MsgFiltersObj.AddFilter(&filter.AuthFilter{})
	utils.GlobalObject.MsgFiltersObj.AddFilter(&filter.ProtoFilter{})

	//init gob for rpc
	gob.Register(&entity.Player{})
	gob.Register(&sync.Map{})
}
