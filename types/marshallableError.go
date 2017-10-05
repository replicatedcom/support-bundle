package types

type MarshallableError struct {
	Message string `json:"error,omitempty"`
}

func (err MarshallableError) Error() string {
	return err.Message
}
