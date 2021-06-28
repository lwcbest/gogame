package fnet

import (
	// "strings"
	"fmt"
	"reflect"
)

type Protocol struct {
	Apis map[string]reflect.Value
}

func NewProtocol() *Protocol {
	return &Protocol{
		Apis: make(map[string]reflect.Value),
	}
}

func (this *Protocol) AddRouter(router interface{}) {
	value := reflect.ValueOf(router)
	tp := reflect.TypeOf(router)
	for i := 0; i < value.NumMethod(); i++ {
		name := tp.Method(i).Name
		if name[0:2] != "On" {
			continue
		}
		name = name[2:]
		// k := strings.Split(name, "_")
		this.Apis[name] = value.Method(i)
		fmt.Println("add router ", name)
	}
}

func (this *Protocol) HandleConnection(fconn IConn, ackChan chan struct{}) {
	for {
		pkg, err := ReadPackage(fconn)
		if err != nil {
			fmt.Println("parse package failed!", err)
			close(ackChan)
			return
		}
		if pkg.pkgType == PKG_HANDSHAKE_ACK {
			fmt.Println("received handshake ack", pkg)
			continue
		}
		if pkg.pkgType == PKG_HEARTBEAT {
			f := this.Apis["HandleData"]
			f.Call([]reflect.Value{reflect.ValueOf("HeartBeat"), reflect.ValueOf([]byte{}), reflect.ValueOf(fconn)})
		}
		if pkg.pkgType != PKG_DATA {
			// fmt.Println("received pkg not data: ", pkg)
			continue
		}
		msg := MsgDecode(pkg.body)
		if msg.msgType != MSG_PUSH && msg.msgType != MSG_RESPONSE {
			fmt.Println("received msg not push and not response: ", msg)
			return
		}
		if f, ok := this.Apis[msg.route]; ok {
			f.Call([]reflect.Value{reflect.ValueOf(msg.data)})
		} else {
			f := this.Apis["HandleData"]
			f.Call([]reflect.Value{reflect.ValueOf(msg.route), reflect.ValueOf(msg.data), reflect.ValueOf(fconn)})
		}
	}
}
