package config

import (
	"errors"
	"github.com/shybily/go-utils"
	"github.com/yookoala/realpath"
	"gopkg.in/ini.v1"
)

var Config *ini.Section

func LoadConfig(file string, section string) error {
	filePath, err := realpath.Realpath(file)
	if err != nil || !utils.FileExists(filePath) {
		return errors.New("file not found")
	}
	iniConf, err := ini.Load(filePath)
	if err != nil {
		return err
	}
	Config = iniConf.Section(section)
	return nil
}

func Val(key string) string {
	return Config.Key(key).Value()
}

func Int(key string) int {
	res, _ := Config.Key(key).Int()
	return res
}

func Bool(key string) bool {
	res, _ := Config.Key(key).Bool()
	return res
}

func Int64(key string) int64 {
	res, _ := Config.Key(key).Int64()
	return res
}
