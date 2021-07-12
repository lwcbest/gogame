package entity

import (
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/service"
	"strings"
	"time"
)

const (
	PLAYER_STATE_IDLE  = 1
	PLAYER_STATE_READY = 2
	PLAYER_STATE_START = 3
	PLAYER_STATE_OVER  = 4

	PLAYER_STATE_WAIT    = 5
	PLAYER_STATE_OFFLINE = 6
)

type Player struct {
	Username string
	Pwd      string
	Uid      string
	Name     string
	AvaUrl   string
	Score    int32
	Level    int32
	Gender   int32

	IsOffLine  bool
	FeServerId string
	BeServerId string
	SessionId  string
	State      int
	Have       []string
	Select     []string
	EnterTime  string
}

func NewPlayer(uid string, name, avaUrl string, gender int32) *Player {
	return &Player{
		Uid:       uid,
		Name:      name,
		AvaUrl:    avaUrl,
		State:     PLAYER_STATE_IDLE,
		Gender:    gender,
		EnterTime: time.Now().String(),
	}
}

func GenPlayer(username, pwd string) *Player {
	return &Player{
		Username: username,
		Uid:      username,
		Pwd:      pwd,
	}
}

func (this *Player) IsReady() bool {
	return this.State == PLAYER_STATE_START || this.State == PLAYER_STATE_READY
}

func (this *Player) GetFakeSessionInfo() iface.ISession {
	return &service.Session{FrontendId: this.FeServerId,Uid: this.Uid}
}

//玩家非计数道具

const (
	PROP_ALPACA_SKIN = "2_0_1_" //羊驼皮肤
	PROP_ALPACA_HAT  = "2_0_2_" //羊驼帽子
)

func (this *Player) SetAlpacaInfo(have []string, sel []string) {
	for _, h := range have {
		if strings.HasPrefix(h, "2_0_") {
			this.Have = append(this.Have, h)
		}
	}
	for _, h := range sel {
		if strings.HasPrefix(h, "2_0_") {
			this.Select = append(this.Select, h)
		}
	}
}

func (this *Player) UpdateSelect(sel []string) {
	for _, s := range sel {
		ll := len(s)
		for i, ss := range this.Select {
			if ss[0:ll-2] == s[0:ll-2] {
				this.Select[i] = s
			}
		}
	}
}
