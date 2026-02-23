package server

import (
	"fmt"
	"io"
	"strings"

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

func respond(cmds core.RedigoCmds, c *core.Client) {
	core.EvalAndRespond(cmds, c)
}
