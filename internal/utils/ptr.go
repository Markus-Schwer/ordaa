package ptr

func To[T any](in T) *T {
	return &in
}
