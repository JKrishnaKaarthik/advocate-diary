package app

type StatusError struct {
	Status int
	Msg    string
}

func (e *StatusError) Error() string { return e.Msg }

func Bad(msg string) error      { return &StatusError{Status: 400, Msg: msg} }
func NotFound() error           { return &StatusError{Status: 404, Msg: "Not found"} }
func Conflict(msg string) error { return &StatusError{Status: 409, Msg: msg} }
