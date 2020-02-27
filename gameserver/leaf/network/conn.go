package network

import (
	"net"
)

type Conn interface {
	ReadPkg() (*Package, error)
	WritePkg(p *Package) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
}