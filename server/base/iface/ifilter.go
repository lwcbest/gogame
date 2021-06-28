package iface

import "github.com/golang/protobuf/proto"

type IFilter interface {
	BeforeDo(*CommonRequest) error
	AfterDo(*CommonRequest, proto.Message) error
}
