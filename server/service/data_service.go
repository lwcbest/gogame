package service

import (
	"gameserver-997/server/base/logger"
	"gameserver-997/server/base/utils"
	"gameserver-997/server/dao/mongo"
)

func StartDataService() {
	utils.GlobalObject.CustomData["DataService"] = &DataService{}
	utils.GlobalObject.CustomData["DataService"].(*DataService).Init()
}

func GetDataService() *DataService {
	return utils.GlobalObject.CustomData["DataService"].(*DataService)
}

type DataService struct {
	M *mongo.MyMongo
}

func (this *DataService) Init() {
	logger.Info("[start data service!!!!]")
	this.M = mongo.InitMongo()
}
