package fnet

import (
	_ "gameserver-997/server/base/logger"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

type WsConn struct {
	conn *websocket.Conn
	buf  []byte
}

func NewWsConn(conn *websocket.Conn) *WsConn {
	return &WsConn{conn, []byte{}}
}

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

func (this *WsConn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *WsConn) Close() error {
	err := this.conn.Close()
	return err
}

func (this *WsConn) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}
