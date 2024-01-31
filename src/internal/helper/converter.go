package helper

func ValToSlice[T any](input *T) []T {
	if input != nil {
		return []T{*input}
	}
	return []T{}
}
