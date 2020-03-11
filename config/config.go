package config

import (
	"github.com/shybily/go-utils"
	"github.com/sirupsen/logrus"
	"github.com/yookoala/realpath"
	"gopkg.in/ini.v1"
	"os"
)

var Config *ini.Section

func LoadConfig(file string, section string) {
	filePath, err := realpath.Realpath(file)
	if err != nil || !utils.FileExists(filePath) {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"file":    file,
			"section": section,
		}).Fatal("load config file failed")
		os.Exit(1)
	}
	iniConf, err := ini.Load(filePath)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("load env file failed")
		os.Exit(1)
	}
	Config = iniConf.Section(section)
	logrus.WithFields(logrus.Fields{"file_path": filePath, "section": section}).Info("load config file success")
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
