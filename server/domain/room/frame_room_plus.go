package room

import (
	"errors"
	"gameserver-997/pb/gopb"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/domain/clock"
	"gameserver-997/server/domain/entity"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
)

const FRAME_TICKER = 10
const FRAME_INTERVAL = 20
const KEY_FRAME = 5

type TFrame = gopb.OnFrame
type Tcommand = gopb.OnFrame_FCommand

type FrameRoomV struct {
	*Room
	//真同步独有的数据
	randSeed     int64 //随机种子
	curFrame     int32 //当前帧
	totalFrame   int32 //总帧数
	frameTicker  *time.Ticker
	commandQueA  []*Tcommand //命令队列A
	commandQueB  []*Tcommand //命令队列B
	frameingQue  int8        //当前使用中的队列
	allCommand   []*Tcommand //所有走过的命令集合 断线重连用
	initCommands []*Tcommand //初始化帧
	allFrames    []*TFrame
}

func createFrameRoomV(roomType int, roomMode string, pool *RoomPool) (*FrameRoomV, error) {
	// constants := conf.GetSDConstantElement(playMod)
	base := createBaseRoom(roomType, roomMode, pool)
	room := &FrameRoomV{
		Room:     base,
		randSeed: time.Now().Unix(),
	}

	room.totalFrame = int32(room.roomConf.EndTime / FRAME_INTERVAL)
	return room, nil
}

func (this *FrameRoomV) checkStart() {
	if this.roomState == ROOM_STATE_READY {
		this.Start()
	}
}

func (this *FrameRoomV) PlayerEnter(player *entity.Player) error {
	err := this.Room.PlayerEnter(player)
	if err != nil {
		return err
	}

	//TODO push msg
	return nil
}

func (this *FrameRoomV) sendAllReady() {
	length := this.GetPlayerCount()
	logger.Fatal("need implement send all ready", length)
	// postions := this.uids
	// msg := &gopb.ReqCommonArgs{Args: make([]*gopb.ReqCommonArg, length)}
	// for i := length - 1; i >= 0; i-- {
	// 	index := rand.Intn(i + 1)
	// 	if index != i {
	// 		postions[i], postions[index] = postions[index], postions[i]
	// 	}
	// 	msg.Args[i] = &gopb.ReqCommonArg{StrVal1: postions[i]}
	// }
	// data, _ := proto.Marshal(msg)
	// this.pushDataByType(PUSH_MG_UIDS, "OnAllPlayerReady", data)
}

//玩家确认开始游戏 ready 之后
func (this *FrameRoomV) PlayerOkToStart(uid string) error {
	logger.Info("isAllReady when playerOkTostart %t\n", this.isAllPlayerReady())
	if this.roomState == ROOM_STATE_RUNNING {
		return errors.New("room already start")
	}
	if !this.isAllPlayerReady() {
		return errors.New("players not all ready")
	}
	p, ok := this.players.Load(uid)
	var player entity.Player
	if ok {
		player = p.(entity.Player)
	}
	player.State = entity.PLAYER_STATE_START
	logger.Info("this.isallplayerstart: ", this.isAllPlayerStart())

	if this.isAllPlayerStart() {
		this.Start()
	}
	return nil
}

func (this *FrameRoomV) PlayerOffline(uid, feServerId string) {
	p, have := this.GetPlayer(uid)
	if have {
		p.IsOffLine = true
	} else {
		return
	}

	if this.roomState == ROOM_STATE_READY {
		if this.isAllPlayerStart() {
			this.Start()
		}
	}
	this.channel.Leave(uid, feServerId)
}

//开始执行帧 第一个玩家进入或者房间开始执行
func (this *FrameRoomV) Start() {
	logger.Info("frameRoom start: ", this.roomState)
	if this.roomState == ROOM_STATE_RUNNING {
		return
	}

	if ok := atomic.CompareAndSwapInt32(&this.roomState, ROOM_STATE_READY, ROOM_STATE_RUNNING); ok == false {
		logger.Info("room already start...")
		return
	}

	this.startTime = time.Now()
	this.frameTicker = time.NewTicker(time.Millisecond * FRAME_TICKER)
	go func() {
		defer this.frameTicker.Stop()
		for {
			select {
			case <-this.frameTicker.C:
				stop := this.doFrame()
				if stop {
					timer := &clock.SimonTimer{}
					timer.Init(time.Second*3, this.endGame)
					logger.Info("stop frame from time")
					return
				}
				//TODO end ticker from base room
				// case <-this.endTicker:
				// 	logger.Info("stop frame from ticker")
				// 	return
			}
		}
	}()
}

//走帧判断
func (this *FrameRoomV) doFrame() bool {
	now := time.Now()
	for t := this.startTime.Add(time.Duration(this.curFrame) * FRAME_INTERVAL * time.Millisecond); t.Sub(now) < 0; {
		this.curFrame++
		if this.curFrame%KEY_FRAME == 0 {
			return this.doKeyFrame()
		}
	}
	return false
}

//关键帧
func (this *FrameRoomV) doKeyFrame() bool {
	msg := &gopb.OnFrame{CurFrame: this.curFrame}
	//可能存在读了frameingQue,然后先执行下面的代码，commandQueA被slice成0然后写入0的位置
	//用一个队列的问题：slice之后append会增加底层数组长度，一直触发重新分配，但是数据是安全的
	this.frameingQue ^= 1
	if this.frameingQue == 1 {
		msg.Commands = this.commandQueA
		this.commandQueA = this.commandQueA[0:0]
	} else {
		msg.Commands = this.commandQueB
		this.commandQueB = this.commandQueB[0:0]
	}
	//期望帧已经过了，补到后面的帧
	for _, command := range msg.Commands {
		if command.ExFrame == 0 {
			if command.EpFrame+KEY_FRAME < this.curFrame {
				command.ExFrame = this.curFrame + (command.EpFrame % KEY_FRAME) - KEY_FRAME
			} else {
				command.ExFrame = command.EpFrame + KEY_FRAME
			}
		}
	}

	this.allCommand = append(this.allCommand, msg.Commands...)
	this.allFrames = append(this.allFrames, msg)
	data, _ := proto.Marshal(msg)
	this.pushDataByType(PUSH_MG_UIDS, "OnFrame", data)
	return this.totalFrame > 0 && this.curFrame >= this.totalFrame
}

//新增命令
func (this *FrameRoomV) SyncFCommand(epFrame, exFrame int32, uid string, ctype int32, paramList []int32) {
	var command *Tcommand = &Tcommand{
		EpFrame:   epFrame,
		ExFrame:   exFrame,
		Uid:       uid,
		Ctype:     ctype,
		ParamList: paramList,
	}
	if this.frameingQue == 0 {
		this.commandQueA = append(this.commandQueA, command)
	} else {
		this.commandQueB = append(this.commandQueB, command)
	}
}

func (this *FrameRoomV) SyncResult(syncUid string, data []byte) error {
	msg := &gopb.ReqCommonArgs{}
	if err := proto.Unmarshal(data, msg); err != nil {
		return err
	}

	p, have := this.GetPlayer(syncUid)
	if have {
		p.State = entity.PLAYER_STATE_OVER
	} else {
		//TODO
	}

	for _, msg := range msg.Args {
		uid := msg.StrVal1
		score := msg.IntVal1
		p, have := this.GetPlayer(uid)
		if have {
			if p.Score == 0 || p.Score == score {
				p.Score = score
				//p.Hit = hit
				//p.Perfect = perfect
			}
		}
	}
	// for _, uid := range this.uids {
	// 	p := this.getPlayer(uid)
	// 	if p.State != entity.PLAYER_STATE_OVER {
	// 		return nil
	// 	}
	// }
	this.endGame()

	return nil
}

func (this *FrameRoomV) endGame() {
	ok := atomic.CompareAndSwapInt32(&this.roomState, ROOM_STATE_RUNNING, ROOM_STATE_OVER)
	if ok == false {
		//logger.Info("room already end by other gorountine ", this.mgUid)
		return
	}
	//logger.Info("endGame success called %s %s", this.mgUid, this.roomKey)
	this.syncGameResult()
	this.saveGameLog()
	this.Room.endGame()
}

func (this *FrameRoomV) saveGameLog() {
	newFrames := make([]*TFrame, 1, len(this.allFrames)+1)
	newFrames[0] = &gopb.OnFrame{
		CurFrame: 0,
		Commands: this.initCommands,
	}
	newFrames = append(newFrames, this.allFrames...)
	// players := make([]*mongo.Player, len(this.uids))
	// for i, uid := range this.uids {
	// 	p := this.getPlayer(uid)
	// 	players[i] = &mongo.Player{
	// 		Uid:    p.Uid,
	// 		Name:   p.Name,
	// 		AvaUrl: p.AvaUrl,
	// 		Select: p.Select,
	// 	}
	// }
	// id, err := mongo.SaveGameLog(this.roomKey, players, newFrames)
	logger.Info("save log when game end: %s %+v", "", "err")
}

//同步房间结算信息
// func (this *FrameRoomV) syncGameResult() {
// 	msg := &gopb.OnGameEnd{}
// 	result := make([]*gopb.OnGameEnd_GameResult, len(this.uids))
// 	syncResult := &ss.GameResult{GameType: "alpaca", RoomId: this.mgUid, RoomKey: this.roomKey.Hex(), ServerId: utils.GlobalObject.Name}
// 	var maxScore int32 = 0
// 	var winUid string
// 	var winAward *gopb.OnGameEnd_RoomAward
// 	dbResult := &mongo.GameAward{CreateTime: time.Now(), Id: this.roomKey, DeviceId: this.mgUid, Done: false}
// 	for index, uid := range this.uids {
// 		p := this.players[uid]
// 		star := calcGameStar("1", int(p.Score))
// 		result[index] = &gopb.OnGameEnd_GameResult{
// 			Uid:     uid,
// 			Score:   p.Score,
// 			Hit:     p.Hit,
// 			Perfect: p.Perfect,
// 			Star:    int32(star),
// 		}
// 		tag := fmt.Sprintf("star%d", star)
// 		ssAward := RandomRoomAward(this.allAwards, tag)
// 		var award *gopb.OnGameEnd_RoomAward
// 		if ssAward != nil {
// 			award := &gopb.OnGameEnd_RoomAward{
// 				Id:         ssAward.CouponId,
// 				CouponName: ssAward.CouponName,
// 				CouponImg:  ssAward.ImgApp,
// 				CouponDesc: ssAward.Description,
// 				Aid: ssAward.Aid,
// 			}
// 			result[index].Award = award
// 		}
// 		if p.Score > maxScore {
// 			maxScore = p.Score
// 			winUid = uid
// 			winAward = award
// 		}

// 		pAward := &mongo.PlayerAward{
// 			Uid:   uid,
// 			Score: p.Score,
// 			Star:  int32(star),
// 		}
// 		if award != nil {
// 			pAward.Award = award.Id
// 			pAward.AwardName = award.CouponName
// 			pAward.Aid = award.Aid
// 		}
// 		dbResult.PAwards = append(dbResult.PAwards, pAward)
// 	}
// 	msg.Result = result
// 	logger.Info("sync game result: %+v %+v\n", result, *dbResult)
// 	if winAward != nil {
// 		watchers := make([]*entity.Watcher, 0)
// 		for _, w := range this.watchers {
// 			if w.TarUid == winUid {
// 				watchers = append(watchers, w)
// 			}
// 		}
// 		if len(watchers) > 0 {
// 			randomIdx := rand.Intn(len(watchers))
// 			luckyWatcher := watchers[randomIdx]
// 			msg.LuckyStar = &gopb.OnGameEnd_LuckyPlayer{
// 				Uid:    luckyWatcher.Uid,
// 				Name:   luckyWatcher.Name,
// 				AvaUrl: luckyWatcher.AvaUrl,
// 				Award: &gopb.OnGameEnd_RoomAward{
// 					Id:         winAward.Id,
// 					CouponName: winAward.CouponName,
// 					CouponImg:  winAward.CouponImg,
// 					CouponDesc: winAward.CouponDesc,
// 				},
// 			}
// 			dbResult.LuckyStar = &mongo.LuckyAward{
// 				Uid:    luckyWatcher.Uid,
// 				ItemId: winAward.Id,
// 				ImgUrl: winAward.CouponImg,
// 				Name:   winAward.CouponName,
// 				Aid: winAward.Aid,
// 			}
// 		}
// 	}
// 	for _, w := range this.watchers {
// 		if w.TarUid == winUid {
// 			coin := 1
// 			if w.VoteN >= 5 {
// 				coin = 0
// 			}
// 			dbResult.WAwards = append(dbResult.WAwards, &mongo.WatcherAward{
// 				Uid:     w.Uid,
// 				CoinNum: coin, //todo random
// 			})
// 		}
// 	}
// 	err := mongo.SaveGameAward(dbResult)
// 	logger.Info("save gameaward %+v %+v", dbResult, err)
// 	data, _ := proto.Marshal(msg)
// 	this.pushDataByType(PUSH_MG_UIDS, "OnGameEnd", data)
// 	//同步小程序服务器游戏结果
// 	ss.SyncGameResult(syncResult)
// }

//工具方法

func convertCreateCommand(onCommands []*gopb.OnFrame_FCommand) []*gopb.ResNetCreateFRoom_FCommand {
	commands := make([]*gopb.ResNetCreateFRoom_FCommand, len(onCommands))
	for i, c := range onCommands {
		commands[i] = (*gopb.ResNetCreateFRoom_FCommand)(c)
	}
	return commands
}

func convertEnterCommand(onCommands []*gopb.OnFrame_FCommand) []*gopb.ResNetEnterFRoom_FCommand {
	commands := make([]*gopb.ResNetEnterFRoom_FCommand, len(onCommands))
	for i, c := range onCommands {
		commands[i] = (*gopb.ResNetEnterFRoom_FCommand)(c)
	}
	return commands
}
