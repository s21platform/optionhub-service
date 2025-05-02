package utils

func TransformToPtr[T any](data T) *T {
	res := new(T)
	*res = data
	return res
}
