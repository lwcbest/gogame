package filter

import (
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"

	"github.com/golang/protobuf/proto"
)

type LogFilter struct {
}

func (this *LogFilter) BeforeDo(req *iface.CommonRequest) error {
	logger.Info("[LogFilter][Req] %+v", req.RealData)
	return nil
}

func (this *LogFilter) AfterDo(req *iface.CommonRequest, resp proto.Message) error {
	logger.Info("[LogFilter][Res] req:%+v res:%+v", req.RealData, resp)
	return nil
}
