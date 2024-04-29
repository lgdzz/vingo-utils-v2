package vingo

import (
	"github.com/lgdzz/vingo-utils/vingo"
)

func Diff[T any](old T, before func(T) T, after func(vingo.DiffBox)) {
	var box = vingo.DiffBox{Old: old}
	box.SetNewAndCompare(before(old))
	after(box)
}
