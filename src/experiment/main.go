package main

import (
	"assert"
	"config"
	"database/sql"
	"fmt"
	"mr"
	"redis"
	"time"
	"utils"
)

func main() {
	config.Initialize()

	aredis.Initialize()

	tmp := mr.IP2Location{
		IPFrom:      100,
		IPTo:        100,
		CityName:    sql.NullString{"Sss", true},
		CountryCode: sql.NullString{"ssss", true},
		CountryName: sql.NullString{"ssssssss", true},
		RegionName:  sql.NullString{"Sswdes", true},
	}

	b, err := utils.InterfaceToByte(tmp)
	assert.Nil(err)
	assert.Nil(aredis.StoreKey("DATATATA", string(b), time.Minute))
	s, err := aredis.GetKey("DATATATA", false, 0)
	assert.Nil(err)
	var tmp2 mr.IP2Location

	fmt.Println(string(b))
	fmt.Println(s)
	assert.Nil(utils.ByteToInterface([]byte(s), &tmp2))

}
