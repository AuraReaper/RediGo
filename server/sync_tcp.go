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

func toArrayString(ai []interface{}) ([]string, error) {
	as := make([]string, len(ai))
	for i := range ai {
		as[i] = ai[i].(string)
	}
	return as, nil
}

func readCommands(c io.ReadWriter) (core.RedigoCmds, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return nil, err
	}

	raw, err := core.Decode(buf[:n])
	if err != nil {
		return nil, err
	}

	values, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", raw)
	}

	var cmds []*core.RedigoCmd = make([]*core.RedigoCmd, 0)
	for _, value := range values {
		tokens, err := toArrayString(value.([]interface{}))
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, &core.RedigoCmd{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		})
	}
	return cmds, nil
}

func respond(cmds core.RedigoCmds, c io.ReadWriter) {
	core.EvalAndRespond(cmds, c)
}

func RunSyncTCPServer() {
	log.Println("starting a synchronous TCP server on", config.Host, config.Port)

	var con_clients int = 0

	// listening to the configured host:port
	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		log.Println("err", err)
		return
	}

	for {
		// blocking call: waiting for the new client to connect
		c, err := lsnr.Accept()
		if err != nil {
			log.Println("err", err)
		}

		// increment the number of concurrent clients
		con_clients += 1

		for {
			// over the socket, continuously read the command and print it out
			cmds, err := readCommands(c)
			if err != nil {
				c.Close()
				con_clients -= 1
				if err == io.EOF {
					break
				}
			}
			respond(cmds, c)
		}
	}
}
