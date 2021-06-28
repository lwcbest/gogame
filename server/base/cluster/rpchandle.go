package cluster

/*
	regest rpc
*/
import (
	"fmt"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/base/utils"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

type RpcMsgHandle struct {
	PoolSize  int32
	TaskQueue []chan *RpcRequest
	Apis      map[string]reflect.Value
}

func NewRpcMsgHandle() *RpcMsgHandle {
	return &RpcMsgHandle{
		PoolSize:  utils.GlobalObject.PoolSize,
		TaskQueue: make([]chan *RpcRequest, utils.GlobalObject.PoolSize),
		Apis:      make(map[string]reflect.Value),
	}
}

/*
处理rpc消息
*/
func (this *RpcMsgHandle) DoMsg(request *RpcRequest) {
	if request.Rpcdata.MsgType == RESPONSE && request.Rpcdata.Key != "" {
		//放回异步结果
		AResultGlobalObj.FillAsyncResult(request.Rpcdata.Key, request.Rpcdata)
		return
	} else {
		//rpc 请求
		if f, ok := this.Apis[request.Rpcdata.Target]; ok {
			//存在
			var r *iface.CommonRequest = &iface.CommonRequest{Fconn: request.Fconn}
			var bSession *iface.BackendSession
			if len(request.Rpcdata.Args) > 0 {
				bSession, ok = request.Rpcdata.Args[0].(*iface.BackendSession)
			}
			if ok {
				r.BSession = bSession
				if len(request.Rpcdata.Args) > 1 {
					r.RpcData = request.Rpcdata.Args[1:]
				}
			} else {
				r.RpcData = request.Rpcdata.Args
			}
			st := time.Now()
			defer func() {
				if err := recover(); err != nil {
					logger.Error("[rpc][error][%v]: %v", request.Rpcdata.Target, err)
					errorStack := ""
					for i := 1; i < 5; i++ {
						pc, file, line, ok := runtime.Caller(i)
						if !ok {
							break
						}
						errorStack += fmt.Sprintf("%s:%d (0x%x) && ", file, line, pc)
					}
					logger.Error("[rpc][error][%v]: %v", request.Rpcdata.Target, errorStack)
				}
			}()
			if request.Rpcdata.MsgType == REQUEST_FORRESULT {
				logger.Debug("[REQUEST_FORRESULT]:%v", request.Rpcdata.Target)
				ret := f.Call([]reflect.Value{reflect.ValueOf(r)})
				if len(ret) == 0 {
					logger.Error("RpcRequest no resp ", request.Rpcdata.Target)
					return
				}
				packdata, err := utils.GlobalObject.RpcCProtoc.GetDataPack().Pack(0, &RpcData{
					MsgType: RESPONSE,
					Result:  ret[0].Interface().(map[string]interface{}),
					Key:     request.Rpcdata.Key,
				})
				if err == nil {
					request.Fconn.Send(packdata)
				} else {
					logger.Error("Pack rpcData failed err: %v", err)
				}
			} else if request.Rpcdata.MsgType == REQUEST_NORESULT {
				// f.Call([]reflect.Value{reflect.ValueOf(request)})
				f.Call([]reflect.Value{reflect.ValueOf(r)})
			}

			logger.Debug(fmt.Sprintf("rpc %s cost total time: %f ms", request.Rpcdata.Target, time.Now().Sub(st).Seconds()*1000))
		} else {
			logger.Error(fmt.Sprintf("not found rpc:  %s", request.Rpcdata.Target))
		}
	}
}

func (this *RpcMsgHandle) DoMsgWitchTimeout(request *RpcRequest) {
	ch := make(chan struct{})
	go func(ch chan struct{}) {
		this.DoMsg(request)
		close(ch)
	}(ch)

	select {
	case <-time.After(time.Second * 30):
		logger.Error("DoMsgWitchTimeout target %s, args %v", request.Rpcdata.Target, request.Rpcdata.Args)
	case <-ch:
		return
	}
}

func (this *RpcMsgHandle) DeliverToMsgQueue(pkg interface{}) {
	request := pkg.(*RpcRequest)
	//add to worker pool
	index := rand.Int31n(utils.GlobalObject.PoolSize)
	taskQueue := this.TaskQueue[index]
	// logger.Debug(fmt.Sprintf("add to rpc pool : %d", index))
	taskQueue <- request
}

func (this *RpcMsgHandle) DoMsgFromGoRoutine(pkg interface{}) {
	request := pkg.(*RpcRequest)
	go this.DoMsg(request)
}

func (this *RpcMsgHandle) AddRouter(router interface{}) {
	value := reflect.ValueOf(router)
	tp := value.Type()
	for i := 0; i < value.NumMethod(); i += 1 {
		name := tp.Method(i).Name

		if _, ok := this.Apis[name]; ok {
			//存在
			panic("repeated rpc " + name)
		}
		this.Apis[name] = value.Method(i)
		logger.Info("add rpc " + name)
	}
}

func (this *RpcMsgHandle) StartWorkerLoop(poolSize int) {
	if utils.GlobalObject.IsThreadSafeMode() {
		this.TaskQueue[0] = make(chan *RpcRequest, utils.GlobalObject.MaxWorkerLen)
		go func() {
			for {
				select {
				case rpcRequest := <-this.TaskQueue[0]:
					this.DoMsg(rpcRequest)
				case delayCall := <-utils.GlobalObject.GetSafeTimer().GetTriggerChannel():
					delayCall.Call()
				}
			}
		}()
	} else {
		for i := 0; i < poolSize; i += 1 {
			c := make(chan *RpcRequest, utils.GlobalObject.MaxWorkerLen)
			this.TaskQueue[i] = c
			go func(index int, taskQueue chan *RpcRequest) {
				// logger.Info(fmt.Sprintf("init rpc thread pool %d.", index))
				for {
					request := <-taskQueue
					// this.DoMsg(request)
					this.DoMsgWitchTimeout(request)
				}

			}(i, c)
		}
	}
}
