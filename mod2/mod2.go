package mod2

import "dad/types"

type module struct {
	age int
}

var A types.Module = module{age: 5}

func (a module) Init() chan string { return make(chan string) }
func (a module) Shutdown()         {}

func (a module) Id() string { return "adam2" }
