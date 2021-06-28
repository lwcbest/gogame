package iface

type ISessionService interface {
	Create(string, string, Iconnection) ISession
	Bind(string, string)
	UnBind(string, string)
	Get(string) ISession
	Remove(string)
	ImportData(string, string, interface{})
	KickBySId(string, string)
	SendData(string, []byte)
	SendDataByUid(string, []byte)
	PushMsgByUids([]string, IMessage)
	GetRootCluster() Icluster
	GetRootServer() Iserver
}

type IChannelService interface {
	AddChannel(name string, channel IChannel)
	GetChannel(name string) IChannel
	DestroyChannel(name string)
	PushMsgByUids(route string, msg interface{}, uidMap map[string]string)
	GetRootCluster() Icluster
	GetRootServer() Iserver
}

type IChannel interface {
	Add(uid string, feServerId string)
	Leave(uid string, feServerId string)
	GetMember(uid string) string
	Destroy()
	PushMessage(route string, args ...interface{})
	PushMessageByUids(uids []string, route string, args ...interface{})
}

type ISession interface {
	GetServerId() string
	GetId() string
	GetUid() string
	Bind(string)
	UnBind(string)
	Set(string, interface{})
	Remove(string)
	Get(string) interface{}
	Send([]byte) error
	SendBatch([][]byte) bool
	Closed(string)
	HeartBeat()
	BackendSession() interface{}
	Kick(string)
	RemoteIp() string
}

type IBackendSession interface {
	GetServerId() string
	GetUid() string
}

type MessageType int32

type IMessage interface {
	GetMsgType() MessageType
	GetMsgId() uint
	GetRoute() string
	GetData() []byte
}
