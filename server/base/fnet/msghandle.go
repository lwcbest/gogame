package fnet

/*
	find msg api
*/
import (
	"fmt"
	"gameserver-997/pb/gopb"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/base/utils"
	"gameserver-997/server/constants"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"

	// "errors"
	"time"
	// "math"
	"runtime/debug"
)

type MsgHandle struct {
	PoolSize  int32
	TaskQueue []chan *utils.Package
	Apis      map[string]reflect.Value
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		PoolSize:  utils.GlobalObject.PoolSize,
		TaskQueue: make([]chan *utils.Package, utils.GlobalObject.PoolSize),
		Apis:      make(map[string]reflect.Value),
	}
}

//一致性路由,保证同一连接的数据转发给相同的goroutine
func (this *MsgHandle) DeliverToMsgQueue(v interface{}) {
	pkg := v.(*utils.Package)
	index := pkg.GetConnection().GetSessionId() % uint32(utils.GlobalObject.PoolSize)
	taskQueue := this.TaskQueue[index]
	taskQueue <- pkg
}

func (this *MsgHandle) DoMsgFromGoRoutine(v interface{}) {
	pkg := v.(*utils.Package)
	go func() {
		this.HandlePackage(pkg)
	}()
}

func (this *MsgHandle) AddRouter(router interface{}) {
	value := reflect.ValueOf(router)
	tp := value.Type()
	for i := 0; i < value.NumMethod(); i += 1 {
		name := tp.Method(i).Name
		k := strings.Split(name, "_")

		if _, ok := this.Apis[k[0]]; ok {
			panic("repeated api " + k[0])
		}
		this.Apis[k[0]] = value.Method(i)
		logger.Info("add api " + name)
	}
}

func (this *MsgHandle) HandleCommonError(req *iface.CommonRequest, err interface{}) {
	if err != nil {
		logger.Error("msg error:%+v", err)
		debug.PrintStack()
	}
	errResp := &gopb.ResError{Code: constants.MSG_CODE.FAIL, Msg: "Server Common Error"}
	replyMsg := utils.BuildMsg(utils.MSG_RESPONSE, req.ReqId, req.Router, errResp)
	msgData, _ := utils.MsgEncode(replyMsg) //todo error检测
	replyPkg := utils.BuildPkg(utils.PKG_DATA, msgData)
	replyData := utils.WritePackage(replyPkg)
	req.Session.Send(replyData)
}

func (this *MsgHandle) HandleMessage(session iface.ISession, body []byte) {
	msg, err := utils.MsgDecode(body)
	if err != nil {
		logger.Error("decodeMsgError: ", err)
		return
	}
	switch msg.GetMsgType() {
	case utils.MSG_NOTIFY, utils.MSG_PUSH, utils.MSG_RESPONSE:
		logger.Info("received a notify or push or response message msg: %v", msg)
	case utils.MSG_REQUEST:
		//根据服务器名判断 e.g. Net_Login/Game_Fight
		routes := strings.Split(msg.GetRoute(), "_")
		tarServer := strings.ToLower(routes[0])

		if len(routes) <= 1 {
			//只有一个啥也不做防止报错
		} else if len(tarServer) > len(utils.GlobalObject.Name) || tarServer != utils.GlobalObject.Name[:len(tarServer)] {
			//RPC 路由
			routes = []string{"Proxy", tarServer, routes[1]}
		} else {
			routes = routes[1:]
		}

		if f, ok := this.Apis[routes[0]]; ok {
			//存在
			st := time.Now()
			var replyMsg *utils.Message
			if len(routes) == 1 {
				//前端接口
				r := &iface.CommonRequest{Session: session, ReqId: msg.GetMsgId(), Router: msg.GetRoute(), Data: msg.GetData(), RouterStr: routes[0]}

				//1. Dd before Filter
				beforeErr := utils.GlobalObject.MsgFiltersObj.DoBeforeActions(r)
				if beforeErr != nil {
					this.HandleCommonError(r, beforeErr)
					return
				}
				//2. Do handler
				values, handleError := utils.XingoTry(f, []reflect.Value{reflect.ValueOf(r)})
				if handleError != nil {
					this.HandleCommonError(r, handleError)
					return
				}
				respProtoMsg := values[0].Interface().(proto.Message)
				//3. Do after filter
				afterErr := utils.GlobalObject.MsgFiltersObj.DoAfterActions(r, respProtoMsg)
				if beforeErr != nil {
					this.HandleCommonError(r, afterErr)
					return
				}

				replyMsg = utils.BuildMsg(utils.MSG_RESPONSE, msg.GetMsgId(), msg.GetRoute(), respProtoMsg)
			} else {
				//RPC后端接口 TODO 这块要重构
				r := &iface.CommonRequest{Session: session, ReqId: msg.GetMsgId(), Router: msg.GetRoute(), Data: msg.GetData(), RouterStr: routes[0]}
				values, handleError := utils.XingoTry(f, []reflect.Value{reflect.ValueOf(routes[1]), reflect.ValueOf(routes[2]), reflect.ValueOf(session), reflect.ValueOf(msg.GetData())})
				if handleError != nil {

					this.HandleCommonError(r, handleError)
					return
				}

				replyMsg = utils.BuildMsgFromData(utils.MSG_RESPONSE, msg.GetMsgId(), msg.GetRoute(), values[0].Interface().([]byte))
			}
			logger.Debug(fmt.Sprintf("Api_%s cost total time: %f ms", msg.GetRoute(), time.Now().Sub(st).Seconds()*1000))

			msgData, _ := utils.MsgEncode(replyMsg) //todo error检测
			replyPkg := utils.BuildPkg(utils.PKG_DATA, msgData)
			replyData := utils.WritePackage(replyPkg)
			session.Send(replyData)

		} else {
			logger.Error(fmt.Sprintf("not found api:  %s ", msg.GetRoute()))
		}
	}
}

// 根据packageType做成不同的处理
func (this *MsgHandle) HandlePackage(pkg *utils.Package) {
	sessionService := utils.GlobalObject.TcpServer.GetSessionService()
	if sessionService == nil {
		logger.Error("[HandlePackage]error, no session service~")
	}

	serverName := utils.GlobalObject.Name
	sessionId := fmt.Sprintf("%s:%d", serverName, pkg.GetConnection().GetSessionId())
	switch pkg.GetPkgType() {
	case utils.PKG_HANDSHAKE:
		session := sessionService.Create(sessionId, serverName, pkg.GetConnection())
		ackPkg := utils.BuildPkg(utils.PKG_HANDSHAKE_ACK, nil)
		session.Send(utils.WritePackage(ackPkg))
		session.HeartBeat()
	case utils.PKG_HANDSHAKE_ACK, utils.PKG_HEARTBEAT:
		connection := pkg.GetConnection()
		session := sessionService.Get(sessionId)
		if session == nil {
			connection.Stop()
		} else {
			ackPkg := utils.BuildPkg(utils.PKG_HEARTBEAT, nil)
			session.HeartBeat()
			session.Send(utils.WritePackage(ackPkg))
		}
	case utils.PKG_DATA:
		connection := pkg.GetConnection()
		session := sessionService.Get(sessionId)
		if session == nil {
			connection.Stop()
		} else {
			this.HandleMessage(session, pkg.GetData())
		}
	case utils.PKG_KICK:
		logger.Info("kick pkg %d")
	}
}

func (this *MsgHandle) StartWorkerLoop(poolSize int) {
	if utils.GlobalObject.IsThreadSafeMode() {
		//线程安全模式所有的逻辑都在一个goroutine处理, 这样可以实现无锁化服务
		// this.TaskQueue[0] = make(chan *PkgAll, utils.GlobalObject.MaxWorkerLen)
		this.TaskQueue[0] = make(chan *utils.Package, utils.GlobalObject.MaxWorkerLen)
		go func() {
			logger.Info("init thread mode workpool.")
			for {
				select {
				case data := <-this.TaskQueue[0]:
					this.HandlePackage(data)
				case delaytask := <-utils.GlobalObject.GetSafeTimer().GetTriggerChannel():
					delaytask.Call()
				}
			}
		}()
	} else {
		for i := 0; i < poolSize; i += 1 {
			// c := make(chan *PkgAll, utils.GlobalObject.MaxWorkerLen)
			c := make(chan *utils.Package, utils.GlobalObject.MaxWorkerLen)
			this.TaskQueue[i] = c
			// go func(index int, taskQueue chan *PkgAll) {
			go func(index int, taskQueue chan *utils.Package) {
				// logger.Info(fmt.Sprintf("init thread pool %d.", index))
				for {
					data := <-taskQueue
					this.HandlePackage(data)
				}
			}(i, c)
		}
	}
}
