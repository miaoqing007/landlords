package client_handle

import (
	"fmt"
	"math/rand"
	"time"
)

func P_Shuff_req() {
	arr := shuffle()
	fmt.Println(arr)
}

func shuffle() []int {
	rand.Seed(time.Now().UnixNano())
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28,
		29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42,
		43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54}
	var i, j int
	var temp int
	for i = len(arr) - 1; i > 0; i-- {
		j = rand.Intn(len(arr))
		if arr[i] == arr[j] {
			continue
		}
		temp = arr[i]
		arr[i] = arr[j]
		arr[j] = temp
	}
	return arr
}
