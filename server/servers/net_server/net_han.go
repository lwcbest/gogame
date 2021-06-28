package net_server

import (
	"fmt"
	"gameserver-997/pb/gopb"
	"gameserver-997/server/base/clusterserver"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/constants"
	"gameserver-997/server/domain/entity"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

type NetHan struct {
}

//RPC
func (this *NetHan) Proxy(serverType string, target string, session iface.ISession, data []byte) (rdata []byte) {
	msg := &gopb.ResError{}
	now1 := time.Now()
	result, err := clusterserver.GlobalClusterServer.RpcSystemCallServerType(session, serverType, target, data)
	if err == nil {
		rdata = result["data"].([]byte)
	} else {
		msg.Code = constants.MSG_CODE.FAIL
		msg.Msg = fmt.Sprintf("rpc error: %s", err.Error())
		rdata, _ = proto.Marshal(msg)
		logger.Info("proxy %s err: %+v cost: %d", target, err, time.Now().Unix()-now1.Unix())
	}
	return
}

func (this *NetHan) Login(request *iface.CommonRequest) (response proto.Message) {
	realReq := request.RealData.(*gopb.ReqNetLogin)
	logger.Info("username:%+v,pwd:%+v", realReq.Username, realReq.Password)
	player := entity.GenPlayer(realReq.Username, realReq.Password)

	//TODO 没有判断密码，后面接入短信验证系统
	logger.Info("[Login]没有判断密码，后面接入短信验证系统")
	res, err := clusterserver.GlobalClusterServer.RpcRandomCallServerType(request.Session, "dc", "GetPlayerInfo", player)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			res, err = clusterserver.GlobalClusterServer.RpcRandomCallServerType(request.Session, "dc", "SavePlayerInfo", player)
		} else {
			logger.Error("[RPC][GetPlayerInfo]error:%v", err)
		}
	}

	//TODO 没有判断断线重连
	logger.Info("[Login]没有判断断线重连")
	player = res["player"].(*entity.Player)
	player.FeServerId = request.Session.GetServerId()
	player.SessionId = request.Session.GetId()
	request.Session.Bind(player.Uid)
	request.Session.Set("player", player)

	//TODO game config reader
	logger.Info("[Login]game config reader")
	response = &gopb.ResNetLogin{
		Code: constants.MSG_CODE.SUCCESS,
		Player: &gopb.ResNetLogin_Player{
			Username: player.Username,
			Uid:      player.Uid,
			Name:     player.Name,
			AvaUrl:   player.AvaUrl,
			Score:    player.Score,
			Level:    player.Level,
		},
		GameConfig: "{}",
	}

	return
}

func (this *NetHan) EnterMatchQueue(request *iface.CommonRequest) (response proto.Message) {
	realReq := request.RealData.(*gopb.ReqNetEnterMatchQueue)
	logger.Info("username:%+v,pwd:%+v", realReq.Level)
	//use level to enter match queue
	rpcRes, err := clusterserver.GlobalClusterServer.RpcCallServerId(request.Session, "match1", "JoinMatchQueue", realReq.Level)
	if err != nil {
		response = &gopb.ResNetEnterMatchQueue{Code: constants.MSG_CODE.FAIL}
		return
	}

	logger.Info("[EnterMatchQueue]queue length:%v", rpcRes["CurLen"])
	response = &gopb.ResNetLogin{Code: constants.MSG_CODE.SUCCESS}
	return
}

// func (this *NetApi) CreateRoom(request *iface.CommonRequest) proto.Message {
// 	code, msg, data := this.CreateCommonRoom(request, constants.ROOM_TYPE_FRAME)
// 	if code != constants.ERROR_CODE.SUCCESS {
// 		return &gopb.ResError{Code: code, Msg: msg}
// 	}
// 	reply := &gopb.ResNetCreateFRoom{}
// 	err := proto.Unmarshal(data, reply)
// 	if err != nil {
// 		return &gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: constants.ERROR_MSG.MARSHALFAIL}
// 	}
// 	return reply
// }

// //创建状态同步房间
// func (this *NetApi) CreateStateRoom(request *iface.CommonRequest) proto.Message {
// 	code, msg, data := this.CreateCommonRoom(request, constants.ROOM_TYPE_STATE)
// 	if code != constants.ERROR_CODE.SUCCESS {
// 		return &gopb.ResError{Code: code, Msg: msg}
// 	}
// 	reply := &gopb.ResNetCreateSRoom{}
// 	err := proto.Unmarshal(data, reply)
// 	if err != nil {
// 		return &gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: constants.ERROR_MSG.MARSHALFAIL}
// 	}
// 	return reply
// }

// func (this *NetApi) CreateTimeRoom(request *iface.CommonRequest) proto.Message {
// 	code, msg, data := this.CreateCommonRoom(request, constants.ROOM_TYPE_TIME)
// 	if code != constants.ERROR_CODE.SUCCESS {
// 		return &gopb.ResError{Code: code, Msg: msg}
// 	}
// 	reply := &gopb.ResNetCreateSRoom{}
// 	err := proto.Unmarshal(data, reply)
// 	if err != nil {
// 		return &gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: constants.ERROR_MSG.MARSHALFAIL}
// 	}
// 	return reply
// }

// func (this *NetApi) CreateCommonRoom(request *iface.CommonRequest, roomType int) (int32, string, []byte) {
// 	req := &gopb.ReqNetCreateRoom{}
// 	if err := proto.Unmarshal(request.Data, req); err != nil {
// 		return constants.ERROR_CODE.MARSHALFAIL, constants.ERROR_MSG.MARSHALFAIL, nil
// 	}

// 	gaming, err := redis.GetPlayerGaming(req.Mac)
// 	ip := request.Session.RemoteIp()
// 	curIp := req.FingerPrint
// 	logger.Info("CreateStateRoom req %+v\n", req, curIp, ip)
// 	if gaming != nil && gaming.Ip != "" && gaming.Ip != curIp && curIp != "" {
// 		return constants.ERROR_CODE.ALREADYIN, "房间已经被占用", nil
// 	}

// 	defer redis.DelRoomCreating(req.Mac)
// 	done, err := redis.SetRoomCreating(req.Mac, time.Minute)
// 	msg := "房间正在创建中"
// 	if done == false {
// 		return constants.ERROR_CODE.CREATING, msg, nil
// 	}
// 	// sessionService := utils.GlobalObject.TcpServer.GetSessionService()
// 	// sessionService.KickOld(req.Mac, request.Session.GetId(), "newLogin")
// 	request.Session.Bind(req.Mac)

// 	var result map[string]interface{}
// 	logger.Info("CreateCommonRoom BackendSession: %+v", request.Session.BackendSession())
// 	if err == nil && gaming.ServerId != "" {
// 		result, err = rpc.RpcCallServerId(request.Session, gaming.ServerId, "CreateRoom", roomType, req.PlayMod, curIp, int(req.ItemId))
// 	}
// 	if gaming.ServerId == "" || err != nil {
// 		result, err = rpc.RpcCallServerType(request.Session, "game", "CreateRoom", roomType, req.PlayMod, curIp, int(req.ItemId))
// 	}
// 	if err != nil {
// 		return constants.ERROR_CODE.FAIL, err.Error(), nil
// 	}

// 	key := constants.SERVER_PREFER + "game"
// 	request.Session.Set(key, result["serverId"])
// 	return constants.ERROR_CODE.SUCCESS, "", result["data"].([]byte)
// }

// func (this *NetApi) EnterRoom(request *iface.CommonRequest) proto.Message {
// 	req := &gopb.ReqNetEnterRoom{}
// 	if err := proto.Unmarshal(request.Data, req); err != nil {
// 		return &gopb.ResError{Code: constants.ERROR_CODE.MARSHALFAIL, Msg: constants.ERROR_MSG.MARSHALFAIL}
// 	}

// 	uid, err := service.AuthToken(req.Token, req.Token)
// 	if err != nil {
// 		uid = req.Token
// 	}
// 	// sessionService := utils.GlobalObject.TcpServer.GetSessionService()
// 	// sessionService.KickOld(uid, request.Session.GetId(), "newLogin")
// 	request.Session.Bind(uid)

// 	gaming, err := redis.GetPlayerGaming(req.DeviceId)
// 	if err != nil || gaming.ServerId == "" {
// 		return &gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: constants.ERROR_MSG.NOROOM}
// 	}

// 	result, err := rpc.RpcCallServerId(request.Session, gaming.ServerId, "EnterRoom", req.DeviceId)
// 	if err != nil {
// 		return &gopb.ResError{Code: constants.ERROR_CODE.MARSHALFAIL, Msg: err.Error()}
// 	}
// 	key := constants.SERVER_PREFER + "game"
// 	request.Session.Set(key, gaming.ServerId)
// 	reply := &gopb.ResNetEnterFRoom{}
// 	proto.Unmarshal(result["data"].([]byte), reply)
// 	return reply
// }

// func (this *NetApi) EnterStateRoom(request *iface.CommonRequest) proto.Message {
// 	code, msg, data := EnterCommonRoom(request)
// 	if code != constants.ERROR_CODE.SUCCESS {
// 		return &gopb.ResError{Code: code, Msg: msg}
// 	}
// 	reply := &gopb.ResNetEnterSRoom{}
// 	proto.Unmarshal(data, reply)
// 	return reply
// }

// func (this *NetApi) EnterTimeRoom(request *iface.CommonRequest) proto.Message {
// 	code, msg, data := EnterCommonRoom(request)
// 	if code != constants.ERROR_CODE.SUCCESS {
// 		return &gopb.ResError{Code: code, Msg: msg}
// 	}
// 	reply := &gopb.ResNetEnterTRoom{}
// 	proto.Unmarshal(data, reply)
// 	return reply
// }

// func EnterCommonRoom(request *iface.CommonRequest) (code int32, msg string, data []byte) {
// 	req := &gopb.ReqNetEnterRoom{}
// 	if err := proto.Unmarshal(request.Data, req); err != nil {
// 		return constants.ERROR_CODE.MARSHALFAIL, constants.ERROR_MSG.MARSHALFAIL, nil
// 	}

// 	//先处理一下玩家登陆数据
// 	uid, err := service.AuthToken(req.Token, req.Token)
// 	if err != nil {
// 		uid = req.Token
// 	}
// 	// sessionService := utils.GlobalObject.TcpServer.GetSessionService()
// 	// sessionService.Kick(uid, "newLogin")
// 	// sessionService.KickOld(uid, request.Session.GetId(), "newLogin")
// 	request.Session.Bind(uid)

// 	gaming, err := redis.GetPlayerGaming(req.DeviceId)
// 	if err != nil || gaming.ServerId == "" {
// 		return constants.ERROR_CODE.FAIL, constants.ERROR_MSG.NOROOM, nil
// 	}

// 	result, err := rpc.RpcCallServerId(request.Session, gaming.ServerId, "EnterRoom", req.DeviceId)
// 	if err != nil {
// 		return constants.ERROR_CODE.MARSHALFAIL, err.Error(), nil
// 	}
// 	key := constants.SERVER_PREFER + "game"
// 	request.Session.Set(key, gaming.ServerId)
// 	return constants.ERROR_CODE.SUCCESS, "", result["data"].([]byte)
// }
