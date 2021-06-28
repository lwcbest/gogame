package fnet

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
)

type WsConn struct {
	conn *websocket.Conn
	buf  []byte
	ip   string
}

func NewWsConn(conn *websocket.Conn, ip string) *WsConn {
	return &WsConn{conn, []byte{}, ip}
}

// 这个方法不符合Reader接口惯例 实际上实现了io.ReadFull接口，这里可以用
func (this *WsConn) Read(data []byte) (n int, err error) {
	readed := 0
	for readed < len(data) {
		if len(this.buf) == 0 {
			_, message, err := this.conn.ReadMessage()
			if err != nil {
				return readed, err
			}
			this.buf = message
		}
		numCopy := copy(data[readed:], this.buf)
		this.buf = this.buf[numCopy:]
		readed += numCopy
		if readed >= len(data) {
			return readed, nil
		}
	}
	return readed, nil
}

func (this *WsConn) Write(data []byte) (n int, err error) {
	if err = this.conn.WriteMessage(websocket.BinaryMessage, data); err == nil {
		return len(data), nil
	} else {
		return 0, err
	}
}

type RemoteAddr struct {
	ip string
}

func (this *RemoteAddr) String() string {
	return this.ip
}

func (this *RemoteAddr) Network() string {
	return "ws"
}

func (this *WsConn) RemoteAddr() net.Addr {
	if this.ip == "" {
		return &RemoteAddr{this.conn.RemoteAddr().String()}
	}
	return &RemoteAddr{this.ip}
}

func (this *WsConn) Close() error {
	this.conn.Close()
	return nil
}

func (this *WsConn) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}

//包含自带方法close不用实现了
