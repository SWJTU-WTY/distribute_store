package util

type MyError struct {
	Msg string
}

func Construct(msg string) MyError {
	return MyError{msg}
}

func (m MyError) Error() string {
	return m.Msg
}
