package fnet

import (
	"io"
	// "gameserver-997/server/base/iface"
)

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
	connection IConn
}

func BuildPackage(pType PackageType, msgLen int, msgData []byte) *Package {
	return &Package{
		pType,
		msgLen,
		msgData,
		nil,
	}
}

func ReadPackage(conn IConn) (*Package, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}
	pkgType := header[0]
	bodyLength := int32(header[1])<<16 | int32(header[2])<<8 | int32(header[3])
	body := make([]byte, bodyLength)
	if _, err := io.ReadFull(conn, body); err != nil {
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
