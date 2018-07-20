package errors

type CmdError struct {
	ExitCode int
	Err      error
}

func (e CmdError) Error() string {
	return e.Err.Error()
}
