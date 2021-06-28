package filter

import (
	"errors"
	"gameserver-997/pb/gopb"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"reflect"

	"github.com/golang/protobuf/proto"
)

type ProtoFilter struct {
}

func (this *ProtoFilter) BeforeDo(req *iface.CommonRequest) error {
	routerMap := map[string]reflect.Type{
		"Login":           reflect.TypeOf(gopb.ReqNetLogin{}),
		"EnterMatchQueue": reflect.TypeOf(gopb.ReqNetEnterMatchQueue{}),
	}

	if t, okk := routerMap[req.RouterStr]; okk {
		protoMsg := reflect.New(t).Interface().(proto.Message)
		req.RealData = protoMsg
		if err := proto.Unmarshal(req.Data, protoMsg); err != nil {
			logger.Error("[ProtoFilter][BeforeDo] %+v", err)
			return errors.New("unmarshal error")
		}
	} else {
		logger.Error("[proto_filter.go] not found api in [RouterMap]:  %v ", req.RouterStr)
		return errors.New("unmarshal error")
	}

	logger.Debug("[ProtoFilter] this is before do.....")
	return nil
}

func (this *ProtoFilter) AfterDo(req *iface.CommonRequest, resp proto.Message) error {
	logger.Debug("[ProtoFilter] this is after do.....")
	return nil
}
