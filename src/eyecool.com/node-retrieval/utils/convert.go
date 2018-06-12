package utils

import "strconv"

func Itof32(i int) float32 {
	maxstr := strconv.Itoa(i)
	f64, err := strconv.ParseFloat(maxstr, 32)
	if err !=nil{
		return 0;
	}
	return float32(f64)
}
