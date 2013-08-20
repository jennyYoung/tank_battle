package main

import (
	"encoding/json"
	"fmt"
	"log"
	"msgpasser"
	"net"
	"time"
)

var names = []string{"alice", "bob", "charlie"}
var localName = "charlie"

//var names = []string {"alice", "bob"}
//var names = []string {"alice"}
var addrs = []string{"unix11.andrew.cmu.edu", "unix12.andrew.cmu.edu", "unix13.andrew.cmu.edu"}

//var addrs = []string {"128.237.230.19", "unix11.andrew.cmu.edu"}
//var addrs = []string {"127.0.0.1"}
var ports = []int{9999}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Version: 1.0")

	var passer msgpasser.Passer
	passer.Init(9999)
	time.Sleep(3 * time.Second)

	for i, addrS := range addrs {
		addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:9999", addrS))
		conn, _ := net.DialUDP("udp", nil, addr)
		data := msgpasser.GameRoomData{"game room info", 0, nil, len(names), names[i], "join", localName}
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
		}
	}

	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	conn, _ := net.DialUDP("udp", nil, addr)

	// To avoid too many input happen
	time.Sleep(3 * time.Second)

	x := 10000
	y := 3
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			msg := msgpasser.Data{fmt.Sprintf("%v %d", localName, i*y+j), true}
			b, _ := json.Marshal(&msg)
			conn.Write(b)
		}
		time.Sleep(300 * time.Millisecond)
	}

	for {
		time.Sleep(30 * time.Second)
		log.Println("I am still alive!")
	}
}
