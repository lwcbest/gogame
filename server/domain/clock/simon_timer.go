package clock

import (
	"gameserver-997/server/base/logger"
	"time"
)

type SimonTimer struct {
	timer      *time.Timer
	closeChan  chan bool
	cb         func()
	cbWhitArgs func(...interface{})
}

func (st *SimonTimer) Init(dur time.Duration, cb func()) {
	if st.timer == nil {
		st.timer = time.NewTimer(dur)
		st.closeChan = make(chan bool, 1)
		st.cb = cb

		go func() {
			for {
				select {
				case <-st.timer.C:
					st.cb()
				case <-st.closeChan:
					return
				}
			}
		}()
	} else {
		logger.Fatal("simon fatal reset")
		st.timer.Reset(dur)
	}
}

func (st *SimonTimer) Reset(dur time.Duration) {
	st.timer.Reset(dur)
}

func (st *SimonTimer) Close() {
	if st.timer != nil {
		logger.Fatal("start close")
		st.closeChan <- true
		logger.Fatal("start close2")
		st.timer.Stop() //关闭定时器
	}
}
