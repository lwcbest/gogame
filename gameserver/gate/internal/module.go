package internal

import (
	"gameserver/conf"
	"gameserver/game"
	"gameserver/leaf/network"
	"gameserver/msg"
)

type Module struct {
	*network.Gate
}

func (m *Module) OnInit() {
	m.Gate = &network.Gate{
		MaxConnNum:      conf.Server.MaxConnNum,
		PendingWriteNum: conf.PendingWriteNum,
		MaxMsgLen:       conf.MaxMsgLen,
		WSAddr:          conf.Server.WSAddr,
		HTTPTimeout:     conf.HTTPTimeout,
		CertFile:        conf.Server.CertFile,
		KeyFile:         conf.Server.KeyFile,
		TCPAddr:         conf.Server.TCPAddr,
		LenMsgLen:       conf.LenMsgLen,
		LittleEndian:    conf.LittleEndian,
		Processor:       msg.Instance,
		SessionChanRPC:  game.ChanRPC,
	}
}
