package foundation

type ExceptionRenderer interface {
	Render(throwable error) string
}
