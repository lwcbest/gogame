package match

import (
	"container/list"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/constants"
	"gameserver-997/server/domain/clock"
	"gameserver-997/server/domain/entity"
	"gameserver-997/server/util/rpc"
	"sync"
	"time"
)

const FIND_TIME time.Duration = time.Second * 2

type MatchPool struct {
	playerList list.List
	locker     sync.Mutex
	simonTimer clock.SimonTimer
	status     int32 // 0 for idle, 1 for run
	rule       int32
}

func (this *MatchPool) JoinPlayer(player entity.Player) int {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.playerList.PushBack(player)
	logger.Info("[JoinPlayer][len %v]push player:%v", this.playerList.Len(), player)
	return this.playerList.Len()
}

func (this *MatchPool) LeavePlayer(uid string) {
	this.locker.Lock()
	defer this.locker.Unlock()

	var next *list.Element
	for e := this.playerList.Front(); e != nil; {
		if e.Value.(entity.Player).Uid == uid {
			next = e.Next()
			this.playerList.Remove(e)
			logger.Info("[LeavePlayer][len %v]remove player:%v", this.playerList.Len(), e)
			e = next
		} else {
			e = e.Next()
		}
	}
}

func (this *MatchPool) ReadPlayers() list.List {
	return this.playerList
}

func (this *MatchPool) Start(rule int32) {
	logger.Fatal("[Start],abcd")
	if this.status == 1 {
		return
	}
	this.rule = rule

	//end test rule

	this.simonTimer = clock.SimonTimer{}
	this.simonTimer.Init(FIND_TIME, this.doMatch)
}

func (this *MatchPool) doMatch() {
	this.locker.Lock()
	defer this.locker.Unlock()

	count := this.playerList.Len()
	logger.Debug("[domatch] count: %v", count)

	if count >= 2 {
		//throw player
		this.compute(2)
	} else {
		logger.Debug("[doMatch]count:%v", count)
	}

	this.simonTimer.Reset(FIND_TIME)
}

func (this *MatchPool) compute(playerCount int) {
	//匹配人

	//减少人
	ps := make([]entity.Player, playerCount)
	p1 := this.playerList.Front()
	ps[0] = p1.Value.(entity.Player)
	var e *list.Element
	for i := 1; i < playerCount; i++ {
		e = p1.Next()
		ps[i] = e.Value.(entity.Player)
		this.playerList.Remove(e)
	}
	this.playerList.Remove(p1)

	//进入房间
	roomType := constants.ROOM_TYPE_FRAME
	roomMode := "grade1"

	rpcRes, err := rpc.RPCRandomCall(nil, "game", "CreateRoom", roomType, roomMode)
	if err != nil {
		logger.Fatal("err:%v", err)
	}
	gameServerName := rpcRes["serverName"]
	roomId := rpcRes["roomId"]

	sessionArray := make([]iface.ISession,0)
	for i := 0; i < playerCount; i++ {
		logger.Info("[Call][PlayerEnterGame]%v", ps[i])
		rpc.RPCCall(nil, ps[i].FeServerId, "PlayerEnterGame", roomId, gameServerName, ps[i].SessionId)
		sessionArray = append(sessionArray,ps[i].GetFakeSessionInfo())
	}

	//push alls
	rpc.PushMsgByUids("OnEnterRoom","ok",sessionArray)
	//channelService := utils.GlobalObject.TcpServer.GetChannelService()
	//channelService.PushMsgByUids()
}
