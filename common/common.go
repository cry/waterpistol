package common

import "fmt"

// Panicf prints out the message before pannicing with err
func Panicf(err error, format string, a ...interface{}) {
	fmt.Printf(format, a)
	panic(err)
}
