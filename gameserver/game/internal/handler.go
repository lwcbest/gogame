package internal

import (
	"myGo/gameserver/leaf/log"
	"myGo/gameserver/leaf/network"
	"myGo/gameserver/msg"
	"reflect"
)

func init(){
	handler(&msg.Hello{},handleHello)
}

func handler(m interface{},h interface{}){
	skeleton.RegisterChanRPC(reflect.TypeOf(m),h)
}

func handleHello(args []interface{}){
	m:=args[0].(*msg.Hello)
	session:=args[1].(network.Session)
	log.Debug("hello %v",m.GetName(),session)
	//session.WriteMsg(&msg.Hello{
	//	Name:proto.String("client"),
	//})
}