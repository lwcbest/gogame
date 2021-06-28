package filter

import (
	"errors"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"

	"github.com/golang/protobuf/proto"
)

type AuthFilter struct {
}

func (this *AuthFilter) BeforeDo(req *iface.CommonRequest) error {
	//忽略相关接口
	switch req.RouterStr {
	case "Login", "Login2":
		return nil
	}
	res := req.Session.Get("player")
	if res == nil {
		return errors.New("[Auth Error] need login first!")
	}
	logger.Info("[AuthFilter][Req] %+v", req.RealData)
	return nil
}

func (this *AuthFilter) AfterDo(req *iface.CommonRequest, resp proto.Message) error {
	logger.Info("[AuthFilter][Res] req:%+v res:%+v", req.RealData, resp)
	return nil
}
