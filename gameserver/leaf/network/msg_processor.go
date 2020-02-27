package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/lwcbest/gogame/gameserver/leaf/chanrpc"
	"github.com/lwcbest/gogame/gameserver/leaf/log"
)

const MSG_Route_Limit = 255

type MessageType int32

const (
	MSG_REQUEST  MessageType = 0
	MSG_NOTIFY   MessageType = 1
	MSG_RESPONSE MessageType = 2
	MSG_PUSH     MessageType = 3
)

type Message struct {
	MsgType MessageType
	Route   string
	Id      uint
	Data    []byte
}

type Processor struct {
	msgInfoMap map[string]*MsgInfo
}

type MsgInfo struct {
	msgType       reflect.Type
	msgRouter     *chanrpc.Server
	msgHandler    MsgHandler
	msgRawHandler MsgHandler
}

type MsgHandler func([]interface{})

type MsgRaw struct {
	msgID      string
	msgRawData json.RawMessage
}

func NewHybridProcessor() *Processor {
	p := new(Processor)
	p.msgInfoMap = make(map[string]*MsgInfo)
	return p
}

// register msg
func (p *Processor) Register(msg interface{}) string {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	msgName := msgType.Elem().Name()
	if msgName == "" {
		log.Fatal("unnamed json message")
	}
	if _, ok := p.msgInfoMap[msgName]; ok {
		log.Fatal("message %v is already registered", msgName)
	}

	info := new(MsgInfo)
	info.msgType = msgType
	p.msgInfoMap[msgName] = info
	return msgName
}

// set router,need register first
func (p *Processor) SetRouter(msg interface{}, msgRouter *chanrpc.Server) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	msgName := msgType.Elem().Name()
	i, ok := p.msgInfoMap[msgName]
	if !ok {
		log.Fatal("message %v not registered", msgName)
	}

	i.msgRouter = msgRouter
}

// set handler
func (p *Processor) SetHandler(msg interface{}, msgHandler MsgHandler) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	msgName := msgType.Elem().Name()
	i, ok := p.msgInfoMap[msgName]
	if !ok {
		log.Fatal("message %v not registered", msgName)
	}

	i.msgHandler = msgHandler
}

func (p *Processor) SetRawHandler(msgName string, msgRawHandler MsgHandler) {
	i, ok := p.msgInfoMap[msgName]
	if !ok {
		log.Fatal("message %v not registered", msgName)
	}

	i.msgRawHandler = msgRawHandler
}

// goroutine safe
func (p *Processor) Route(msg interface{}, userData interface{}) error {
	// raw
	if msgRaw, ok := msg.(MsgRaw); ok {
		i, ok := p.msgInfoMap[msgRaw.msgID]
		if !ok {
			return fmt.Errorf("message %v not registered", msgRaw.msgID)
		}
		if i.msgRawHandler != nil {
			i.msgRawHandler([]interface{}{msgRaw.msgID, msgRaw.msgRawData, userData})
		}
		return nil
	}

	// json
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return errors.New("json message pointer required")
	}
	msgName := msgType.Elem().Name()
	i, ok := p.msgInfoMap[msgName]
	if !ok {
		return fmt.Errorf("message %v not registered", msgName)
	}
	if i.msgHandler != nil {
		i.msgHandler([]interface{}{msg, userData})
	}
	if i.msgRouter != nil {
		i.msgRouter.Go(msgType, msg, userData)
	}
	return nil
}

func (p *Processor) HandlePackage(s *Session, pkg *Package) *Message {
	log.Debug("handle package", pkg)
	switch pkg.pkgType {
	case PKG_HANDSHAKE:
		sendHandShake(s)
	case PKG_HEARTBEAT, PKG_HANDSHAKE_ACK:
		s.heartbeat.Handle()
	case PKG_DATA:
		_, msg := MsgDecode(pkg.body)
		return &msg
	case PKG_KICK:
		log.Release("kick pkg %d", s.sid)
	}

	return nil
}

func sendHandShake(s *Session) {
	handShakeJson := &HandshakeJson{
		Code: 200,
		Sys: Sys{
			Heartbeat: HeartbeatDuration,
			Protos: Protos{
				Req:  ProtoBuf{Name: "req"},
				Res:  ProtoBuf{Name: "res"},
				Push: ProtoBuf{Name: "push"},
			},
		},
	}

	body, _ := json.Marshal(handShakeJson)
	pkg := &Package{
		pkgType: PKG_HANDSHAKE,
		length:  len(body),
		body:    body,
	}
	log.Debug("send hand shake:", pkg)
	s.WriteMsg(pkg)
}

/// goroutine safe
/// ---------------header------------------------------
/// |  flag  | message id  | routeLength |   route    |
/// | 1byte  |  0-4bytes   |    1byte    | 0-256bytes |
/// ---------------------------------------------------
///
/// ---------------------flag--------------------
/// | -----000 | -----001 | -----010 | -----011 |
/// |  request |  notify  | response |   push   |
/// ---------------------------------------------
func MsgEncode(msg Message) (error, []byte) {
	msgType := msg.MsgType
	route := msg.Route
	id := msg.Id
	msgData := msg.Data
	routeLength := len([]byte(route))
	if routeLength > MSG_Route_Limit {
		log.Error("route is too long...")
	}

	//1.encode head
	var bys []byte
	idLen := 0
	if id > 0 {
		bys = encodeUInt32(id)
		idLen = len(bys)
	}

	totolLen := 1 + idLen + 1 + routeLength
	head := make([]byte, totolLen)
	flag := byte(msgType)
	head[0] = flag
	offset := 1
	if id > 0 {
		head, offset = writeBytes(head, offset, bys)
	}

	head[offset] = byte(routeLength)
	offset++
	head, offset = writeBytes(head, offset, []byte(route))

	//2.encode body
	body := msgData
	if offset != len(head) {
		fmt.Println("ji suan cuo wu@@!!")
		return errors.New("ji suan cuo wu"), nil
	}

	result := make([]byte, offset+len(body))
	for i := 0; i < offset; i++ {
		result[i] = head[i]
	}
	for i := 0; i < len(body); i++ {
		result[offset+i] = body[i]
	}

	return nil, result
}

//base 128 varints
func encodeUInt32(n uint) []byte {
	byteList := make([]byte, 0)
	for n != 0 {
		tmp := n % 128
		next := n >> 7
		if next != 0 {
			tmp = tmp + 128
		}
		byteList = append(byteList, byte(tmp))
		n = next
	}
	return byteList
}

func writeBytes(buffer []byte, offset int, bytes []byte) ([]byte, int) {
	for i := 0; i < len(bytes); i++ {
		buffer[offset] = bytes[i]
		offset++
	}
	return buffer, offset
}

/// ---------------header------------------------------
/// |  flag  | message id  | routeLength |   route    |
/// | 1byte  |  0-4bytes   |    1byte    | 0-256bytes |
/// ---------------------------------------------------
///
/// ---------------------flag--------------------
/// | -----000 | -----001 | -----010 | -----011 |
/// |  request |  notify  | response |   push   |
/// ---------------------------------------------
func MsgDecode(bytes []byte) (error, Message) {
	id := uint32(0)
	offset := 1
	msgType := MessageType(bytes[0])
	switch msgType {
	case MSG_REQUEST, MSG_RESPONSE:
		_id, idLength := decodeUInt32(bytes, offset)
		offset += idLength
		id = _id
	case MSG_NOTIFY, MSG_PUSH:
		break
	}

	routeLength := int(bytes[offset])
	offset++
	route := string(bytes[offset : offset+routeLength])
	offset += routeLength

	//Decode body
	body := make([]byte, len(bytes)-offset)
	for i := 0; i < len(body); i++ {
		body[i] = bytes[i+offset]
	}

	msg := Message{
		MsgType: msgType,
		Route:   route,
		Id:      uint(id),
		Data:    body,
	}
	return nil, msg
}

func decodeUInt32(data []byte, offset int) (uint32, int) {
	n := uint32(0)
	length := 0
	for i := offset; i < len(data); i++ {
		length++
		m := uint32(data[i])
		n = n + ((m & 0x7f) * uint32(math.Pow(2, float64(7*(i-offset)))))
		if m < 128 {
			break
		}
	}
	return n, length
}
