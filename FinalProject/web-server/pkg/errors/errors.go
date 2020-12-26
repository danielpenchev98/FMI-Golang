package errors

import "github.com/pkg/errors"

//ClientError represents a problem with client request
type ClientError struct {
	Err error
}

func (e *ClientError) Error() string {
	return e.Err.Error()
}

func NewClientError(description string) *ClientError {
	return &ClientError{
		Err: errors.New(description),
	}
}

type ItemNotFoundError struct {
	Err error
}

func (e *ItemNotFoundError) Error() string {
	return e.Err.Error()
}

func NewItemNotFoundError(description string) *ItemNotFoundError {
	return &ItemNotFoundError{
		Err: errors.New(description),
	}
}

//ServerError represents a problem with server
type ServerError struct {
	Err error
}

func (e *ServerError) Error() string {
	return e.Err.Error()
}

func NewServerError(description string, err error) *ServerError {
	return &ServerError{
		Err: errors.Wrapf(err, description),
	}
}
