package iface

import (
	"github.com/golang/protobuf/proto"
)

/*
* 这里是一些公用的包，例如rpc传递需要类型定义
 */

//handle remote公用的唯一参数
type CommonRequest struct {
	Session   ISession
	BSession  *BackendSession
	Data      []byte
	RpcData   []interface{}
	Fconn     IWriter
	RealData  proto.Message
	RouterStr string
	ReqId     uint
	Router    string
}

type BackendSession struct {
	Uid        string
	FrontendId string
	Setting    map[string]interface{}
}
