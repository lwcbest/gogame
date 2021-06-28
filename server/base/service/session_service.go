package service

import (
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"sync"
	"time"
)

//SessionService session func...
type SessionService struct {
	RootCluster iface.Icluster
	RootServer  iface.Iserver
	Sessions    sync.Map
	UidMap      sync.Map
	uid2sids    sync.Map
}

func (this *SessionService) GetRootServer() iface.Iserver {
	return this.RootServer
}

func (this *SessionService) GetRootCluster() iface.Icluster {
	return this.RootCluster
}

func (this *SessionService) Create(sid string, frontendId string, socket iface.Iconnection) iface.ISession {
	session := &Session{
		Id:             sid,
		FrontendId:     frontendId,
		socket:         socket,
		sessionService: this,
		closeChan:      make(chan int),
	}

	this.Sessions.Store(sid, session)
	return session
}

func (this *SessionService) getSession(sid string) iface.ISession {
	session, ok := this.Sessions.Load(sid)
	if !ok {
		return nil
	}
	return session.(iface.ISession)
}

func (this *SessionService) getSidByUid(uid string) []string {
	// sid, ok := this.UidMap.Load(uid)
	// if !ok {
	// 	return nil
	// }
	sids, ok := this.uid2sids.Load(uid)
	if !ok {
		return nil
	}
	return sids.([]string)
}

func (this *SessionService) Get(sid string) iface.ISession {
	return this.getSession(sid)
}

func (this *SessionService) Bind(sid string, uid string) {
	session := this.getSession(sid)
	if session == nil {
		return
	}
	// this.UidMap.Store(uid, sid)
	sids := this.getSidByUid(uid)
	this.uid2sids.Store(uid, append(sids, sid))
}

func (this *SessionService) UnBind(sid string, uid string) {
	session := this.getSession(sid)
	if session == nil {
		return
	}
	if session.GetUid() != uid {
		return
	}

	session.UnBind(uid)
	sids := this.getSidByUid(uid)
	for i, v := range sids {
		if v == sid {
			this.uid2sids.Store(uid, append(sids[0:i], sids[i+1:]...))
		}
	}
}

func (this *SessionService) Remove(sid string) {
	session := this.getSession(sid)
	if session == nil {
		return
	}
	if session.GetUid() != "" {
		this.UnBind(sid, session.GetUid())
	}
	this.Sessions.Delete(sid)
}

func (this *SessionService) ImportData(sid string, key string, value interface{}) {
	session := this.getSession(sid)
	if session == nil {
		return
	}

	session.Set(key, value)
}

func (this *SessionService) KickBySId(sid string, reason string) {
	session := this.getSession(sid)
	if session == nil {
		return
	}

	// if session.GetUid() != "" {
	// 	this.UnBind(sid, session.GetUid())
	// }

	session.Kick(reason)
	this.Remove(sid)
}

func (this *SessionService) SendData(sid string, data []byte) {
	session := this.getSession(sid)
	if session == nil {
		return
	}
	session.Send(data)
}

func (this *SessionService) SendDataByUid(uid string, data []byte) {
	// sid := this.getSidByUid(uid)
	// if sid == "" {
	// 	logger.Info("uid no sid", uid)
	// 	return
	// }
	sids := this.getSidByUid(uid)
	for _, sid := range sids {
		session := this.getSession(sid)
		if session == nil {
			logger.Info("SendDataByUid session is nil sid:%s uid:%s", sid, uid)
		} else {
			session.Send(data)
		}

	}
}

func (this *SessionService) PushMsgByUids(uids []string, msg iface.IMessage) {
	// msgData, err := utils.MsgEncode(msg)
	// if err != nil {
	// 	logger.Error("PushMsg Error encode: ", err)
	// 	return
	// }

	// pkg := utils.BuildPackage(utils.PKG_DATA, len(msgData), msgData)
	// pkgData := utils.WritePackage(pkg)
	// counted := make(map[string]bool)
	// for _, uid := range uids {
	// 	if !counted[uid] {
	// 		this.SendDataByUid(uid, pkgData)
	// 		counted[uid] = true
	// 	}

	// }
}

type Session struct {
	Id         string
	FrontendId string
	Uid        string
	Setting    sync.Map

	socket         iface.Iconnection
	sessionService *SessionService
	state          string

	closeChan chan int
	timer     *time.Timer
}

func (this *Session) GetServerId() string {
	return this.FrontendId
}

func (this *Session) GetUid() string {
	return this.Uid
}

func (this *Session) GetId() string {
	return this.Id
}

func (this *Session) BackendSession() interface{} {
	setting := make(map[string]interface{})
	this.Setting.Range(func(k, v interface{}) bool {
		setting[k.(string)] = v
		return true
	})
	bSession := &iface.BackendSession{
		Uid:        this.Uid,
		FrontendId: this.FrontendId,
		Setting:    setting,
	}

	return bSession
}

func (this *Session) Bind(uid string) {
	if this.Uid != "" {
		this.sessionService.UnBind(this.Id, this.Uid)
	}
	this.Uid = uid
	this.sessionService.Bind(this.Id, uid)
	//emit bind event
}

func (this *Session) UnBind(uid string) {
	this.Uid = ""
	//emit unbind event
}

func (this *Session) Set(key string, value interface{}) {
	this.Setting.Store(key, value)
}

func (this *Session) Remove(key string) {
	this.Setting.Delete(key)
}

func (this *Session) Get(key string) interface{} {
	res, _ := this.Setting.Load(key)
	return res
}

func (this *Session) Send(data []byte) error {
	err := this.socket.Send(data)
	if err != nil {
		logger.Error("session.send failed sid：%s err %v", this.Id, err)
	}
	return err
}

func (this *Session) SendBatch(datas [][]byte) bool {
	return false
	//TODO return this.socket.send(msgs)
}

func (this *Session) Closed(reason string) {
	logger.Info("session on [%s] is closed with session id: %s %s", this.FrontendId, this.Id, this.socket.RemoteIp(), this.Uid, reason)
	if this.state == "closed" {
		return
	}
	this.state = "closed"
	close(this.closeChan)
	this.sessionService.Remove(this.Id)
}

func (this *Session) Kick(reason string) {
	logger.Info("session on [%s] is kicked with session id: %s %s", this.FrontendId, this.Id, this.socket.RemoteIp(), this.Uid, reason)
	this.socket.Stop()
}

func (this *Session) RemoteIp() string {
	return this.socket.RemoteIp()
}

var HtDur time.Duration = 6000 * time.Millisecond

func (this *Session) HeartBeat() {
	if this.timer == nil {
		this.timer = time.NewTimer(HtDur)
		go func() {
			select {
			case <-this.timer.C:
				this.Kick("heartbeat timeout")
			case <-this.closeChan:
				return
			}
		}()
	} else {
		this.timer.Reset(HtDur)
	}
}
