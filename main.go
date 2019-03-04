package main

import (
	"dad/mod1"
	"dad/mod2"
	"dad/types"
	"fmt"
)

var modules = []types.Module{mod1.A, mod2.A}

func main() {

	for _, val := range modules {
		fmt.Println(val.Id())
	}

}
