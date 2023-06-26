package internal

import "math/rand"

func CopySlice[T any](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	return result
}

func RandomElement[T any](slice []T) T {
	return slice[rand.Intn(len(slice))]
}
