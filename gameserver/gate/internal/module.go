package internal

import (
	"myGo/gameserver/conf"
	"myGo/gameserver/game"
	"myGo/gameserver/leaf/network"
	"myGo/gameserver/msg"
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
		Processor:       msg.Processor,
		SessionChanRPC:  game.ChanRPC,
	}
}
