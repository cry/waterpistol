package types

type Module interface {
	Id() string
	Init() chan string
	Shutdown()
}
