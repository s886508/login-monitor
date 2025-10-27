package input

type Input interface {
	Init() error
	FetchEvent() (string, error)
	Close()
}
