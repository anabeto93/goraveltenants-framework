package contracts

type Listener interface {
	Handle(args ...interface{}) error
	Name() string
}
