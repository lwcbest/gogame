package main

import (
	"fmt"
	"io"
	"os"

	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	gopb "gameserver-997/pb/gopb"
	"gameserver-997/server/base/logger"
	"gameserver-997/unit_test/fnet"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

//封装网络
type TcpClient struct {
	connection fnet.IConn
	addr       string
	stopChan   chan struct{}
}

func (this *TcpClient) Start() {

}

//封装hangler
type Handler struct {
	Apis     map[uint32]reflect.Value
	DataPack *fnet.PBDataPack
}

func (this *Handler) AddRouter(router interface{}) {
	value := reflect.ValueOf(router)
	tp := reflect.TypeOf(router)
	for i := 0; i < value.NumMethod(); i++ {
		name := tp.Method(i).Name
		if name[0:2] != "On" {
			continue
		}
		name = name[2:]
		k := strings.Split(name, "_")
		index, err := strconv.Atoi(k[1])
		if err != nil {
			panic("error api: " + name)
		}
		this.Apis[uint32(index)] = value.Method(i)
		fmt.Println("add router ", name)
	}
}

func (this *Handler) HandleClient(client *TcpClient) {
	for {
		headdata := make([]byte, this.DataPack.GetHeadLen())
		if _, err := io.ReadFull(client.connection, headdata); err != nil {
			fmt.Println("error:", err.Error())
			close(client.stopChan)
			return
		}
		pkgHead, err := this.DataPack.Unpack(headdata)
		fmt.Println("pkgHead: ", pkgHead)
		if err != nil {
			fmt.Println("error:", err.Error())
			close(client.stopChan)
			return
		}
		//data
		pkg := pkgHead.(*fnet.PkgData)
		if pkg.Len > 0 {
			pkg.Data = make([]byte, pkg.Len)
			if _, err := io.ReadFull(client.connection, pkg.Data); err != nil {
				fmt.Println("error:", err.Error())
				close(client.stopChan)
				return
			}
		}
		fmt.Println("pkg: ", pkg)
		if f, ok := this.Apis[pkg.MsgId]; ok {
			f.Call([]reflect.Value{reflect.ValueOf(pkg.Data)})
		} else {
			fmt.Println("received message no handler: ", pkg.MsgId)
		}
	}
}

type NewApi struct {
	reqId uint
}

func (this *NewApi) HandShake(connection fnet.IConn) {
	pkg := fnet.BuildPackage(fnet.PKG_HANDSHAKE, 0, nil)
	bytes := fnet.WritePackage(pkg)
	fmt.Println("before write write")
	connection.Write(bytes)
}

var heartTime int64 = 0

func (this *NewApi) HeartBeat(connection fnet.IConn) {
	ticker := time.NewTicker(time.Millisecond * 5000)
	go func() {
		for {
			select {
			case <-ticker.C:
				// if heartTime == 0 {
				// 	heartTime = time.Now().UnixNano() / 1000000
				// } else {
				// 	fmt.Printf("hearttime is not zero")
				// }
				pkg := fnet.BuildPackage(fnet.PKG_HEARTBEAT, 0, nil)
				bytes := fnet.WritePackage(pkg)
				connection.Write(bytes)
			}
		}
	}()
}

func (this *NewApi) GetNetAddr(connection fnet.IConn) {
	this.reqId++
	msg := fnet.BuildMsg(fnet.MSG_REQUEST, this.reqId, "Net_Addr", nil)
	msgData, _ := fnet.MsgEncode(msg)
	pkg := fnet.BuildPackage(fnet.PKG_DATA, len(msgData), msgData)
	bytes := fnet.WritePackage(pkg)
	connection.Write(bytes)
}

func (this *NewApi) OnVersion_Addr(data []byte) {
	msg := &gopb.ResVersionAddr{}
	if err := proto.Unmarshal(data, msg); err == nil {
		fmt.Println("Onnetaddr success\n", msg)
	} else {
		fmt.Printf("Onnetaddr %v\n fauled", data)
	}
}

// func (this *NewApi) Login(connection fnet.IConn, user string) {
// 	info := &gopb.ReqNetLogin{
// 		Uid:   user,
// 		Token: "pass0",
// 	}
// 	this.reqId++
// 	msg := fnet.BuildMsg(fnet.MSG_REQUEST, this.reqId, "Net_Login", info)
// 	msgData, _ := fnet.MsgEncode(msg)
// 	pkg := fnet.BuildPackage(fnet.PKG_DATA, len(msgData), msgData)
// 	bytes := fnet.WritePackage(pkg)
// 	connection.Write(bytes)
// 	fmt.Println("Login time: ", time.Now().UnixNano()/1000)
// }

func GetRouteName(obj interface{}) (string, string) {
	typeString := reflect.TypeOf(obj).String()
	typeString = typeString[5:]
	pfStart := -1
	serverNameStart := -1
	interNameStart := -1
	for i := 0; i < len(typeString); i++ {
		if typeString[i] >= 65 && typeString[i] <= 90 {
			//首字母大写
			if pfStart < 0 {
				pfStart = i
				continue
			}
			if serverNameStart < 0 {
				serverNameStart = i
				continue
			}

			if interNameStart < 0 {
				interNameStart = i
				break
			}
		}
	}

	pfStr := typeString[pfStart:serverNameStart]
	serverStr := typeString[serverNameStart:interNameStart]
	interStr := typeString[interNameStart:]
	routeName := serverStr + "_" + interStr
	fmt.Println("[GetRouteName]: ", routeName)

	return pfStr, routeName
}

func (this *NewApi) Net_Login(connection fnet.IConn, username string, pwd string) {
	logger.Fatal("start login test: ", time.Now().UnixNano()/1000)
	info := &gopb.ReqNetLogin{Username: username, Password: pwd}
	_, routeName := GetRouteName(info)
	this.reqId++
	msg := fnet.BuildMsg(fnet.MSG_REQUEST, this.reqId, routeName, info)
	msgData, _ := fnet.MsgEncode(msg)
	pkg := fnet.BuildPackage(fnet.PKG_DATA, len(msgData), msgData)
	bytes := fnet.WritePackage(pkg)
	connection.Write(bytes)
	logger.Fatal("login test: ", time.Now().UnixNano()/1000)
}

func (this *NewApi) Net_EnterMatchQueue(connection fnet.IConn) {
	logger.Info("Net_EnterMatchQueue: ", time.Now().UnixNano()/1000)
	info := &gopb.ReqNetEnterMatchQueue{Level: 1}
	_, routeName := GetRouteName(info)
	this.reqId++
	msg := fnet.BuildMsg(fnet.MSG_REQUEST, this.reqId, routeName, info)
	msgData, _ := fnet.MsgEncode(msg)
	pkg := fnet.BuildPackage(fnet.PKG_DATA, len(msgData), msgData)
	bytes := fnet.WritePackage(pkg)
	connection.Write(bytes)
	logger.Info("Net_EnterMatchQueue: ", time.Now().UnixNano()/1000)
}

func (this *NewApi) Net_CreateTimeRoom(connection fnet.IConn) {
	info := &gopb.ReqNetCreateRoom{PlayMod: "richMan", Mac: "18-93-7F-E1-8D-CD"}
	this.reqId++
	msg := fnet.BuildMsg(fnet.MSG_REQUEST, this.reqId, "Net_CreateTimeRoom", info)
	msgData, _ := fnet.MsgEncode(msg)
	pkg := fnet.BuildPackage(fnet.PKG_DATA, len(msgData), msgData)
	bytes := fnet.WritePackage(pkg)
	connection.Write(bytes)
	fmt.Println("CreateRoom time: ", time.Now().UnixNano()/1000)
}

func (this *NewApi) Net_EnterRoom(connection fnet.IConn, uid, roomId string) {
	info := &gopb.ReqNetEnterRoom{Token: uid, DeviceId: roomId}
	this.reqId++
	fmt.Println(this.reqId, "+++++++++++++++++++")
	msg := fnet.BuildMsg(fnet.MSG_REQUEST, this.reqId, "EnterTimeRoom", info)
	msgData, _ := fnet.MsgEncode(msg)
	pkg := fnet.BuildPackage(fnet.PKG_DATA, len(msgData), msgData)
	bytes := fnet.WritePackage(pkg)
	connection.Write(bytes)
}

func (this *NewApi) Game_RandomAward(connection fnet.IConn) {
	info := &gopb.ReqGamRandomAward{Start: 18, Steps: 1}
	this.reqId++
	msg := fnet.BuildMsg(fnet.MSG_REQUEST, this.reqId, "Game_RandomAward", info)
	msgData, _ := fnet.MsgEncode(msg)
	pkg := fnet.BuildPackage(fnet.PKG_DATA, len(msgData), msgData)
	bytes := fnet.WritePackage(pkg)
	connection.Write(bytes)
}

func (this *NewApi) Game_PlayerReady(connection fnet.IConn) {
	this.reqId++
	info := &gopb.ReqGameFPlayerReday{Select: []string{"2_0_1_2"}}
	msg := fnet.BuildMsg(fnet.MSG_REQUEST, this.reqId, "Game_PlayerReady", info)
	msgData, _ := fnet.MsgEncode(msg)
	pkg := fnet.BuildPackage(fnet.PKG_DATA, len(msgData), msgData)
	bytes := fnet.WritePackage(pkg)
	connection.Write(bytes)
}

func (this *NewApi) Game_SyncScore(connection fnet.IConn) {
	info := &gopb.ReqCommonArgs{Args: make([]*gopb.ReqCommonArg, 1)}
	info.Args[0] = &gopb.ReqCommonArg{
		StrVal1: "5f9bd3f7e8a5671dccc580b5",
		IntVal1: 210,
		IntVal2: 0,
		IntVal3: 0,
	}
	msg := fnet.BuildMsg(fnet.MSG_REQUEST, this.reqId, "Game_SyncResult", info)
	msgData, _ := fnet.MsgEncode(msg)
	pkg := fnet.BuildPackage(fnet.PKG_DATA, len(msgData), msgData)
	bytes := fnet.WritePackage(pkg)
	connection.Write(bytes)
}

func (this *NewApi) SendCommand(connection fnet.IConn, command string) {
	this.reqId++
	msg := fnet.BuildMsg(fnet.MSG_REQUEST, this.reqId, command, nil)
	msgData, _ := fnet.MsgEncode(msg)
	pkg := fnet.BuildPackage(fnet.PKG_DATA, len(msgData), msgData)
	bytes := fnet.WritePackage(pkg)
	connection.Write(bytes)
}

func (this *NewApi) SendFrame(connection fnet.IConn, i int32) {
	this.reqId++
	info := &gopb.ReqGameSyncFCommand{EpFrame: i, Uid: "11111", Ctype: 2, ParamList: []int32{1}}
	msg := fnet.BuildMsg(fnet.MSG_REQUEST, this.reqId, "Game_SyncFCommand", info)
	msgData, _ := fnet.MsgEncode(msg)
	pkg := fnet.BuildPackage(fnet.PKG_DATA, len(msgData), msgData)
	bytes := fnet.WritePackage(pkg)
	connection.Write(bytes)
}

var lastSererTime int64
var lastCliTime int64
var frameTime int64 = 0
var loginTime int64 = 0

func (this *NewApi) OnHandleData(route string, data []byte, connection fnet.IConn) {
	var msg proto.Message
	switch route {
	case "Net_Login":
		msg = &gopb.ResNetLogin{}
	case "Net_CreateRoom":
		msg = &gopb.ResNetCreateFRoom{}
	case "Net_CreateStateRoom":
		msg = &gopb.ResNetCreateSRoom{}
	case "Net_EnterStateRoom":
		msg = &gopb.ResNetEnterSRoom{}
	case "Net_CreateTimeRoom":
		msg = &gopb.ResNetCreateSRoom{}

	case "OnPlayerWaiting":
		msg = &gopb.PlayerInfo{}
	case "OnGameEnd":
		msg = &gopb.OnGameEnd{}
	case "OnPlayerEnter":
		msg = &gopb.OnPlayerEnterFRoom{}
	case "OnPlayerClientEnter":
		msg = &gopb.OnPlayerClientEnter{}
	case "OnPlayerReady":
		msg = &gopb.OnPlayerReady{}
	case "OnAllPlayerReady":
		// heartTime = time.Now().UnixNano() / 1000000
		// fmt.Println("send gameoktostart: ", time.Now().UnixNano()/1000000)
		this.SendCommand(connection, "Game_OkToStart")
		// loginTime = time.Now().Unix()
		// this.Game_SyncScore(connection)
		return
	case "OnFrame":
		if loginTime > 0 && time.Now().Unix()-loginTime > 30 {
			loginTime = 0
			this.Game_SyncScore(connection)
		}
		fmt.Printf("frameTime: %d\n", time.Now().UnixNano()/1000000-frameTime)
		frameTime = time.Now().UnixNano() / 1000000
		msg = &gopb.OnFrame{}
	case "HeartBeat":
		// fmt.Printf("heartBeat cost %d\n", time.Now().UnixNano()/1000000-heartTime)
		// heartTime = time.Now().UnixNano() / 1000000
		return
	case "Game_OkToStart":
		msg = &gopb.ResError{}
		fmt.Printf("heartBeat cost %d\n", time.Now().UnixNano()/1000000-heartTime)
		fallthrough
	case "Game_RandomAward":
		msg = &gopb.ResGamRandomAward{}
	default:
		msg = &gopb.ResError{}
	}
	if err := proto.Unmarshal(data, msg); err == nil {
		logger.Info("handledata", route, msg, time.Now().UnixNano()/1000000, data)
	} else {
		logger.Error("handledata error %s %v %+v\n", route, data, err)
	}
}

func main() {
	connection, err := dialWebSocket()
	if err != nil {
		os.Exit(1)
	}

	//---------------------pkg msg两层协议
	protocol := fnet.NewProtocol()
	api := &NewApi{}
	protocol.AddRouter(api)
	ackChan := make(chan struct{})
	go protocol.HandleConnection(connection, ackChan)
	api.HandShake(connection)
	api.HeartBeat(connection)
	if len(os.Args) < 2 {
		logger.Error("need args!!!")
		return
	}
	if os.Args[1] == "0" {
		api.Net_Login(connection, "luoyonglin2", "111222333")
		api.Net_EnterMatchQueue(connection)

	} else if os.Args[1] == "1" {
		api.Net_Login(connection, "luoyonglin3", "111222333")
		api.Net_EnterMatchQueue(connection)
	} else if os.Args[1] == "2" {
		api.Net_Login(connection, "luoyonglin3", "111222333")
		api.Net_EnterMatchQueue(connection)
	} else if os.Args[1] == "3" {

		//api.Net_EnterRoom(connection, "", "30-9C-23-50-CB-73")
	} else {
		runAis()
	}
	// api.Game_SyncState(connection)
	<-ackChan
}

var uids = []string{"5f44762df1f8ed333c281078", "5f448e9af1f8ed333c281079", "5f44d621f1f8ed4b5f0fb26e", "5f44d1aaf1f8ed4b5f0fb26c", "5f44d61ef1f8ed4b5f0fb26d"}

func runAis() {
	for i := 0; i < 7; i++ {
		go runAi(uids[i%5])
	}
	ch := make(chan struct{})
	<-ch
}

func runAi(uid string) {
	for {
		connection, err := dialWebSocket()
		if err != nil {
			return
		}
		protocol := fnet.NewProtocol()
		api := &NewApi{reqId: 100}
		protocol.AddRouter(api)
		ackChan := make(chan struct{})
		go protocol.HandleConnection(connection, ackChan)
		api.HandShake(connection)
		api.HeartBeat(connection)
		api.Net_EnterRoom(connection, uid, "30-9C-23-50-CB-73")
		fmt.Println(api.reqId, "+++++++++++")
		time.Sleep(time.Second * 150)
		connection.Close()
		time.Sleep(time.Second * 50)
	}
}

func dialTcp(addr string) (fnet.IConn, error) {
	servAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println("resolve addr error ", servAddr, err.Error())
		return nil, err
	}
	connection, err := net.DialTCP("tcp", nil, servAddr)
	if err != nil {
		fmt.Println("dial tcp error: ", err.Error(), servAddr)
		return nil, err
	}
	return connection, nil
}

func dialWebSocket() (fnet.IConn, error) {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		fmt.Printf("dial websocket err %+v", err)
		return nil, err
	}
	wsConn := fnet.NewWsConn(conn)
	return wsConn, nil
}