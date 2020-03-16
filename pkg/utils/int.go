package utils

import "math"

func Int64AvgN(price int64, n int64) []int64 {
	res := make([]int64, n)
	for n > 0 {
		p := price / n
		if n > 1 {
			p = int64(math.Ceil(float64(p/10)) * 10)
		}
		res[n-1] = p
		price -= p
		n -= 1
	}
	return res
}