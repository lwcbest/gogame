package gamedata

import (
	"reflect"

	"github.com/lwcbest/gogame/gameserver/leaf/log"
	"github.com/lwcbest/gogame/gameserver/leaf/recordfile"
)

func readRf(st interface{}) *recordfile.RecordFile {
	rf, err := recordfile.New(st)
	if err != nil {
		log.Fatal("%v", err)
	}
	fn := reflect.TypeOf(st).Name() + ".txt"
	err = rf.Read("gamedata/" + fn)
	if err != nil {
		log.Fatal("%v: %v", fn, err)
	}

	return rf
}
