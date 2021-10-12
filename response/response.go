package response

type Response interface {
	Data() interface{}
	Error() error
	Status() string
	HTTPStatusCode() int
	Message() string
	Meta() interface{}
}
