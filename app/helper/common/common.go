package common

import "sync"

func GetSyncMapLen(m sync.Map) int {
	length := 0
	m.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}
