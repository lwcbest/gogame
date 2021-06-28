package fnet

import (
	"fmt"
	"gameserver-997/server/base/logger"

	// "gameserver-997/server/base/utils"
	"reflect"
	// "strconv"
	"errors"
	"math"

	"github.com/golang/protobuf/proto"
	// "runtime/debug"
)

const MSG_ROUTE_LIMIT = 255

type MessageType int32

const (
	MSG_REQUEST  MessageType = 0
	MSG_NOTIFY   MessageType = 1
	MSG_RESPONSE MessageType = 2
	MSG_PUSH     MessageType = 3
)

type Message struct {
	msgType MessageType
	msgId   uint
	route   string
	data    []byte
}

func BuildMsg(mType MessageType, reqId uint, route string, data proto.Message) *Message {
	var dataByte []byte
	if data == nil {

	} else {
		var err error
		dataByte, err = proto.Marshal(data)
		if err != nil {
			fmt.Printf("proto marshal %v error: %v\n", reflect.TypeOf(data), err)
			panic(err)
		}
	}

	return &Message{
		msgType: mType,
		msgId:   reqId,
		route:   route,
		data:    dataByte,
	}
}

func MsgEncode(msg *Message) ([]byte, error) {
	msgType := msg.msgType
	route := msg.route
	id := msg.msgId
	msgData := msg.data
	routeLength := len([]byte(route))
	if routeLength > MSG_ROUTE_LIMIT {
		logger.Error("route is too long...")
		//todo return err
	}

	var bys []byte
	idLen := 0
	switch msgType {
	case MSG_REQUEST, MSG_RESPONSE:
		bys = encodeUInt32(id)
		idLen = len(bys)
	case MSG_NOTIFY, MSG_PUSH:
		break
	}

	totalLen := 1 + idLen + 1 + routeLength + len(msgData)
	head := make([]byte, totalLen)
	flag := byte(msgType)
	head[0] = flag
	offset := 1
	if idLen > 0 {
		copy(head[offset:], bys)
		offset += idLen
	}

	head[offset] = byte(routeLength)
	offset++
	numCopy := copy(head[offset:], []byte(route))
	offset += numCopy
	if offset != totalLen-len(msgData) {
		return nil, errors.New("headlength not match")
	}
	copy(head[offset:], msgData)
	return head, nil
}

func MsgDecode(bytes []byte) *Message {
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
	body := bytes[offset:] //body这里复用一个底层数组应该是ok的

	msg := Message{
		msgType: msgType,
		route:   route,
		msgId:   uint(id),
		data:    body,
	}
	return &msg
}

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
