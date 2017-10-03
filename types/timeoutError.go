package types

type TimeoutError struct {
	Message string `json:"timeout_error,omitempty"`
}

func (err TimeoutError) Error() string {
	return err.Message
}
