package service

import (
	"gameserver-997/server/base/utils"
	"gameserver-997/server/domain/match"
)

func StartMatchService() {
	utils.GlobalObject.CustomData["MatchService"] = &MatchService{}
	utils.GlobalObject.CustomData["MatchService"].(*MatchService).Init()
}

type MatchService struct {
	MatchPool *match.MatchPool
}

func (this *MatchService) Init() {
	this.MatchPool = &match.MatchPool{}
	this.MatchPool.Start(0)
}
