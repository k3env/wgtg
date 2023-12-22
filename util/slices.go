package util

import (
	"errors"
	"slices"
)

func Find[T interface{}](slice []*T, comparator func(*T) bool) (*T, error) {
	index := slices.IndexFunc(slice, comparator)
	if index == -1 {
		return nil, errors.New("item not found")
	}
	return slice[index], nil
}
