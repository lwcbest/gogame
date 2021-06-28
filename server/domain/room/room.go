package room

import (
	"errors"
	"gameserver-997/pb/gopb"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/base/utils"
	"gameserver-997/server/domain/entity"
	conf "gameserver-997/server/sd_conf"
	"gameserver-997/server/util"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	ROOM_STATE_IDLE    = 1 //等待玩家进入
	ROOM_STATE_READY   = 2 //等待玩家准备
	ROOM_STATE_RUNNING = 3 //运行中
	ROOM_STATE_OVER    = 4 //游戏结束
	ROOM_STATE_DESTROY = 5 //已经销毁
)

type Room struct {
	roomId   string
	roomType int
	roomMode string
	roomConf *RoomConf

	players  sync.Map
	channel  iface.IChannel
	roomPool *RoomPool

	maxWait        time.Duration
	maxRunning     time.Duration
	createTime     time.Time //创建房间的时间
	firstEnterTime time.Time //第一个玩家进入的时间
	forceReadyTime time.Time //强制开始的时间
	startTime      time.Time //开始游戏的时间

	roomState int32
}

type RoomConf struct {
	EndTime     int64
	PrepareTime int64
	WaitingTime int64
}

func createBaseRoom(roomType int, roomMode string, pool *RoomPool) *Room {
	roomname := pool.PoolId + "_" + strconv.FormatInt(time.Now().Unix(), 10) + "_" + strconv.Itoa(util.RandomSome(10000))
	channelService := utils.GlobalObject.TcpServer.GetChannelService()
	channel := channelService.GetChannel(roomname)
	room := &Room{
		roomType:   roomType,
		roomPool:   pool,
		roomMode:   roomMode,
		roomId:     roomname,
		createTime: time.Now(),
		channel:    channel,
	}

	room.roomConf = readConfig(roomMode)
	return room
}

func readConfig(roomMode string) *RoomConf {
	cons := conf.GetSDConstantElement(roomMode)
	logger.Fatal("need implement cons:", cons)
	roomConf := &RoomConf{
		// EndTime:     cons.EndTime.sGetInt(),
		// PrepareTime: cons.PrepareTime.GetInt(),
		// WaitingTime: cons.WaitingTime.GetInt(),
	}
	return roomConf
}

func (this *Room) GetRoomId() string {
	return this.roomId
}

func (this *Room) GetPlayerCount() int {
	count := 0
	this.players.Range(func(k interface{}, v interface{}) bool {
		count++
		return true
	})

	return count
}

func (this *Room) GetPlayer(uid string) (entity.Player, bool) {
	p, ok := this.players.Load(uid)
	var player entity.Player
	if ok {
		player = p.(entity.Player)
		return player, true
	}

	return player, false
}

func (this *Room) GetState() int32 {
	return this.roomState
}

func (this *Room) onPlayerEnter(uid string) {
	if this.firstEnterTime.IsZero() {
		this.firstEnterTime = time.Now()
	}
}

func (this *Room) PlayerEnter(player *entity.Player) error {
	logger.Info("[PlayerEnter]player:%v", player)
	this.players.Store(player.Uid, player)
	return nil
}

//玩家准备
func (this *Room) PlayerReady(uid string) error {
	res, have := this.players.Load(uid)
	if !have {
		return errors.New("no target player set ready~")
	}
	player := res.(*entity.Player)
	player.State = entity.PLAYER_STATE_READY

	//push msg
	return nil
}

func (this *Room) isAllPlayerReady() bool {
	// readyNum := 0
	// for _, uid := range this.uids {
	// 	p := this.getPlayer(uid)
	// 	//判断都ready的时候不在线的就不算了防止玩家没进来
	// 	if p.State != entity.PLAYER_STATE_READY && p.State != entity.PLAYER_STATE_START {

	// 	} else {
	// 		readyNum = readyNum + 1
	// 	}
	// }
	// return readyNum >= this.readyPlayer
	return false
}

func (this *Room) isAllPlayerStart() bool {
	//okNum := 0
	// for _, uid := range this.uids {
	// 	p := this.getPlayer(uid)
	// 	if p.State == entity.PLAYER_STATE_START || p.IsOffLine == true {
	// 		okNum = okNum + 1
	// 	} else {
	// 		return false
	// 	}
	// }
	return true
}

//游戏主动结束 todo确认是否应该在派生类里做
func (this *Room) endGame() {
	// for _, uid := range this.uids {
	// 	p := this.getPlayer(uid)
	// 	if p != nil {
	// 		this.channel.Leave(uid, p.FeServerId)
	// 		this.players.Delete(uid)
	// 	}
	// 	this.roomPool.RemovePlayer(uid)
	// 	go redis.DelPlayerGaming(uid)
	// }
	// //p := this.getPlayer(this.mgUid)
	// //p.SyncState = nil
	// this.uids = []string{}
	// this.destroy()
}

//销毁
func (this *Room) destroy() {

}

func (this *Room) Destroy() {
	this.destroy()
}

//-------------------------游戏服务端透传方法---------------------------

//同步房间结算信息
func (this *Room) syncGameResult() {
	// var maxScore int32 = 0
	// var winUid string
	// // var winAward string
	// var maxStar int
	// dbResult := &mongo.GameAward{CreateTime: time.Now(), Id: this.roomKey, DeviceId: this.mgUid, Aid: this.aid, GameKey: this.playMod, GameName: this.gameName, Done: false}
	// for _, uid := range this.uids {
	// 	p := this.getPlayer(uid)
	// 	star := this.calcGameStar(int(p.Score))
	// 	// ssAward := randomRoomAward(this.allAwards, star)

	// 	if p.Score > maxScore {
	// 		maxScore = p.Score
	// 		winUid = uid
	// 		maxStar = star
	// 	}

	// 	pAward := &mongo.PlayerAward{
	// 		Uid:       uid,
	// 		Score:     p.Score,
	// 		Star:      int32(star),
	// 		NickName:  p.Name,
	// 		AvatarUrl: p.AvaUrl,
	// 		Gender:    p.Gender,
	// 	}
	// 	// if ssAward != "" {
	// 	// 	pAward.Award = ssAward
	// 	// 	pAward.Aid = this.aid
	// 	// }
	// 	dbResult.PAwards = append(dbResult.PAwards, pAward)
	// }
	// logger.Info("sync game result: %+v\n", *dbResult)

	// awardUpvotes := make(map[string]bool)
	// var watchers []*entity.Watcher //所有助力了获胜者的人
	// if maxStar > 0 {
	// 	watchers = make([]*entity.Watcher, 0)
	// 	this.watchers.Range(func(key, val interface{}) bool {
	// 		w := val.(*entity.Watcher)
	// 		if w.TarUid == winUid {
	// 			watchers = append(watchers, w)
	// 		}
	// 		return true
	// 	})

	// 	if len(watchers) > 0 {
	// 		randomIdx := rand.Intn(len(watchers))
	// 		// luckyWatcher := watchers[randomIdx]
	// 		// dbResult.LuckyStar = &mongo.LuckyAward{
	// 		// 	Uid:    luckyWatcher.Uid,
	// 		// 	ItemId: randomLuckyAward(this.allAwards, maxStar),
	// 		// 	Aid:    this.aid,
	// 		// }
	// 		i := 0
	// 		for {
	// 			randomIdx += i
	// 			i++
	// 			if randomIdx >= len(watchers) {
	// 				randomIdx = 0
	// 			}
	// 			awardUpvotes[watchers[randomIdx].Uid] = true
	// 			if i >= this.conf.UpvoteNum {
	// 				break
	// 			}
	// 		}
	// 	}
	// }
	// this.watchers.Range(func(key, val interface{}) bool {
	// 	w := val.(*entity.Watcher)
	// 	if w.TarUid == winUid {
	// 		coin := 1
	// 		if w.VoteN >= 5 {
	// 			coin = 0
	// 		}
	// 		watcherAward := &mongo.WatcherAward{
	// 			Uid:       w.Uid,
	// 			CoinNum:   coin,
	// 			NickName:  w.Name,
	// 			AvatarUrl: w.AvaUrl,
	// 			Aid:       this.aid,
	// 		}
	// 		if awardUpvotes[w.Uid] {
	// 			// watcherAward.ItemId = randomLuckyAward(this.allAwards, maxStar)
	// 			//watcherAward.Star = randomLuckyStar(maxStar)
	// 		}
	// 		dbResult.WAwards = append(dbResult.WAwards, watcherAward)
	// 	}
	// 	return true
	// })
	// //存库
	// err := mongo.SaveGameAward(dbResult)
	// logger.Info("save gameaward %+v %+v", *dbResult, err)

	// data := "abc"

	// this.pushDataByType(PUSH_MG_UIDS, "OnGameEnd", []byte(data))

}

//------------------------------------push方法-----------------------------------
//玩家进入，推送给已经在的玩家
func (this *Room) pushPlayerEnter(p *entity.Player) {
	player := &gopb.OnPlayerEnterFRoom_PlayerInfo{
		Uid:     p.Uid,
		Name:    p.Name,
		AvaUrl:  p.AvaUrl,
		IsReady: p.IsReady(),
		Have:    p.Have,
		Select:  p.Select,
		Gender:  p.Gender,
	}

	msg := &gopb.OnPlayerEnterFRoom{
		LeftTIme: (int64(this.maxWait) - int64(time.Now().Sub(this.firstEnterTime))) / 1e6,
		Player:   player,
	}
	data, _ := proto.Marshal(msg)
	this.pushDataByType(PUSH_MG_UIDS, "OnPlayerEnter", data)
}

func (this *Room) pushPlayerClientEnter(uid string) {
	msg := &gopb.OnPlayerClientEnter{Uid: uid}
	data, _ := proto.Marshal(msg)
	time.AfterFunc(time.Millisecond*100, func() {
		this.pushDataByType(PUSH_MG_UIDS, "OnPlayerClientEnter", data)
	})
}

func (this *Room) pushPlayerWait(p *entity.Player) {
	pMsg := &gopb.PlayerInfo{
		Uid:    p.Uid,
		Name:   p.Name,
		AvaUrl: p.AvaUrl,
		Gender: p.Gender,
	}
	pData, _ := proto.Marshal(pMsg)
	this.pushDataByType(PUSH_MG, "OnPlayerWaiting", pData)
}

//玩家准备，推送给其他玩家
func (this *Room) pushPlayerReady(p *entity.Player) {
	msg := &gopb.OnPlayerReady{
		Uid:    p.Uid,
		Select: p.Select,
		Gender: p.Gender,
	}
	data, _ := proto.Marshal(msg)
	this.pushDataByType(PUSH_MG_UIDS, "OnPlayerReady", data)
}

//根据类型选择推送给谁
const (
	PUSH_ALL     = 0
	PUSH_MG      = 1
	PUSH_UIDS    = 2
	PUSH_MG_UIDS = 3
	PUSH_WAIT    = 4
	PUSH_MG_WAIT = 5
	PUSH_NONE    = 6
)

func (this *Room) pushDataByType(cType int32, route string, data []byte) {
	// var uids []string
	// switch cType {
	// case PUSH_ALL: //同步所有
	// 	uids = append([]string{this.mgUid}, this.uids...)
	// 	this.channel.PushMessage(route, data)
	// case PUSH_MG: //同步大屏
	// 	uids = []string{this.mgUid}
	// 	this.channel.PushMessageByUids(uids, route, data)
	// case PUSH_UIDS: //同步玩家
	// 	uids = append(uids, this.uids...)
	// 	this.channel.PushMessageByUids(this.uids, route, data)
	// case PUSH_MG_UIDS: //同步大屏和玩家
	// 	uids = append([]string{this.mgUid}, this.uids...)
	// 	this.channel.PushMessageByUids(append(this.uids, this.mgUid), route, data)
	// case PUSH_WAIT: //同步候补
	// 	uids = append(uids, this.waitings...)
	// 	this.channel.PushMessageByUids(this.waitings, route, data)
	// case PUSH_MG_WAIT: //同步大屏和候补
	// 	uids = append([]string{this.mgUid}, this.waitings...)
	// 	this.channel.PushMessageByUids(append(this.waitings, this.mgUid), route, data)
	// }
	// this.addState(3, route, uids, data)
}

//-------------------------------------------------------常量和基础方法----------------------

// func calcGameStar(id string, score int) int {
// 	conf := conf.GetSDScoreToStarElement(id)
// 	if conf == nil {
// 		return 0
// 	}
// 	if score < conf.Star1 {
// 		return 0
// 	} else if score < conf.Star2 {
// 		return 1
// 	} else if score < conf.Star3 {
// 		return 2
// 	}
// 	return 3
// }

func ContainTag(tag, tar string, awardMode string) bool {
	if awardMode != "" && awardMode == "1" { //频率券忽略，旧的没有这个值需要兼容
		return false
	}
	tagArr := strings.Split(tag, ",")
	for _, v := range tagArr {
		if v == tar {
			return true
		}
	}
	return false
}

func ParseTag2Star(tag string) int {
	tagArr := strings.Split(tag, ",")
	for _, v := range tagArr {
		if v == "star1" {
			return 1
		} else if v == "star2" {
			return 2
		} else if v == "star3" {
			return 3
		}
	}
	return 0
}
