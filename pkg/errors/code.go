package errors

// normal code of application.
const (
	InvalidCode int = iota + 1000001
	NotFoundCode
	CannotExecCode
	InternalCode
)

// hertz code.
const (
	RouteNotFoundCode int = iota + 1100001
	BindErrorCode
)
