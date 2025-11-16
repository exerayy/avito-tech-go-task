package helper

func RemoveElement[T comparable](s []T, element T) []T {
	for i, el := range s {
		if el == element {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
