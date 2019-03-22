package common

import (
	"fmt"
)

// Panicf prints out the message before pannicing with err
func Panicf(err error, format string, a ...interface{}) {
	fmt.Printf(format, a)
	panic(err)
}

// Panic prints out the message before pannicing with no message
func Panic(a ...interface{}) {
	fmt.Println(a)
	panic("")
}
