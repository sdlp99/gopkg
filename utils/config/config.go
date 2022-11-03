package config

import (
	"github.com/sdlp99/sdpkg/utils/gm"
	"github.com/sdlp99/sdpkg/utils/str"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	gConfig string = "{}"
	gPath   string = "./"
)

func GetPath() string {
	return gPath
}

func init() {

	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	gConfig = string(file)

	enStr := gjson.Get(gConfig, "enSM4").Str
	str1 := strings.Split(enStr, ",")
	for _, val2 := range str1 {
		gConfig, err = sjson.Set(gConfig, val2,
			gm.EnSM4(gjson.Get(gConfig, val2).Str))
	}
	ioutil.WriteFile("./config.json", str.S2b(gConfig), os.ModeAppend)

	for _, val2 := range str1 {
		gConfig, err = sjson.Set(gConfig, val2,
			gm.DeSM4(gjson.Get(gConfig, val2).Str))
	}
}

func GetConfig(key string, defaultVal string) string {
	res := gjson.Get(gConfig, key)
	if res.Value() == nil {
		return defaultVal
	}
	return res.String()

}

func GetConfigInt(key string, defaultVal int) int {
	res := gjson.Get(gConfig, key)
	if res.Value() == nil {
		return defaultVal
	}
	return int(gjson.Get(gConfig, key).Int())
}
