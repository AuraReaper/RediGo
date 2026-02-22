package main

import (
	"flag"
	"log"

	"github.com/AuraReaper/redigo/config"
	"github.com/AuraReaper/redigo/server"
)

func main() {
	setupFlags()
	log.Println("rolling the RediGo")
	server.RunSyncTCPServer()
}

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the redigo server")
	flag.IntVar(&config.Port, "port", 7379, "port for the redigo server")
	flag.Parse()
}
