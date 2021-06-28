package timer

import (
	"fmt"
	"gameserver-997/server/base/logger"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func test(a ...interface{}) {
	fmt.Println(a[0], "============", a[1])
}

var (
	tt = int64(0)
)

func Test(t *testing.T) {

	s := NewSafeTimerScheduel()
	go func() {
		for {
			df := <-s.GetTriggerChannel()
			df.Call()
			atomic.AddInt64(&tt, -1)
		}
	}()
	go func() {
		i := 0
		for i < 50000 {
			s.CreateTimer(int64(rand.Int31n(3600*1e3)), test, []interface{}{22, 33})
			atomic.AddInt64(&tt, 1)
			time.Sleep(1 * time.Second)
			i += 1
		}
	}()
	go func() {
		ii := 0
		for ii < 50000 {
			s.CreateTimer(int64(rand.Int31n(3600*1e3)), test, []interface{}{22, 33})
			atomic.AddInt64(&tt, 1)
			time.Sleep(1 * time.Second)
			ii += 1
		}
	}()

	for {
		time.Sleep(60 * time.Second)
		logger.Info("last timer: ", atomic.LoadInt64(&tt))
	}
}
