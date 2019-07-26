package main

import (
	"flag"
	"fmt"
	"log"

	kpb "./kunpengBattle"
)

var ip string
var port int
var teamID int

func init() {
	flag.IntVar(&teamID, "teamID", 0, "TeamID")
	flag.StringVar(&ip, "ip", "127.0.0.1", "ServerIP")
	flag.IntVar(&port, "port", 6001, "Server Port")
}

func main() {
	flag.Parse()
	fmt.Println("testing:", ip, port, teamID)

	client := kpb.NewKunPengBattleClient(teamID, "Hadrianl")
	err := client.Connect(ip, port)
	if err != nil {
		log.Panicln("Connection Failed!!!")
	}

	client.Start()

}
