package clusterserver

import (
	"errors"
	"fmt"
	"gameserver-997/server/base/cluster"
	"gameserver-997/server/base/fnet"
	"gameserver-997/server/base/fserver"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/base/service"
	"gameserver-997/server/base/utils"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"
)

type ClusterServer struct {
	Name           string
	RemoteNodesMgr *cluster.ChildMgr //子节点有
	ChildsMgr      *cluster.ChildMgr //root节点有
	MasterObj      *fnet.TcpClient
	httpServerMux  *http.ServeMux
	NetServer      iface.Iserver
	RootServer     iface.Iserver
	TelnetServer   iface.Iserver
	Cconf          *cluster.ClusterConf
	modules        map[string][]interface{} //所有模块统一管理
	sync.RWMutex
}

func DoCSConnectionLost(fconn iface.Iconnection) {
	logger.Error("node disconnected from " + utils.GlobalObject.Name)
	//子节点掉线
	nodename, err := fconn.GetProperty("child")
	if err == nil {
		GlobalClusterServer.RemoveChild(nodename.(string))
	}
}

func DoCCConnectionLost(fconn iface.Iclient) {
	//父节点掉线
	rname, err := fconn.GetProperty("remote")
	if err == nil {
		GlobalClusterServer.RemoveRemote(rname.(string))
		logger.Error("remote " + rname.(string) + " disconnected from " + utils.GlobalObject.Name)
	}
}

//reconnected to master
func ReConnectMasterCB(fconn iface.Iclient) {
	rpc := cluster.NewChild(utils.GlobalObject.Name, GlobalClusterServer.MasterObj)
	response, err := rpc.CallChildForResult("TakeProxy", utils.GlobalObject.Name)
	if err == nil {
		roots, ok := response.Result["roots"]
		if ok {
			for _, root := range roots.([]string) {
				GlobalClusterServer.ConnectToRemote(root)
			}
		}
	} else {
		panic(fmt.Sprintf("reconnected to master error: %s", err))
	}
}

func NewClusterServer(name, path string) *ClusterServer {
	// logger.SetPrefix(fmt.Sprintf("[%s]", strings.ToUpper(name)))
	cconf, err := cluster.NewClusterConf(path)
	if err != nil {
		panic("cluster conf error!!!")
	}

	GlobalClusterServer = &ClusterServer{
		Name:           name,
		Cconf:          cconf,
		RemoteNodesMgr: cluster.NewChildMgr(),
		ChildsMgr:      cluster.NewChildMgr(),
		modules:        make(map[string][]interface{}, 0),
		httpServerMux:  http.NewServeMux(),
	}

	serverconf, ok := GlobalClusterServer.Cconf.Servers[name]
	if !ok {
		panic(fmt.Sprintf("no server %s in clusterconf!!!", name))
	}

	utils.GlobalObject.Name = name
	utils.GlobalObject.OnClusterClosed = DoCSConnectionLost
	utils.GlobalObject.OnClusterCClosed = DoCCConnectionLost
	utils.GlobalObject.RpcCProtoc = cluster.NewRpcClientProtocol()

	if utils.GlobalObject.PoolSize > 0 {
		//init rpc worker pool
		utils.GlobalObject.RpcCProtoc.InitWorker(int32(utils.GlobalObject.PoolSize))
	}
	if serverconf.NetPort > 0 {
		utils.GlobalObject.Protoc = fnet.NewProtocol()
	}
	if serverconf.RootPort > 0 {
		utils.GlobalObject.RpcSProtoc = cluster.NewRpcServerProtocol()
	}

	if serverconf.Log != "" {
		utils.GlobalObject.LogName = serverconf.Log
		utils.ReSettingLog()
	}

	//telnet debug tool
	if serverconf.DebugPort > 0 {
		if serverconf.Host != "" {
			GlobalClusterServer.TelnetServer = newTcpServer(GlobalClusterServer, "telnet_server", "tcp4", serverconf.Host, serverconf.DebugPort, 100, cluster.NewTelnetProtocol())
		} else {
			GlobalClusterServer.TelnetServer = newTcpServer(GlobalClusterServer, "telnet_server", "tcp4", "127.0.0.1", serverconf.DebugPort, 100, cluster.NewTelnetProtocol())
		}
	}
	return GlobalClusterServer
}

func (this *ClusterServer) StartClusterServer() {
	serverconf, ok := this.Cconf.Servers[utils.GlobalObject.Name]
	if !ok {
		panic("no server in clusterconf!!!")
	}
	//自动发现注册modules api
	modules, ok := this.modules[serverconf.Module]
	if ok {
		//api
		if serverconf.NetPort > 0 {
			for _, m := range modules[0].([]interface{}) {
				if m != nil {
					this.AddRouter(m)
				}
			}
		}
		//http
		if len(serverconf.Http) > 0 || len(serverconf.Https) > 0 {
			for _, m := range modules[1].([]interface{}) {
				if m != nil {
					this.AddHttpRouter(m)
				}
			}
		}
		//rpc
		for _, m := range modules[2].([]interface{}) {
			if m != nil {
				this.AddRpcRouter(m)
			}
		}
	}

	//http server
	if len(serverconf.Http) > 0 {
		//staticfile handel
		if len(serverconf.Http) == 2 {
			this.httpServerMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(serverconf.Http[1].(string)))))
		}
		httpserver := &http.Server{
			Addr:           fmt.Sprintf(":%d", int(serverconf.Http[0].(float64))),
			Handler:        this.httpServerMux,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 1 << 20, //1M
		}
		httpserver.SetKeepAlivesEnabled(true)
		go httpserver.ListenAndServe()
		logger.Info(fmt.Sprintf("http://%s:%d start", serverconf.Host, int(serverconf.Http[0].(float64))))
	} else if len(serverconf.Https) > 2 {
		//staticfile handel
		if len(serverconf.Https) == 4 {
			this.httpServerMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(serverconf.Https[3].(string)))))
		}
		httpserver := &http.Server{
			Addr:           fmt.Sprintf(":%d", int(serverconf.Https[0].(float64))),
			Handler:        this.httpServerMux,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 1 << 20, //1M
		}
		httpserver.SetKeepAlivesEnabled(true)
		go httpserver.ListenAndServeTLS(serverconf.Https[1].(string), serverconf.Https[2].(string))
		logger.Info(fmt.Sprintf("http://%s:%d start", serverconf.Host, int(serverconf.Https[0].(float64))))
	}
	//tcp server or ws server
	if serverconf.UseWs {
		if serverconf.NetPort > 0 {
			utils.GlobalObject.TcpPort = serverconf.NetPort
			this.NetServer = newWsServer(this, "xingocluster_ws_server", "ws", serverconf.Host, serverconf.NetPort,
				utils.GlobalObject.MaxConn, utils.GlobalObject.Protoc)
			this.NetServer.Start()
		}
	} else if serverconf.NetPort > 0 {
		utils.GlobalObject.TcpPort = serverconf.NetPort
		if serverconf.Host != "" {
			this.NetServer = newTcpServer(this, "xingocluster_net_server", "tcp4", serverconf.Host, serverconf.NetPort,
				utils.GlobalObject.MaxConn, utils.GlobalObject.Protoc)
		} else {
			this.NetServer = newTcpServer(this, "xingocluster_net_server", "tcp4", serverconf.Host, serverconf.NetPort,
				utils.GlobalObject.MaxConn, utils.GlobalObject.Protoc)
		}
		this.NetServer.Start()
	}
	if serverconf.RootPort > 0 {
		if serverconf.Host != "" {
			this.RootServer = newTcpServer(this, "xingocluster_root_server", "tcp4", serverconf.Host, serverconf.RootPort,
				utils.GlobalObject.IntraMaxConn, utils.GlobalObject.RpcSProtoc)
		} else {
			this.RootServer = newTcpServer(this, "xingocluster_root_server", "tcp4", serverconf.Host, serverconf.RootPort,
				utils.GlobalObject.IntraMaxConn, utils.GlobalObject.RpcSProtoc)
		}
		this.RootServer.Start()
	}
	//telnet
	if this.TelnetServer != nil {
		logger.Info(fmt.Sprintf("telnet tool start: %s:%d.", serverconf.Host, serverconf.DebugPort))
		this.TelnetServer.Start()
	}

	//master
	this.ConnectToMaster()

	logger.Info("xingo cluster start success.")
	// close
	this.WaitSignal()
	this.MasterObj.Stop(true)
	if this.RootServer != nil {
		this.RootServer.Stop()
	}

	if this.NetServer != nil {
		this.NetServer.Stop()
	}

	if this.TelnetServer != nil {
		this.TelnetServer.Stop()
	}
	logger.Info("xingo cluster stoped.")
}

func (this *ClusterServer) WaitSignal() {
	signal.Notify(utils.GlobalObject.ProcessSignalChan, os.Kill, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	sig := <-utils.GlobalObject.ProcessSignalChan
	//尝试主动通知master checkalive
	rpc := cluster.NewChild(utils.GlobalObject.Name, this.MasterObj)
	rpc.CallChildNotForResult("ChildOffLine", utils.GlobalObject.Name)

	logger.Info(fmt.Sprintf("server exit. signal: [%s]", sig))
}

func (this *ClusterServer) GetServerAddr(serverType string) ([]string, error) {
	rpc := cluster.NewChild(utils.GlobalObject.Name, this.MasterObj)
	rpcData, err := rpc.CallChildForResult("OnlineNodes", serverType)
	if err != nil {
		return nil, err
	}
	urls := rpcData.Result["urls"]
	if urls != nil {
		return urls.([]string), nil
	}
	return nil, nil
}

func (this *ClusterServer) ConnectToMaster() {
	master := fnet.NewReConnTcpClient(this.Cconf.Master.Host, this.Cconf.Master.RootPort, utils.GlobalObject.RpcCProtoc, 1024, 60, ReConnectMasterCB)
	this.MasterObj = master
	master.Start()
	//注册到master
	rpc := cluster.NewChild(utils.GlobalObject.Name, this.MasterObj)
	response, err := rpc.CallChildForResult("TakeProxy", utils.GlobalObject.Name)
	if err == nil {
		roots, ok := response.Result["roots"]
		logger.Info("ConnectToRemote -------", utils.GlobalObject.Name, response.Result)
		if ok {
			for _, root := range roots.([]string) {
				this.ConnectToRemote(root)
			}
		}
	} else {
		panic(fmt.Sprintf("connected to master error: %s", err))
	}
}

func (this *ClusterServer) ConnectToRemote(rname string) {
	rserverconf, ok := this.Cconf.Servers[rname]
	if ok {
		//处理master掉线，重新通知的情况
		if _, err := this.GetRemote(rname); err != nil {
			rserver := fnet.NewTcpClient(rserverconf.Host, rserverconf.RootPort, utils.GlobalObject.RpcCProtoc)
			this.RemoteNodesMgr.AddChild(rname, rserver)
			rserver.Start()
			rserver.SetProperty("remote", rname)
			//takeproxy
			child, err := this.RemoteNodesMgr.GetChild(rname)
			if err == nil {
				child.CallChildNotForResult("TakeProxy", utils.GlobalObject.Name)
			}
		} else {
			logger.Info("Remote connection already exist!")
		}
	} else {
		//未找到节点
		logger.Error("ConnectToRemote error. " + rname + " node can`t found!!!")
	}
}

func (this *ClusterServer) AddRouter(router interface{}) {
	if utils.GlobalObject.Protoc != nil {
		utils.GlobalObject.Protoc.AddRpcRouter(router)
	}
}

func (this *ClusterServer) AddRpcRouter(router interface{}) {
	utils.GlobalObject.RpcCProtoc.AddRpcRouter(router)
	if utils.GlobalObject.RpcSProtoc != nil {
		utils.GlobalObject.RpcSProtoc.AddRpcRouter(router)
	}
}

/*
子节点连上来回调
*/
func (this *ClusterServer) AddChild(name string, writer iface.IWriter) {
	this.Lock()
	defer this.Unlock()

	this.ChildsMgr.AddChild(name, writer)
	writer.SetProperty("child", name)
}

/*
子节点断开回调
*/
func (this *ClusterServer) RemoveChild(name string) {
	this.Lock()
	defer this.Unlock()

	this.ChildsMgr.RemoveChild(name)
}

func (this *ClusterServer) RemoveRemote(name string) {
	this.Lock()
	defer this.Unlock()

	this.RemoteNodesMgr.RemoveChild(name)
}

func (this *ClusterServer) GetRemote(name string) (*cluster.Child, error) {
	this.RLock()
	defer this.RUnlock()

	return this.RemoteNodesMgr.GetChild(name)
}

/*
注册模块到分布式服务器
*/
func (this *ClusterServer) AddModule(mname string, apimodule interface{}, httpmodule interface{}, rpcmodule interface{}) {
	//this.modules[mname] = []interface{}{module, rpcmodule}
	if _, ok := this.modules[mname]; ok {
		this.modules[mname][0] = append(this.modules[mname][0].([]interface{}), apimodule)
		this.modules[mname][1] = append(this.modules[mname][1].([]interface{}), httpmodule)
		this.modules[mname][2] = append(this.modules[mname][2].([]interface{}), rpcmodule)
	} else {
		this.modules[mname] = []interface{}{[]interface{}{apimodule}, []interface{}{httpmodule}, []interface{}{rpcmodule}}
	}
}

/*
注册http的api到分布式服务器
*/
func (this *ClusterServer) AddHttpRouter(router interface{}) {
	value := reflect.ValueOf(router)
	tp := value.Type()
	for i := 0; i < value.NumMethod(); i += 1 {
		name := tp.Method(i).Name
		uri := fmt.Sprintf("/%s", strings.ToLower(strings.Replace(name, "Handle", "", 1)))
		this.httpServerMux.HandleFunc(uri,
			utils.HttpRequestWrap(uri, value.Method(i).Interface().(func(http.ResponseWriter, *http.Request))))
		logger.Info("add http url: " + uri)
	}
}

func (this *ClusterServer) RpcCallServerId(session iface.ISession, serverId string, target string, args ...interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	tarServer, err := GlobalClusterServer.RemoteNodesMgr.GetChild(serverId)
	if err != nil {
		return response, err
	}
	var newArgs []interface{}
	if session != nil {
		newArgs = append([]interface{}{session.BackendSession()}, args...)
	} else {
		newArgs = append([]interface{}{args[0]}, args[1:]...)
	}
	rpcData, err := tarServer.CallChildForResult(target, newArgs...)
	if err != nil {
		return response, err
	} else if err := rpcData.Result["err"]; err != nil {
		return response, errors.New(err.(string))
	}
	return rpcData.Result, nil
}

func (this *ClusterServer) RpcRandomCallServerType(session iface.ISession, serverType string, target string, args ...interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var tarServer *cluster.Child
	var err error

	tarServer = GlobalClusterServer.RemoteNodesMgr.GetRandomChild(serverType)

	if tarServer == nil {
		return response, errors.New(fmt.Sprintf("target server by type %s not found", serverType))
	}
	var newArgs []interface{}
	if session != nil {
		newArgs = append([]interface{}{session.BackendSession()}, args...)
	} else {
		newArgs = append([]interface{}{args[0]}, args[1:]...)
	}
	rpcData, err := tarServer.CallChildForResult(target, newArgs...)
	if err != nil {
		return response, err
	} else if err := rpcData.Result["err"]; err != nil {
		return response, errors.New(err.(string))
	}
	rpcData.Result["serverName"] = tarServer.GetName()
	return rpcData.Result, nil
}

func (this *ClusterServer) RpcSystemCallServerType(session iface.ISession, serverType string, target string, args ...interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	var tarServer *cluster.Child
	var err error
	preferKey := fmt.Sprintf("SERVER_PREFER:%s", serverType)
	if serverName := session.Get(preferKey); serverName != nil {
		tarServer, err = GlobalClusterServer.RemoteNodesMgr.GetChild(serverName.(string))
	}
	if tarServer == nil {
		tarServer = GlobalClusterServer.RemoteNodesMgr.GetRandomChild(serverType)
	}
	if tarServer == nil {
		return response, errors.New(fmt.Sprintf("target server by type %s not found", serverType))
	}
	var newArgs []interface{}
	if session != nil {
		newArgs = append([]interface{}{session.BackendSession()}, args...)
	} else {
		newArgs = append([]interface{}{args[0]}, args[1:]...)
	}
	rpcData, err := tarServer.CallChildForResult(target, newArgs...)
	if err != nil {
		return response, err
	} else if err := rpcData.Result["err"]; err != nil {
		return response, errors.New(err.(string))
	}
	return rpcData.Result, nil
}

func (this *ClusterServer) RpcPushServerId(session iface.ISession, serverId string, target string, args ...interface{}) error {
	tarServer, err := GlobalClusterServer.RemoteNodesMgr.GetChild(serverId)
	if err != nil {
		return err
	}
	
	var newArgs []interface{}
	if session != nil {
		newArgs = append([]interface{}{session.BackendSession()}, args...)
	} else {
		newArgs = append([]interface{}{args[0]}, args[1:]...)
	}

	err = tarServer.CallChildNotForResult(target, newArgs...)
	if err != nil {
		return err
	}
	return nil
}

func (this *ClusterServer) RpcPushServerName(serverName string, target string, args ...interface{}) error {
	var tarServer *cluster.Child
	var err error
	tarServer = GlobalClusterServer.RemoteNodesMgr.GetRandomChild(serverName)
	if tarServer == nil {
		return errors.New(fmt.Sprintf("target server by type %s not found", serverName))
	}
	err = tarServer.CallChildNotForResult(target, args...)
	if err != nil {
		return err
	}
	return nil
}

func newWsServer(cluster *ClusterServer, name string, version string, ip string, port int, maxConn int, protoc iface.IServerProtocol) iface.Iserver {
	wss := &fserver.WsServer{
		Name:          name,
		IPVersion:     version,
		IP:            ip,
		Port:          port,
		MaxConn:       maxConn,
		ConnectionMgr: fnet.NewConnectionMgr(),
		Protoc:        protoc,
		GenNum:        utils.NewUUIDGenerator(name),
		RootCluster:   cluster,
	}
	ss := &service.SessionService{RootCluster: cluster, RootServer: wss}
	cs := &service.ChannelService{RootCluster: cluster, RootServer: wss}
	wss.SetService(ss, cs)
	utils.GlobalObject.TcpServer = wss

	return wss
}

func newTcpServer(cluster *ClusterServer, name string, version string, ip string, port int, maxConn int, protoc iface.IServerProtocol) iface.Iserver {
	s := &fserver.Server{
		Name:          name,
		IPVersion:     version,
		IP:            ip,
		Port:          port,
		MaxConn:       maxConn,
		ConnectionMgr: fnet.NewConnectionMgr(),
		Protoc:        protoc,
		GenNum:        utils.NewUUIDGenerator(name),
	}

	ss := &service.SessionService{RootCluster: cluster, RootServer: s}
	cs := &service.ChannelService{RootCluster: cluster, RootServer: s}
	s.SetService(ss, cs)
	utils.GlobalObject.TcpServer = s

	return s
}
