package internal

func CopySlice[T any](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)
	return result
}
