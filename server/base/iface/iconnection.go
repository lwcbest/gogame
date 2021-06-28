package iface

import (
	"net"
)

type Iconnection interface {
	Start()
	Stop()
	GetConnection() IConn
	GetSessionId() uint32
	Send([]byte) error
	SendBuff([]byte) error
	RemoteAddr() net.Addr
	RemoteIp() string
	LostConnection()
	GetProperty(string) (interface{}, error)
	SetProperty(string, interface{})
	RemoveProperty(string)
}
