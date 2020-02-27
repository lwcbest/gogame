package network

import (
	"io"

	"github.com/lwcbest/gogame/gameserver/leaf/log"
)

type PkgParser struct {
}

func NewPkgParser() *PkgParser {
	p := new(PkgParser)
	return p
}

/// goroutine safe
/// -------header-------|-----body------
/// |  type  |  length  |     msg      |
/// | 1byte  |  3bytes  | length bytes |
/// --------------------|---------------
///
/// ---------------------type----------------------------------
/// | ----0001   | ----0010 | ----0011  | ----0100 | ----0101 |
/// |  handshake |  hs ack  | heartbeat |   data   |   kick   |
/// -----------------------------------------------------------
func (p *PkgParser) Read(conn *TCPConn) (*Package, error) {
	header := make([]byte, 4)

	// read header
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}
	pkgType := header[0]
	bodyLength := header[1]<<16 | header[2]<<8 | header[3]
	log.Debug("2.read header fullend?:", header, pkgType, bodyLength)
	body := make([]byte, bodyLength)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, err
	}
	log.Debug("3.read body", body)

	pkg := &Package{
		pkgType: PackageType(pkgType),
		length:  int(bodyLength),
		body:    body,
	}

	log.Debug("4.pkg= ", pkg)
	return pkg, nil
}

// goroutine safe
func (p *PkgParser) Write(conn *TCPConn, pkg *Package) error {
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

	for {
		if index >= length || body == nil {
			break
		}
		buf[index] = body[index-4]
		index++
	}

	conn.Write(buf)
	return nil
}

type PackageType int32

const (
	PKG_HANDSHAKE     PackageType = 1
	PKG_HANDSHAKE_ACK PackageType = 2
	PKG_HEARTBEAT     PackageType = 3
	PKG_DATA          PackageType = 4
	PKG_KICK          PackageType = 5
)

type Package struct {
	pkgType PackageType
	length  int
	body    []byte
}
