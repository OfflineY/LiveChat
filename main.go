package main

import (
	"LiveChat/db"
	"LiveChat/service"
	"LiveChat/util"
	"LiveChat/ws"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var dbConn *service.RoomDetail

func main() {
	util.Paint("v3.0 BETA")
	fmt.Println("__init__")

	// 基础配置文件
	if util.IsIniExist("./conf.ini") {
		log.Println("Conf exists, skip initialization.")
	} else {
		log.Println("Conf does not exist, start initialize.")
	}

	cfg := util.IniLoad("./conf.ini")

	// 数据库连接
	databaseName := cfg.Section("MongoDB").Key("databaseName").Value()
	applyUrl := cfg.Section("MongoDB").Key("applyUrl").Value()
	maxPoolSize, err := cfg.Section("MongoDB").Key("maxPoolSize").Uint64()
	if err != nil {
		log.Println(err)
	}

	conn := db.Conn(
		applyUrl,
		maxPoolSize,
	)

	// TODO：通过数据库存储内容恢复房间

	// roomHubs 房间 hub 集合
	var roomHubs = make(map[string]*ws.Hub)

	dbConn = &service.RoomDetail{
		RoomHubs:     roomHubs,
		Conn:         conn,
		DatabaseName: databaseName,
	}

	for Room := range roomHubs {
		service.RecoverRoom(Room, dbConn)
	}

	// TODO
	service.CreateRoom("test_room_001", dbConn)

	fmt.Println("__main__")

	var wsAddr = flag.String(
		"wsAddr",
		":8080",
		"http.ListenAndServe",
	)

	var webAddr = flag.String(
		"webAddr",
		":5212",
		"http.ListenAndServe",
	)

	go service.RunWeb(*webAddr, dbConn)

	log.Println("Web port runs ->", *webAddr)

	log.Println("Websocket port runs ->", *wsAddr)

	err = http.ListenAndServe(*wsAddr, nil)
	if err != nil {
		log.Fatal("Error: ListenAndServe: ", err)
	}

}
