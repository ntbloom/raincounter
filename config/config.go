// Package config sets and retrieves configuration values
package config

import (
	"github.com/ntbloom/rainbase/pkg/config/configkey"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// config files
const (
	configDir     = "/etc/rainbase/"
	mainConfig    = "rainbase"
	secretsConfig = "secrets"
)

// Configure process config files and set log level
func Configure() {
	// get the base config
	viper.SetConfigName(mainConfig)
	viper.AddConfigPath(configDir)
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("config not loaded: %s", err)
	}

	// bring in secrets
	viper.SetConfigName(secretsConfig)
	err = viper.MergeInConfig()
	if err != nil {
		logrus.Fatal("secrets not loaded")
	}

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
	logrus.SetLevel(lev)
	logrus.Infof("logger set to %s level", logrus.GetLevel())
}
