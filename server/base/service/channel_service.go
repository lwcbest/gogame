package service

import (
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"sync"
)

type ChannelService struct {
	RootCluster iface.Icluster
	RootServer  iface.Iserver
	Channels    sync.Map
}

func (this *ChannelService) GetRootServer() iface.Iserver {
	return this.RootServer
}

func (this *ChannelService) GetRootCluster() iface.Icluster {
	return this.RootCluster
}

func (this *ChannelService) AddChannel(name string, channel iface.IChannel) {
	this.Channels.Store(name, channel)
}

func (this *ChannelService) GetChannel(name string) iface.IChannel {
	channel := &Channel{
		Name:           name,
		UserCount:      0,
		channelService: this,
	}
	cha, _ := this.Channels.LoadOrStore(name, channel)
	return cha.(iface.IChannel)
}

func (this *ChannelService) DestroyChannel(name string) {
	this.Channels.Delete(name)
}

func (this *ChannelService) PushMsgByUids(route string, msg interface{}, sessionArray []iface.ISession) {
	//rpc send serverId,msg
	serverUserMap := make(map[string][]string)
	for _, session := range sessionArray {
		uid := session.GetUid()
		serverId := session.GetServerId()
		if serverUserMap[serverId] == nil {
			serverUserMap[serverId] = make([]string, 0)
		}
		serverUserMap[serverId] = append(serverUserMap[serverId], uid)
	}

	for serverId, uids := range serverUserMap {
		newArgs := []interface{}{uids, route, msg}
		this.RootCluster.RpcPushServerName(serverId, "PushMessage", newArgs...)
	}
}

type Channel struct {
	Name           string
	Groups         sync.Map // group map for uids. key: feServerid, value: [uid,uid]
	UidMap         sync.Map
	UserCount      int
	channelService *ChannelService
}

func (this *Channel) Add(uid string, feServerId string) {
	group, have := this.Groups.LoadOrStore(feServerId, []string{uid})
	g := group.([]string)
	if have {
		this.Groups.Store(feServerId, append(g, uid))
	}

	this.UidMap.Store(uid, feServerId)
	this.UserCount++
}

func (this *Channel) Leave(uid string, feServerId string) {
	group, _ := this.Groups.Load(feServerId)
	if group == nil {
		return
	}
	g := group.([]string)
	logger.Info("channelLeave user %s %s", uid, feServerId)
	for i := 0; i < len(g); i++ {
		if g[i] == uid {
			this.Groups.Store(feServerId, append(g[0:i], g[i+1:]...))
			break
		}
	}

	has := false
	for i := 0; i < len(g); i++ {
		if g[i] == uid {
			has = true
			break
		}
	}
	if !has {
		this.UidMap.Delete(uid)
		this.UserCount--
	}
}

func (this *Channel) GetMember(uid string) string {
	v, _ := this.UidMap.Load(uid)
	return v.(string)
}

func (this *Channel) Destroy() {
	//防止channel本身还在被引用继续发消息
	this.UidMap = sync.Map{}
	this.Groups = sync.Map{}
	this.channelService.DestroyChannel(this.Name)
}

func (this *Channel) PushMessage(route string, args ...interface{}) {
	logger.Error("Need implement~~~")
	// this.Groups.Range(func(serverId, group) {
	// 	newArgs := append([]interface{}{group, route}, args...)
	// 	logger.Info("push message in channel", serverId, route, newArgs)
	// 	err := this.channelService.RootCluster.RpcPushServerName(serverId, "PushMessage", newArgs...)
	// 	if err != nil {
	// 		logger.Error("PushMessage failed!", route, serverId, group, err)
	// 	}
	// }, end)
}
