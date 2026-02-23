package core

type RedigoCmd struct {
	Cmd  string
	Args []string
}

type RedigoCmds []*RedigoCmd
