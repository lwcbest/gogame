package utils

import (
	"gameserver-997/server/base/iface"

	"github.com/golang/protobuf/proto"
)

type MsgFilters struct {
	filters       []iface.IFilter
	beforeActions []func(*iface.CommonRequest) error
	afterActions  []func(*iface.CommonRequest, proto.Message) error
}

func (this *MsgFilters) AddFilter(ft iface.IFilter) {
	this.filters = append(this.filters, ft)
	this.beforeActions = append(this.beforeActions, ft.BeforeDo)
	this.afterActions = append(this.afterActions, ft.AfterDo)
}

func (this *MsgFilters) DoBeforeActions(req *iface.CommonRequest) error {
	for i := len(this.beforeActions) - 1; i > -1; i-- {
		err := this.beforeActions[i](req)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *MsgFilters) DoAfterActions(req *iface.CommonRequest, resMsg proto.Message) error {
	for i := 0; i < len(this.afterActions); i++ {
		err := this.afterActions[i](req, resMsg)
		if err != nil {
			return err
		}
	}
	return nil
}
