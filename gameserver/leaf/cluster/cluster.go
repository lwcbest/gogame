package cluster

import (
	"math"
	"time"

	"github.com/lwcbest/gogame/gameserver/leaf/conf"
	"github.com/lwcbest/gogame/gameserver/leaf/network"
)

var (
	server  *network.TCPServer
	clients []*network.TCPClient
)

func Init() {
	if conf.ListenAddr != "" {
		server = new(network.TCPServer)
		server.Addr = conf.ListenAddr
		server.MaxConnNum = int(math.MaxInt32)
		server.PendingWriteNum = conf.PendingWriteNum
		server.NewSession = newSession

		server.Start()
	}

	for _, addr := range conf.ConnAddrs {
		client := new(network.TCPClient)
		client.Addr = addr
		client.ConnNum = 1
		client.ConnectInterval = 3 * time.Second
		client.PendingWriteNum = conf.PendingWriteNum
		client.NewSession = newSession

		client.Start()
		clients = append(clients, client)
	}
}

func Destroy() {
	if server != nil {
		server.Close()
	}

	for _, client := range clients {
		client.Close()
	}
}

func newSession(conn *network.TCPConn) *network.Session {
	return &network.Session{}
}
