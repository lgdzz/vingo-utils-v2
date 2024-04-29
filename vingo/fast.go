package vingo

func Diff[T any](old T, before func(T) T, after func(DiffBox)) {
	var box = DiffBox{Old: old}
	box.SetNewAndCompare(before(old))
	after(box)
}
