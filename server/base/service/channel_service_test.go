package service

import (
	"gameserver-997/server/base/iface"
	"testing"
)

func TestChannelService_PushMsgByUids(t *testing.T) {
	channelService := &ChannelService{}
	session1 := &Session{Uid: "user1", FrontendId: "server1"}
	session2 := &Session{Uid: "user2", FrontendId: "server2"}
	session3 := &Session{Uid: "user3", FrontendId: "server1"}
	channelService.PushMsgByUids("abc", "this is msg", []iface.ISession{session1, session2, session3})
}