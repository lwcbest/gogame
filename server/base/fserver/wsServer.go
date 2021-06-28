package fserver

import (
	"fmt"
	"gameserver-997/server/base/fnet"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"

	"gameserver-997/server/base/timer"
	"gameserver-997/server/base/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

func init() {
	utils.GlobalObject.Protoc = fnet.NewProtocol()
	// --------------------------------------------init log start
	utils.ReSettingLog()
	// --------------------------------------------init log end
}

type WsServer struct {
	Name           string
	IPVersion      string
	IP             string
	Port           int
	MaxConn        int
	GenNum         *utils.UUIDGenerator
	ConnectionMgr  iface.Iconnectionmgr
	Protoc         iface.IServerProtocol
	SessionService iface.ISessionService
	ChannelService iface.IChannelService
	RootCluster    iface.Icluster
}

func (this *WsServer) SetService(s iface.ISessionService, c iface.IChannelService) {
	this.SessionService = s
	this.ChannelService = c
}

func (this *WsServer) handleConnection(conn *fnet.WsConn) {
	// conn.SetNoDelay(true)
	// conn.SetKeepAlive(true)
	// conn.SetDeadline(time.Now().Add(time.Minute * 2))

	var fconn *fnet.Connection
	if this.Protoc == nil {
		fconn = fnet.NewConnection(conn, this.GenNum.GetUint32(), utils.GlobalObject.Protoc)
	} else {
		fconn = fnet.NewConnection(conn, this.GenNum.GetUint32(), this.Protoc)
	}
	fconn.SetProperty(fnet.XINGO_CONN_PROPERTY_NAME, this.Name)
	fconn.Start()
}

func (this *WsServer) Start() {
	// todo handle websocket by listen http
	utils.GlobalObject.TcpServers[this.Name] = this
	go func() {
		this.Protoc.InitWorker(utils.GlobalObject.PoolSize)
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			conn, err := (&websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}).Upgrade(w, r, nil)
			if err != nil {
				http.NotFound(w, r)
				logger.Info("upgrade failed: ", err)
				return
			}
			conn.SetCloseHandler(func(code int, msg string) error {
				logger.Info("wsConn is closed code %d msg %s", code, msg)
				return nil
			})

			ip := r.Header.Get("X-Real_IP")
			wsConn := fnet.NewWsConn(conn, ip)
			if this.ConnectionMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
			} else {
				go this.handleConnection(wsConn)
			}
		})
		http.ListenAndServe(fmt.Sprintf("%s:%d", this.IP, this.Port), nil)
	}()
}

func (this *WsServer) GetName() string {
	return this.Name
}

func (this *WsServer) GetSessionService() iface.ISessionService {
	return this.SessionService
}

func (this *WsServer) GetChannelService() iface.IChannelService {
	return this.ChannelService
}

func (this *WsServer) GetConnectionMgr() iface.Iconnectionmgr {
	return this.ConnectionMgr
}

func (this *WsServer) GetConnectionQueue() chan interface{} {
	return nil
}

func (this *WsServer) Stop() {
	logger.Info("stop xingo server ", this.Name)
	if utils.GlobalObject.OnServerStop != nil {
		utils.GlobalObject.OnServerStop()
	}
}

func (this *WsServer) AddRouter(router interface{}) {
	logger.Info("AddRouter")
	utils.GlobalObject.Protoc.GetMsgHandle().AddRouter(router)
}

func (this *WsServer) CallLater(durations time.Duration, f func(v ...interface{}), args ...interface{}) {
	delayTask := timer.NewTimer(durations, f, args)
	delayTask.Run()
}

func (this *WsServer) CallWhen(ts string, f func(v ...interface{}), args ...interface{}) {
	loc, err_loc := time.LoadLocation("Local")
	if err_loc != nil {
		logger.Error("CallWhen err %v", err_loc)
		return
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", ts, loc)
	now := time.Now()
	if err == nil {
		if now.Before(t) {
			this.CallLater(t.Sub(now), f, args...)
		} else {
			logger.Error("CallWhen time before now")
		}
	} else {
		logger.Error("callwhen err %v", err)
	}
}

func (this *WsServer) CallLoop(durations time.Duration, f func(v ...interface{}), args ...interface{}) {
	go func() {
		delayTask := timer.NewTimer(durations, f, args)
		for {
			time.Sleep(delayTask.GetDurations())
			delayTask.GetFunc().Call()
		}
	}()
}

func (this *WsServer) WaitSignal() {
	signal.Notify(utils.GlobalObject.ProcessSignalChan, os.Kill, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	sig := <-utils.GlobalObject.ProcessSignalChan
	logger.Info(fmt.Sprintf("server exit. signal: [%s]", sig))
	this.Stop()
}

func (this *WsServer) Serve() {
	this.Start()
	this.WaitSignal()
}
