package network

type HandshakeJson struct {
	Code int `json:"code"`
	Sys  Sys `json:"sys"`
}
type Sys struct {
	Heartbeat int    `json:"heartbeat"`
	Protos    Protos `json:"protos"`
}
type Protos struct {
	Req  ProtoBuf `json:"req"`
	Res  ProtoBuf `json:"res"`
	Push ProtoBuf `json:"push"`
}

type ProtoBuf struct {
	Name string
	Content string
}
