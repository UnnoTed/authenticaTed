package errors

// ID is a custom type for errors ids
type ID int

// Error is similar to the default error struct from golang
// but this one specified a ID
type Error struct {
	ID      ID
	Message string
}

// FromError builds a "Error" from native golang error
func FromError(err error) *Error {
	e := new(Error)
	e.Message = err.Error()
	return e
}

// New returns a error compactible with the default Go "error"
func New(msg string) *Error {
	e := new(Error)
	e.Message = msg
	return e
}

// New creates a new Error
func NewFromID(id ID) *Error {
	e := new(Error)
	e.ID = id
	e.Message = Messages[id]

	return e
}

// Error gives support to native golang error
// it's used for testing
func (e *Error) Error() string {
	return e.Message
}

// String returns the error message
func (e *Error) String() string {
	return e.Message
}

// JSON converts the error to json
// it's used on APIs
func (e *Error) JSON() map[string]interface{} {
	i := map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"id":      e.ID,
			"message": e.Message,
		},
	}

	return i
}
