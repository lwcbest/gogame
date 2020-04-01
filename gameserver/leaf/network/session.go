package network

import (
	"net"
	"reflect"

	"gameserver/leaf/log"

	"github.com/golang/protobuf/proto"
)

type Session struct {
	sid       int
	conn      Conn
	gate      *Gate
	userData  interface{}
	heartbeat *Heartbeat
}

func (s *Session) Run() {
	for {
		pkg, err := s.conn.ReadPkg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}
		if s.gate.Processor != nil {
			msg := s.gate.Processor.HandlePackage(s, pkg)
			if msg == nil {
				continue
			}
			err = s.gate.Processor.Route(msg, s)
			if err != nil {
				log.Error("route message error: %v", err)
				break
			}
		}
	}
}

func (s *Session) OnClose() {
	if s.gate.SessionChanRPC != nil {
		err := s.gate.SessionChanRPC.Call0("CloseSession", s)
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (s *Session) WriteRes(route string, reqId uint, res proto.Message) {
	log.Release("[Response][%d][%s] %v", reqId, route, res)
	//resMsg := res.(*proto.Message)
	msg := BuildMsg(MSG_REQUEST, reqId, route, res)
	s.WriteMsg(msg)
}

func (s *Session) WriteMsg(msg *Message) {
	pkg := BuildPackage(msg, PKG_DATA)
	s.WritePkg(pkg)
}

func (s *Session) WritePkg(pkg *Package) {
	if s.gate.Processor != nil {
		err := s.conn.WritePkg(pkg)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(pkg), err)
		}
	}
}

func (a *Session) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *Session) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *Session) Close() {
	a.conn.Close()
}

func (a *Session) Destroy() {
	a.conn.Destroy()
}

func (a *Session) UserData() interface{} {
	return a.userData
}

func (a *Session) SetUserData(data interface{}) {
	a.userData = data
}
