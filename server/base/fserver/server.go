package fserver

import (
	"fmt"
	"gameserver-997/server/base/fnet"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/base/timer"
	"gameserver-997/server/base/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	utils.GlobalObject.Protoc = fnet.NewProtocol()
	// --------------------------------------------init log start
	utils.ReSettingLog()
	// --------------------------------------------init log end
}

type Server struct {
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

func NewServer() iface.Iserver {
	s := &Server{
		Name:          utils.GlobalObject.Name,
		IPVersion:     "tcp4",
		IP:            "0.0.0.0",
		Port:          utils.GlobalObject.TcpPort,
		MaxConn:       utils.GlobalObject.MaxConn,
		ConnectionMgr: fnet.NewConnectionMgr(),
		Protoc:        utils.GlobalObject.Protoc,
		GenNum:        utils.NewUUIDGenerator(""),
	}
	utils.GlobalObject.TcpServer = s

	return s
}

func NewTcpServer(name string, version string, ip string, port int, maxConn int, protoc iface.IServerProtocol) iface.Iserver {
	s := &Server{
		Name:          name,
		IPVersion:     version,
		IP:            ip,
		Port:          port,
		MaxConn:       maxConn,
		ConnectionMgr: fnet.NewConnectionMgr(),
		Protoc:        protoc,
		GenNum:        utils.NewUUIDGenerator(name),
	}
	utils.GlobalObject.TcpServer = s

	return s
}

func (this *Server) SetService(s iface.ISessionService, c iface.IChannelService) {
	this.SessionService = s
	this.ChannelService = c
}

func (this *Server) handleConnection(conn *net.TCPConn) {
	conn.SetNoDelay(true)
	conn.SetKeepAlive(true)
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

func (this *Server) Start() {
	utils.GlobalObject.TcpServers[this.Name] = this
	go func() {
		this.Protoc.InitWorker(utils.GlobalObject.PoolSize)
		tcpAddr, err := net.ResolveTCPAddr(this.IPVersion, fmt.Sprintf("%s:%d", this.IP, this.Port))
		if err != nil {
			logger.Fatal("ResolveTCPAddr err: ", err)
			return
		}
		ln, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			logger.Error("ListenTCP failed addr:%v, err:%v", tcpAddr, err)
			panic(err)
		}
		logger.Info(fmt.Sprintf("start xingo server %s...", this.Name))
		for {
			conn, err := ln.AcceptTCP()
			if err != nil {
				logger.Error("AcceptTCP err %v", err)
				continue
			}
			//max client exceed
			if this.ConnectionMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
			} else {
				go this.handleConnection(conn)
			}
		}
	}()
}

func (this *Server) GetName() string {
	return this.Name
}

func (this *Server) GetConnectionMgr() iface.Iconnectionmgr {
	return this.ConnectionMgr
}

func (this *Server) GetConnectionQueue() chan interface{} {
	return nil
}

func (this *Server) GetSessionService() iface.ISessionService {
	return this.SessionService
}

func (this *Server) GetChannelService() iface.IChannelService {
	return this.ChannelService
}

func (this *Server) Stop() {
	logger.Info("stop xingo server ", this.Name)
	if utils.GlobalObject.OnServerStop != nil {
		utils.GlobalObject.OnServerStop()
	}
}

func (this *Server) AddRouter(router interface{}) {
	logger.Info("AddRouter")
	utils.GlobalObject.Protoc.GetMsgHandle().AddRouter(router)
}

func (this *Server) CallLater(durations time.Duration, f func(v ...interface{}), args ...interface{}) {
	delayTask := timer.NewTimer(durations, f, args)
	delayTask.Run()
}

func (this *Server) CallWhen(ts string, f func(v ...interface{}), args ...interface{}) {
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
		logger.Error("CallWhen err %v", err)
	}
}

func (this *Server) CallLoop(durations time.Duration, f func(v ...interface{}), args ...interface{}) {
	go func() {
		delayTask := timer.NewTimer(durations, f, args)
		for {
			time.Sleep(delayTask.GetDurations())
			delayTask.GetFunc().Call()
		}
	}()
}

func (this *Server) WaitSignal() {
	signal.Notify(utils.GlobalObject.ProcessSignalChan, os.Kill, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	sig := <-utils.GlobalObject.ProcessSignalChan
	logger.Info(fmt.Sprintf("server exit. signal: [%s]", sig))
	this.Stop()
}

func (this *Server) Serve() {
	this.Start()
	this.WaitSignal()
}
