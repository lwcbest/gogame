package fnet

import (
	"io"
	"net"
	"time"
)

type IConn interface {
	io.Reader
	io.Writer
	Close() error
	RemoteAddr() net.Addr
	SetReadDeadline(t time.Time) error
}