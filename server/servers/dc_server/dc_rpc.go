package dc_server

import (
	"gameserver-997/server/base/iface"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/domain/entity"
	"gameserver-997/server/service"

	"gopkg.in/mgo.v2/bson"
)

type DcRpc struct{}

func (this *DcRpc) GetPlayerInfo(request *iface.CommonRequest) map[string]interface{} {
	logger.Info("[DcRpc][GetPlayerInfo]%v", request.RpcData)
	resp := make(map[string]interface{})
	player := request.RpcData[0].(*entity.Player)
	dataService := service.GetDataService()
	row := make(map[string]interface{})
	err := dataService.M.FindOne("player", bson.M{"username": player.Username, "pwd": player.Pwd}, row)
	if err != nil {
		resp["err"] = err.Error()
		return resp
	}

	resp["player"] = &entity.Player{
		Uid:      row["_id"].(bson.ObjectId).String(),
		Username: row["username"].(string),
		Pwd:      row["pwd"].(string),
		Name:     row["name"].(string),
		AvaUrl:   row["avaurl"].(string),
		Score:    int32(row["score"].(int)),
		Level:    int32(row["level"].(int)),
	}

	return resp
}

func (this *DcRpc) SavePlayerInfo(request *iface.CommonRequest) map[string]interface{} {
	resp := make(map[string]interface{})
	player := request.RpcData[0].(*entity.Player)
	dataService := service.GetDataService()
	err := dataService.M.Upsert("player", bson.M{"Uid": player.Uid}, player)
	if err != nil {
		resp["err"] = err.Error()
	}

	resp["player"] = player

	return resp
}
