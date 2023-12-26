package errors

import (
	"errors"
	"strings"
)

var IPNotFoundError = NewError("IP address for interface not found")
var InvalidIPError = NewError("Invalid IP address")
var NoAllocableIPsError = NewError("IP pool is full")
var NoWGInterfaces = NewError("Wireguard interfaces not found!")

type errorStruct struct {
	err  error
	text string
}

func NewError(text string) errorStruct {
	s := errorStruct{
		err:  errors.New(strings.ToLower(text)),
		text: text,
	}
	return s
}
func (e errorStruct) String() string {
	return e.text
}
func (e errorStruct) Error() string {
	return e.text
}
