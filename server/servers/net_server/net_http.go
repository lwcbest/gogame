package net_server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gameserver-997/server/util"
)

type NetHttp struct {
}

func (this *NetHttp) Hello(w http.ResponseWriter, r *http.Request) {
	queryFrom, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil && len(queryFrom["name"]) > 0 {
		w.Write([]byte(fmt.Sprintf("hello %s", queryFrom["name"][0])))
	} else {
		w.Write([]byte("hello unknow body"))
	}
}

//----------------------------------------------------以上都是测试方法，后面是真正有用的逻辑----------------------
type reqRoomStatus struct {
	RoomId   string `json:"roomId"`
	ServerId string `json:"serverId"`
	Uid      string `json:"uid"`
}

type ResRoomInfo struct {
	Status int `json:"status"` //1进房间2进候补3满人
}

type ReqEnterRoom struct {
	RoomId   string    `json:"roomId"`
	ServerId string    `json:"serverId"`
	UserInfo *UserInfo `json:"userInfo"`
}

type UserInfo struct {
	Uid       string   `json:"uid"`
	NickName  string   `json:"nickName"`
	AvatarUrl string   `json:"avatarUrl"`
	Gender    int      `json:"gender"`
	Have      []string `json:"have,omitempty"`
	Select    []string `json:"select,omitempty"`
}

type reqRoomVoteInfo struct {
	DeviceId string `json:"deviceId"`
	Uid      string `json:"uid"`
}

type ResRoomVoteInfo struct {
	RoomKey string   `json:"roomKey"`
	Code    int      `json:"code"` //游戏是否进行中 //0房间不存在 1房间未开始 2房间进行中
	Uids    []string `json:"uids,omitempty"`
	VoteUid string   `json:"voteUid"`
}

type reqEnterWatch struct {
	DeviceId string    `json:"deviceId"`
	TarUid   string    `json:"tarUid"`
	VoteN    int       `json:"voteN"`
	UserInfo *UserInfo `json:"userInfo"`
}

//点赞
type reqUpVote struct {
	DeviceId string `json:"deviceId"`
	Uid      string `json:"uid"`
	Val      int    `json:"val"`
}

type reqNotifyEnterQueue struct {
	DeviceId string `json:"deviceId"`
	UserId   string `json:"userId"`
	IsElec   bool   `json:"isElec"`
	IsTimes  bool   `json:"isTimes"`
}

type reqNotifyQuitQueue struct {
	DeviceId string `json:"deviceId"`
	UserId   string `json:"userId"`
}

type reqNotifyRoomInfo struct {
	Devices []string `json:"devices"`
	Conf    string   `json:"conf"`
}

type reqDeviceStatus struct {
	DeviceId string `json:"deviceId"`
}

type resDeviceStatus struct {
	DeviceId  string `json:"deviceId"`
	Status    int    `json:"status"` //1游戏中 2未开始
	PlayerNum int    `json:"playerNum"`
}

func (this *NetHttp) GetGameLog(w http.ResponseWriter, r *http.Request) {
	//id := r.URL.Query().Get("roomKey")
	// gameLog, err := mongo.GetGameLog(id)
	// if err != nil {
	// 	writeResp(w, err.Error(), nil)
	// } else {
	// 	writeResp(w, "", gameLog)
	// }

}

func readBody(reader io.ReadCloser, v interface{}) string {
	var errMsg = ""
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		errMsg = "read failed!"
	}
	if err := json.Unmarshal(body, v); err != nil {
		errMsg = err.Error()
	}
	return errMsg
}

func writeResp(w http.ResponseWriter, errMsg, ret interface{}) {
	if errMsg != "" {
		data, _ := json.Marshal(map[string]interface{}{
			"code": -1,
			"msg":  errMsg,
		})
		w.Write(data)
	} else {
		data, _ := json.Marshal(map[string]interface{}{
			"code": http.StatusOK,
			"msg":  "Success",
			"data": ret,
		})
		w.Write(data)
	}
}

func checkSign(r *http.Request) string {
	signature := r.Header.Get("signature")
	timestamp := r.Header.Get("timestamp")
	appKey := r.Header.Get("appKey")
	if len(signature) == 0 || len(timestamp) == 0 || len(appKey) == 0 {
		return "auth data required!"
	}

	intTime, err := strconv.ParseInt(timestamp, 10, 32)
	if intTime+3*60 < time.Now().Unix() || err != nil || intTime > time.Now().Unix() {
		return "auth time invalid!"
	}

	//todo use appsecret from config
	md5Str := fmt.Sprintf("appkey=%s?timestamp=%s?appsecret=%s", appKey, timestamp, "whateverxxxxxx")
	sign := util.Md5(md5Str)
	if sign != signature {
		return "signature not invalid!"
	}
	return ""
}
