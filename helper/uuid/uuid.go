package uuid

import (
	"github.com/google/uuid"
	"strings"
)

func GetUUID() string {
	result := strings.ReplaceAll(uuid.New().String(), "-", "")
	return result[8:24]
}
