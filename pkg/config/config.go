// Package config sets and retrieves configuration values
package config

import (
	"path"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// config files
const (
	configDir  = "/etc/raincounter/"
	mainConfig = "insecure"
)

var (
	RegularFile string
)

// Configure process config files and set log level
func Configure() {
	logrus.Info("Pulling in viper config")
	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}

	var (
		regular   string
		directory string
	)
	if RegularFile == "" {
		regular = mainConfig
		directory = configDir
	} else {
		regular = path.Base(RegularFile)
		directory = path.Dir(RegularFile)
	}
	getConfigFiles(regular, directory)

	// set the log level
	SetLogger()
}

// SetLogger sets the log level from config or a level if you specify
func SetLogger(level ...string) {
	params := len(level)
	var lev logrus.Level
	var err error

	switch params {
	case 0:
		lev, err = logrus.ParseLevel(viper.GetString(configkey.Loglevel))
		if err != nil {
			panic(err)
		}
	case 1:
		lev, err = logrus.ParseLevel(level[0])
		if err != nil {
			panic(err)
		}
	default:
		panic("specify a log level in config file")
	}
	logrus.SetReportCaller(true)
	logrus.SetLevel(lev)
	logrus.Infof("logger set to %s level", logrus.GetLevel())
}

func getConfigFiles(regular, directory string) {
	logrus.Infof("using config=%s, directory=%s", regular, directory)
	viper.SetConfigType("yml")
	viper.AddConfigPath(directory)
	viper.SetConfigName(regular)
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("config not loaded: %s", err)
	}
}
