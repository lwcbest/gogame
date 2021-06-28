package fnet

import (
	"errors"
	"fmt"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"sync"
)

type ConnectionMgr struct {
	connections map[uint32]iface.Iconnection
	conMrgLock  sync.RWMutex
}

func (this *ConnectionMgr) Add(conn iface.Iconnection) {
	this.conMrgLock.Lock()
	defer this.conMrgLock.Unlock()
	this.connections[conn.GetSessionId()] = conn
	logger.Debug(fmt.Sprintf("Total connection: %d", len(this.connections)))
}

func (this *ConnectionMgr) Remove(conn iface.Iconnection) error {
	this.conMrgLock.Lock()
	defer this.conMrgLock.Unlock()
	_, ok := this.connections[conn.GetSessionId()]
	if ok {
		delete(this.connections, conn.GetSessionId())
		logger.Info("len of connections: %d", len(this.connections))
		return nil
	} else {
		return errors.New("not found!!")
	}

}

func (this *ConnectionMgr) Get(sid uint32) (iface.Iconnection, error) {
	this.conMrgLock.Lock()
	defer this.conMrgLock.Unlock()
	v, ok := this.connections[sid]
	if ok {
		delete(this.connections, sid)
		return v, nil
	} else {
		return nil, errors.New("not found!!")
	}
}

func (this *ConnectionMgr) Len() int {
	this.conMrgLock.Lock()
	defer this.conMrgLock.Unlock()
	return len(this.connections)
}

func NewConnectionMgr() *ConnectionMgr {
	return &ConnectionMgr{
		connections: make(map[uint32]iface.Iconnection),
	}
}
