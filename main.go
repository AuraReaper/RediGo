package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/AuraReaper/redigo/config"
	"github.com/AuraReaper/redigo/server"
)

func main() {
	setupFlags()
	log.Println("rolling the RediGo")

	var sigs chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	var wg sync.WaitGroup
	wg.Add(2)

	go server.RunAsyncTCPServer(&wg)
	go server.WaitForSignal(&wg, sigs)
}

func setupFlags() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the redigo server")
	flag.IntVar(&config.Port, "port", 7379, "port for the redigo server")
	flag.Parse()
}
