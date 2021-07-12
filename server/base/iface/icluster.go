package iface

type Icluster interface {
	RpcCallServerId(session ISession, serverId string, target string, args ...interface{}) (map[string]interface{}, error)
	RpcRandomCallServerType(session ISession, serverType string, target string, args ...interface{}) (map[string]interface{}, error)
	RpcSystemCallServerType(session ISession, serverType string, target string, args ...interface{}) (map[string]interface{}, error)
	RpcPushServerId(session ISession, serverType string, target string, args ...interface{}) error
	RpcPushServerName(serverName string, target string, args ...interface{}) error
}
