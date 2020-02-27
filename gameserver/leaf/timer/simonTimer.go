package timer

import (
	"time"
)

type SimonTimer struct {
	timer     *time.Timer
	closeChan chan bool
	cb        func()
}

func (st *SimonTimer) Init(dur time.Duration, cb func()) {
	if st.timer == nil {
		st.timer = time.NewTimer(dur)
		st.closeChan = make(chan bool)
		st.cb = cb

		go func() {
			select {
			case <-st.timer.C:
				st.cb()
			case <-st.closeChan:
				return
			}
		}()
	}
}

func (st *SimonTimer) Reset(dur time.Duration) {
	st.timer.Reset(dur)
}

func (st *SimonTimer) Close() {
	st.timer.Stop() //关闭定时器
	st.closeChan<-true
}