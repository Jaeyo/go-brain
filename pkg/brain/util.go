package brain

type numbers interface {
	int | int64
}

func maxKey[T1 numbers, T2 any](m map[T1]T2) T1 {
	var max T1
	for key := range m {
		if key > max {
			max = key
		}
	}
	return max
}
