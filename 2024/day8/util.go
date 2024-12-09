package main

type Numbered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func Abs[T Numbered](n T) T {
	if n < 0 {
		return n - n - n
	}
	return n
}
