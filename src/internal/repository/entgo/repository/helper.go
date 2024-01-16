package repository

func toUint64s[T any](list []*T, fn func(item *T) uint64) []uint64 {
	var id []uint64
	for _, item := range list {
		id = append(id, fn(item))
	}
	return id
}
