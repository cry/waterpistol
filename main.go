package main

import (
	"malware/mod1"
	"malware/mod2"
	"malware/types"
	"fmt"
)

var modules = []types.Module{mod1.A, mod2.A}

func main() {

	for _, val := range modules {
		fmt.Println(val.Id())
	}

}
