package qqai

type Error struct {
	Ret int
	Msg string
}

func (e *Error) Error() string {
	return e.Msg
}

func (e *Error) ErrorCode() int {
	return e.Ret
}
