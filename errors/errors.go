package errors

import (
	"errors"
	"strings"
)

var IPNotFoundError = NewError("IP address for interface not found")
var InvalidIPError = NewError("Invalid IP address")
var NoAllocableIPsError = NewError("IP pool is full")
var NoWGInterfaces = NewError("Wireguard interfaces not found!")

type TypeError struct {
	err  error
	text string
}

func NewError(text string) TypeError {
	s := TypeError{
		err:  errors.New(strings.ToLower(text)),
		text: text,
	}
	return s
}
func (e TypeError) String() string {
	return e.text
}
func (e TypeError) Error() string {
	return e.text
}
