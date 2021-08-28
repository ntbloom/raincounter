package main

import (
	_ "github.com/ntbloom/raincounter/pkg/common/mqtt"
)

//// connect to mqtt
//func connectToMQTT() paho.Client {
//	client, err := mqtt.NewConnection(mqtt.NewBrokerConfig())
//	if err != nil {
//		panic(err)
//	}
//	return client
//}
//
//// connect to the localdb postgresql
//func connectToDatabase() *database.Sqlite {
//	db, err := database.NewSqlite(viper.GetString(configkey.DatabaseRemoteDevFile), true)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}

func main() {
}
