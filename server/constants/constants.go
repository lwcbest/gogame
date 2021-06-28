package constants

type MsgCode struct {
	SUCCESS int32
	FAIL    int32
}

var MSG_CODE = MsgCode{
	SUCCESS: 1,
	FAIL:    0,
}

type ErrorMsg struct {
	SUCCESS     string
	MARSHALFAIL string
	NOTAVAIABLE string
	NOROOM      string
}

var ERROR_MSG = ErrorMsg{
	SUCCESS:     "SUCCESS",
	MARSHALFAIL: "unable to unmarshal data",
	NOTAVAIABLE: "no net server avaiable",
	NOROOM:      "游戏尚未开始，请稍后再试！",
}

//前端server保存的
const SERVER_PREFER = "SERVER_PREFER:"

//http方法
type Methods struct {
	GET  string
	POST string
}

var METHOD = Methods{
	GET:  "GET",
	POST: "POST",
}

const (
	ROOM_TYPE_STATE = 1
	ROOM_TYPE_FRAME = 2
	ROOM_TYPE_TIME  = 3
)

//数据库表名
const (
	CollectionGameAward   = "game_award"
	CollectionUserCoupon  = "user_coupon"
	CollectionDailyCoupon = "daily_coupon"
	CollectionGameHistory = "game_history"
)
