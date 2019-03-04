package mod1

import "dad/types"

type virus struct {
	age int
}

var A types.Module = virus{age: 5}

func (a virus) Init() chan string { return make(chan string) }
func (a virus) Shutdown()         {}

func (a virus) Id() string { return "adam" }
