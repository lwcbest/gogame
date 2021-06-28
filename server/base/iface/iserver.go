package iface

import (
	"time"
)

type Iserver interface {
	Start()
	Stop()
	Serve()
	GetConnectionMgr() Iconnectionmgr
	GetSessionService() ISessionService
	GetChannelService() IChannelService
	GetConnectionQueue() chan interface{}
	AddRouter(router interface{})
	CallLater(durations time.Duration, f func(v ...interface{}), args ...interface{})
	CallWhen(ts string, f func(v ...interface{}), args ...interface{})
	CallLoop(durations time.Duration, f func(v ...interface{}), args ...interface{})
	GetName() string
	SetService(s ISessionService, c IChannelService)
}
