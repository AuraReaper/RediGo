package server

import (
	"log"
	"net"
	"syscall"
	"time"

	"github.com/AuraReaper/redigo/config"
	"github.com/AuraReaper/redigo/core"
)

var conClients int = 0
var cronFrequency time.Duration = 1 * time.Second
var lastCronExecTime time.Time = time.Now()

func RunAsyncTCPServer() error {
	log.Println("starting an async TCP server on", config.Host, config.Port)

	maxClients := 20000

	// Create EPOLL Event Objects to hold events
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, maxClients)

	// Create a socket
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFD)

	// Set the Socket operate in a non-blocking mode
	if err := syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	// Bind the IP and port
	ip4 := net.ParseIP(config.Host)
	if err := syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		return err
	}

	// Start Listening
	if err := syscall.Listen(serverFD, maxClients); err != nil {
		return err
	}

	// creating EPOLL instance
	epollFD, err := syscall.EpollCreate(0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(epollFD)

	var sockertServerEvent syscall.EpollEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFD),
	}

	// Listen to read events on the Server itself
	if err := syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &sockertServerEvent); err != nil {
		return err
	}

	for {
		if time.Now().After(lastCronExecTime.Add(cronFrequency)) {
			core.DeleteExpiredKeys()
			lastCronExecTime = time.Now()
		}

		nevents, e := syscall.EpollWait(epollFD, events[:], -1)
		if e != nil {
			continue
		}

		for i := 0; i < nevents; i++ {
			// if the socket server itself is ready for an IO
			if int(events[i].Fd) == serverFD {
				// accept the incoming connection from a client
				fd, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("err", err)
					continue
				}

				// increase the number of concurrent clients count
				conClients++
				syscall.SetNonblock(serverFD, true)

				// add this new TCP connection to be monitored
				var socketClientEvent syscall.EpollEvent = syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(fd),
				}

				if err := syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketClientEvent); err != nil {
					log.Fatal(err)
				}
			} else {
				comm := core.FDComm{Fd: int(events[i].Fd)}
				cmds, err := readCommand(comm)
				if err != nil {
					syscall.Close(int(events[i].Fd))
					conClients--
					continue
				}

				respond(cmds, comm)
			}
		}
	}
}
