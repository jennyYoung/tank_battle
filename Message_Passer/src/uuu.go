package main

import (
	"log"
	"msgpasser"
	"os"
	"os/exec"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Version: 1.0")

	var passer msgpasser.Passer
	passer.Init(9999)
	time.Sleep(3 * time.Second)

	cmd := exec.Command("python", "bSserver.py", os.Args[1], os.Args[2])
	go cmd.Run()
	/*if err != nil {
		log.Println(err)
	}*/

	time.Sleep(3 * time.Second)

	cmd2 := exec.Command("java", "-jar", "ttt.jar", os.Args[1], os.Args[2])
	go cmd2.Run()
	/*if err != nil {
		log.Println(err)
	}*/

	for {
		time.Sleep(30 * time.Second)
		log.Println("I am still alive!")
	}
}
