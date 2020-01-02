package es

// Error -
type Error string

func (e Error) Error() string { return string(e) }

const (
	// ErrMarshalEvent -
	ErrMarshalEvent = Error("unable to marshal event")

	// ErrUnmarshalEvent -
	ErrUnmarshalEvent = Error("unable to unmarshal event")

	// ErrUnknownEventType -
	ErrUnknownEventType = Error("unknown event type")
)
