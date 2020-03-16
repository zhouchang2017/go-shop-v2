package utils

func Max(a ...int64) (max int64) {
	max = a[0]
	for _, value := range a {
		if value > max {
			max = value
		}
	}
	return max
}

func Min(a ...int64) (min int64) {
	min = a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
	}
	return min
}
