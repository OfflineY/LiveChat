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

func main() {
	util.Paint("v3.0 BETA")

	fmt.Println("__init__")

	if util.IsIniExist("./conf.ini") {
		log.Println("Conf exists, skip initialization.")
	} else {
		log.Println("Conf does not exist, start initialize.")
		util.IniInitial("./conf.ini")
	}

	cfg := util.IniLoad("./conf.ini")

	databaseName := cfg.Section("MongoDB").Key("databaseName").Value()
	applyUrl := cfg.Section("MongoDB").Key("applyUrl").Value()
	maxPoolSize, err := cfg.Section("MongoDB").Key("maxPoolSize").Uint64()
	if err != nil {
		log.Println(err)
	}

	conn := db.Conn(applyUrl, maxPoolSize)

	// groups hub
	var roomHubs = make(map[string]*ws.Hub)

	dbConn := &service.DatabaseConn{
		Conn: conn,
		Name: databaseName,
	}

	rooms := &service.Rooms{
		RoomHubs: roomHubs,
	}

	for Room := range roomHubs {
		rooms.RecoverRoom(Room, dbConn)
	}

	rooms.CreateRoom("test_room_001", dbConn)

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

	go service.RunWeb(*webAddr, dbConn, rooms)

	log.Println("Web port runs ->", *webAddr)

	log.Println("Websocket port runs ->", *wsAddr)

	err = http.ListenAndServe(*wsAddr, nil)
	if err != nil {
		log.Fatal("Error: ListenAndServe: ", err)
	}

}
