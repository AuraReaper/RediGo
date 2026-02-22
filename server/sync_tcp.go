package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/AuraReaper/redigo/config"
	"github.com/AuraReaper/redigo/core"
)

func RunSyncTCPServer() {
	log.Println("rolling a synchronous TCP server on", config.Host, config.Port)

	var conClients int = 0 // count of concurrent clients

	// listening on the configures host:port
	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {
		// blocking call: waiting for the new client to connect
		c, err := lsnr.Accept()
		if err != nil {
			panic(err)
		}

		// increment the no of concurrent clients
		conClients += 1
		log.Println("client connected with address:", c.RemoteAddr(), "concurrent clients", conClients)

		for {
			// over the socket, continously read the command and print it out
			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				conClients -= 1
				log.Println("client disconnected", c.RemoteAddr(), "concurrent clients", conClients)
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}
			respond(cmd, c)
		}
	}
}

func readCommand(c net.Conn) (*core.RedigoCmd, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return nil, err
	}

	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}

	return &core.RedigoCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respond(cmd *core.RedigoCmd, c net.Conn) {
	err := core.EvalAndRespond(cmd, c)
	if err != nil {
		respondError(err, c)
	}
}

func respondError(err error, c net.Conn) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}
