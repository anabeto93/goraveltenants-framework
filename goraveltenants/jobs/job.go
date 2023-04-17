package jobs

type Job interface {
	Execute(args ...interface{}) error
}
