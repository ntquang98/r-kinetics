package pointer

func GetValueIfNotNil[T any](v *T) T {
	var rs T
	if v != nil {
		rs = *v
	}
	return rs
}

func GetPointer[T any](v T) *T {
	return &v
}
