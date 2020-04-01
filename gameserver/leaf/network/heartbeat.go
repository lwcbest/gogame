package network

import (
	"time"

	"github.com/lwcbest/gogame/gameserver/leaf/log"
	"github.com/lwcbest/gogame/gameserver/leaf/timer"
)

const HeartbeatDuration int = 3000

type Heartbeat struct {
	timeout   int
	heartbeat int
	session   *Session
	ch        chan bool
	closed    bool
	st        *timer.SimonTimer
}

func (h *Heartbeat) Close() {
	h.closed = true
	h.st.Close()
}

func (h *Heartbeat) Handle() {
	if h.st == nil {
		h.st = &timer.SimonTimer{}
		h.st.Init(time.Duration(h.timeout)*time.Millisecond, func() {
			h.session.Close()
			h.session.Destroy()
			log.Release("Heartbeat timeout")
		})
	} else {
		h.st.Reset(time.Duration(h.timeout) * time.Millisecond)
	}

	if h.closed == false {
		SendHeartbeat(h.session)
	}
}

func SendHeartbeat(s *Session) {
	pkg := &Package{
		pkgType: PKG_HEARTBEAT,
		length:  0,
	}

	log.Debug("send heartbeat:", pkg)
	s.WritePkg(pkg)
}
