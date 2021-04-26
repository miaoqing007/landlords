package util

import (
	"strconv"
)

func SubString(s string, start, end int) string {
	souce := []rune(s)
	n := len(souce)
	if start < 0 || end > n || start > end {
		return ""
	}
	return string(souce[start:end])
}

func SortArrayStringBig2Small(array []string) []string {
	for i := 0; i < len(array)-1; i++ {
		for j := i; j < len(array); j++ {
			vi, _ := strconv.Atoi(array[i][1:])
			vj, _ := strconv.Atoi(array[j][1:])
			if vi < vj { //从大到小  从小到大 a[i]>a[j]
				array[i], array[j] = array[j], array[i]
			}
		}
	}
	return array
}

func SortArrayStringSmall2Big(array []string) []string {
	for i := 0; i < len(array)-1; i++ {
		for j := i; j < len(array); j++ {
			vi, _ := strconv.Atoi(array[i][1:])
			vj, _ := strconv.Atoi(array[j][1:])
			if vi > vj { //从大到小  从小到大 a[i]>a[j]
				array[i], array[j] = array[j], array[i]
			}
		}
	}
	return array
}
