package main

import (
	"log"
	"msgpasser"
	"time"
)

var names = []string{"alice", "bob"}

//var names = []string {"alice", "bob"}
//var names = []string {"alice"}
var addrs = []string{"128.237.247.177", "128.2.247.16"}

//var addrs = []string {"unix12.andrew.cmu.edu", "unix13.andrew.cmu.edu"}
//var addrs = []string {"127.0.0.1"}
var ports = []int{9999}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Version: 1.0")

	var passer msgpasser.Passer
	passer.Init(9999)
	time.Sleep(3 * time.Second)

/*	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	conn, _ := net.DialUDP("udp", nil, addr)
	data := msgpasser.GameRoomData{"game room info", 0, nil, len(names), names[index]}
	data.Players = make([]msgpasser.PlayerInfo, len(names))
	for i := range names {
		data.Players[i].Ip = addrs[i]
		data.Players[i].Name = names[i]
	}
	b, err := json.Marshal(&data)
	log.Println(string(b))
	if err != nil {
		log.Println(err)
	}
	_, err = conn.Write(b)
	if err != nil {
		log.Println(err)
	}*/

	for {
		time.Sleep(30 * time.Second)
		log.Println("I am still alive!")
	}
}
