package utils

import (
	// "fmt"
	"errors"
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"io"
	"math"
	"reflect"

	"github.com/golang/protobuf/proto"
)

//--------------------------------------------package-----------------------------------------------------

type PackageType int32

const (
	PKG_HANDSHAKE     PackageType = 1
	PKG_HANDSHAKE_ACK PackageType = 2
	PKG_HEARTBEAT     PackageType = 3
	PKG_DATA          PackageType = 4
	PKG_KICK          PackageType = 5
)

/// -------header-------|-----body------
/// |  type  |  length  |     msg      |
/// | 1byte  |  3bytes  | length bytes |
/// --------------------|---------------
///
/// ---------------------type----------------------------------
/// | ----0001   | ----0010 | ----0011  | ----0100 | ----0101 |
/// |  handshake |  hs ack  | heartbeat |   data   |   kick   |
/// -----------------------------------------------------------
type Package struct {
	pkgType    PackageType
	length     int
	body       []byte
	connection iface.Iconnection
}

func (this *Package) GetPkgType() PackageType {
	return this.pkgType
}

func (this *Package) GetConnection() iface.Iconnection {
	return this.connection
}

func (this *Package) GetData() []byte {
	return this.body
}

func BuildPackage(pType PackageType, msgLen int, msgData []byte) *Package {
	return &Package{
		pType,
		msgLen,
		msgData,
		nil,
	}
}

func ReadPackage(conn iface.Iconnection) (*Package, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(conn.GetConnection(), header); err != nil {
		return nil, err
	}
	pkgType := header[0]
	bodyLength := int32(header[1])<<16 | int32(header[2])<<8 | int32(header[3])
	body := make([]byte, bodyLength)
	if _, err := io.ReadFull(conn.GetConnection(), body); err != nil {
		return nil, err
	}

	pkg := &Package{
		pkgType:    PackageType(pkgType),
		length:     int(bodyLength),
		body:       body,
		connection: conn,
	}
	return pkg, nil
}

func WritePackage(pkg *Package) []byte {
	pkgType := pkg.pkgType
	body := pkg.body
	length := 4
	if body != nil {
		length += len(body)
	}

	buf := make([]byte, length)
	index := 0
	buf[index] = byte(pkgType)
	index++
	buf[index] = byte(len(body) >> 16 & 0xFF)
	index++
	buf[index] = byte(len(body) >> 8 & 0xFF)
	index++
	buf[index] = byte(len(body) & 0xFF)
	index++

	copy(buf[4:], body)
	return buf
}

func BuildPkg(pkgType PackageType, data []byte) *Package {
	length := len(data)
	pkg := &Package{pkgType, length, data, nil}
	return pkg
}

//------------------------------------------------------------------message-------------------------------------------------------
const MSG_ROUTE_LIMIT = 255

const (
	MSG_REQUEST  iface.MessageType = 0
	MSG_NOTIFY   iface.MessageType = 1
	MSG_RESPONSE iface.MessageType = 2
	MSG_PUSH     iface.MessageType = 3
)

type Message struct {
	msgType iface.MessageType
	msgId   uint
	route   string
	data    []byte
}

//为了使用接口这里需要包装成方法
func (this *Message) GetMsgType() iface.MessageType {
	return this.msgType
}

func (this *Message) GetMsgId() uint {
	return this.msgId
}

func (this *Message) GetRoute() string {
	return this.route
}

func (this *Message) GetData() []byte {
	return this.data
}

func BuildMsg(mType iface.MessageType, reqId uint, route string, data proto.Message) *Message {
	dataByte, err := proto.Marshal(data)
	if err != nil {
		logger.Error("proto marshal %v error: %v", reflect.TypeOf(data), err)
	}

	return &Message{
		msgType: mType,
		msgId:   reqId,
		route:   route,
		data:    dataByte,
	}
}

func BuildMsgFromData(mType iface.MessageType, reqId uint, route string, data []byte) *Message {
	return &Message{
		msgType: mType,
		msgId:   reqId,
		route:   route,
		data:    data,
	}
}

func MsgEncode(msg iface.IMessage) ([]byte, error) {
	msgType := msg.GetMsgType()
	route := msg.GetRoute()
	id := msg.GetMsgId()
	msgData := msg.GetData()
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

func MsgDecode(bytes []byte) (*Message, error) {
	id := uint32(0)
	offset := 1
	msgType := iface.MessageType(bytes[0])
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
	return &msg, nil
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
