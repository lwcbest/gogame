package network

import (
	"time"

	"gameserver/leaf/chanrpc"
)

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       *Processor
	SessionChanRPC  *chanrpc.Server

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool
	CurId        int
}

func (gate *Gate) Run(closeSig chan bool) {
	gate.CurId = 1
	//var wsServer *WSServer
	//if gate.WSAddr != "" {
	//	wsServer = new(WSServer)
	//	wsServer.Addr = gate.WSAddr
	//	wsServer.MaxConnNum = gate.MaxConnNum
	//	wsServer.PendingWriteNum = gate.PendingWriteNum
	//	wsServer.MaxMsgLen = gate.MaxMsgLen
	//	wsServer.HTTPTimeout = gate.HTTPTimeout
	//	wsServer.CertFile = gate.CertFile
	//	wsServer.KeyFile = gate.KeyFile
	//	wsServer.NewSession = func(conn *WSConn) *Session {
	//		session := &Session{conn: conn, gate: gate}
	//		if gate.SessionChanRPC != nil {
	//			gate.SessionChanRPC.Go("NewSession", session)
	//		}
	//		return session
	//	}
	//}

	var tcpServer *TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.NewSession = func(conn *TCPConn) *Session {
			gate.CurId++
			session := &Session{
				sid:  gate.CurId,
				conn: conn,
				gate: gate,
			}
			session.heartbeat = &Heartbeat{
				timeout:   HeartbeatDuration * 2,
				heartbeat: HeartbeatDuration,
				session:   session,
				ch:        make(chan bool),
				closed:    false,
			}
			if gate.SessionChanRPC != nil {
				gate.SessionChanRPC.Go("NewSession", session)
			}
			return session
		}
	}

	//if wsServer != nil {
	//	wsServer.Start()
	//}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	//if wsServer != nil {
	//	wsServer.Close()
	//}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}
