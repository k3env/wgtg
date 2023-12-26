package util

import (
	"slices"
)

func Find[T interface{}](slice []*T, comparator func(*T) bool) *T {
	index := slices.IndexFunc(slice, comparator)
	if index == -1 {
		return nil
	}
	return slice[index]
}

func Has[T interface{}](slice []*T, fn func(*T) bool) bool {
	index := slices.IndexFunc(slice, fn)
	return index != -1
}
