package game_server

// "gameserver-997/server/base/logger"

type GameHan struct {
}

// func (this *GameHan) SyncState(request *iface.CommonRequest) map[string]interface{} {
// 	msg := &gopb.ReqGameSyncState{}
// 	data := (request.RpcData[0]).([]byte)
// 	if err := proto.Unmarshal(data, msg); err != nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: constants.ERROR_MSG.MARSHALFAIL})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	session := request.BSession
// 	r := room.GetRoomPool().GetPlayerRoom(session.Uid)
// 	if r == nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.NOROOM, Msg: constants.ERROR_MSG.NOROOM})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	if err := room.GetRoomPool().SyncState(session.Uid, msg.Tag, msg.SType, msg.Route, msg.Data); err != nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: err.Error()})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.SUCCESS})
// 	return map[string]interface{}{"data": ret}
// }

// func (this *GameHan) SyncResult(request *iface.CommonRequest) map[string]interface{} {
// 	data := (request.RpcData[0]).([]byte)
// 	session := request.BSession
// 	r := room.GetRoomPool().GetPlayerRoom(session.Uid)
// 	if r == nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.NOROOM, Msg: constants.ERROR_MSG.NOROOM})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	if err := room.GetRoomPool().SyncResult(session.Uid, data); err != nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: err.Error()})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.SUCCESS})
// 	return map[string]interface{}{"data": ret}
// }

// func (this *GameHan) SyncFCommand(request *iface.CommonRequest) map[string]interface{} {
// 	now1 := time.Now()
// 	msg := &gopb.ReqGameSyncFCommand{}
// 	data := (request.RpcData[0]).([]byte)
// 	if err := proto.Unmarshal(data, msg); err != nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: constants.ERROR_MSG.MARSHALFAIL})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	session := request.BSession
// 	r := room.GetRoomPool().GetPlayerRoom(session.Uid)
// 	if r == nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.NOROOM, Msg: constants.ERROR_MSG.NOROOM})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	if err := room.GetRoomPool().SyncFCommand(session.Uid, msg.EpFrame, msg.ExFrame, msg.Uid, msg.Ctype, msg.ParamList); err != nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: err.Error()})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.SUCCESS})
// 	logger.Info("SyncFCommand cost time: %d", time.Now().UnixNano()-now1.UnixNano())
// 	return map[string]interface{}{"data": ret}
// }

// func (this *GameHan) PlayerReady(request *iface.CommonRequest) map[string]interface{} {
// 	session := request.BSession
// 	msg := &gopb.ReqGameFPlayerReday{}
// 	data := (request.RpcData[0]).([]byte)
// 	if err := proto.Unmarshal(data, msg); err != nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: constants.ERROR_MSG.MARSHALFAIL})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	r := room.GetRoomPool().GetPlayerRoom(session.Uid)
// 	if r == nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.NOROOM, Msg: constants.ERROR_MSG.NOROOM})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	logger.Info("player ready: %v %v", data, msg)
// 	if err := room.GetRoomPool().PlayerReady(session.Uid, msg.Gender, msg.Select...); err != nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: err.Error()})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.SUCCESS})
// 	return map[string]interface{}{"data": ret}
// }

// func (this *GameHan) OkToStart(request *iface.CommonRequest) map[string]interface{} {
// 	session := request.BSession
// 	r := room.GetRoomPool().GetPlayerRoom(session.Uid)
// 	if r == nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.NOROOM, Msg: constants.ERROR_MSG.NOROOM})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	if err := room.GetRoomPool().PlayerOkToStart(session.Uid); err != nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: err.Error()})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.SUCCESS})
// 	return map[string]interface{}{"data": ret}
// }

// func (this *GameHan) RandomAward(request *iface.CommonRequest) map[string]interface{} {
// 	session := request.BSession
// 	msg := &gopb.ReqGamRandomAward{}
// 	data := (request.RpcData[0]).([]byte)
// 	if err := proto.Unmarshal(data, msg); err != nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: constants.ERROR_MSG.MARSHALFAIL})
// 		return map[string]interface{}{"data": ret}
// 	}

// 	data, err := room.GetRoomPool().RandomAward(session.Uid, msg.Start, msg.Steps)
// 	if err != nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: err.Error()})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	return map[string]interface{}{"data": data}
// }

// func (this *GameHan) QuitRoom(request *iface.CommonRequest) map[string]interface{} {
// 	resp := make(map[string]interface{})
// 	ret := &gopb.ResError{Code: constants.ERROR_CODE.SUCCESS}
// 	err := room.GetRoomPool().PlayerQuit(request.BSession.Uid)
// 	if err != nil {
// 		ret.Code = constants.ERROR_CODE.FAIL
// 		ret.Msg = err.Error()
// 	}
// 	data, _ := proto.Marshal(ret)
// 	resp["data"] = data
// 	return resp
// }

// func (this *GameHan) KickPlayer(request *iface.CommonRequest) map[string]interface{} {
// 	msg := &gopb.ReqCommonArg{}
// 	data := (request.RpcData[0]).([]byte)
// 	if err := proto.Unmarshal(data, msg); err != nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: constants.ERROR_MSG.MARSHALFAIL})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	err := room.GetRoomPool().PlayerQuit(msg.StrVal1)
// 	if err != nil {
// 		ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.FAIL, Msg: err.Error()})
// 		return map[string]interface{}{"data": ret}
// 	}
// 	ret, _ := proto.Marshal(&gopb.ResError{Code: constants.ERROR_CODE.SUCCESS})
// 	return map[string]interface{}{"data": ret}
// }
