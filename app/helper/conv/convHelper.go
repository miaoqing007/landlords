package conv

import (
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"strconv"
)

func ParseInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		glog.Error(errors.WithStack(err))
		return 0
	}
	return v
}

func ParseUint(s string) uint {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		glog.Error(errors.WithStack(err))
		return 0
	}
	return uint(v)
}

func ParseUint64(s string) uint64 {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		glog.Error(errors.WithStack(err))
		return 0
	}
	return v
}

func ParseInt64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		glog.Error(errors.WithStack(err))
		return 0
	}
	return v
}

func ParseInt32(s string) int32 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		glog.Error(errors.WithStack(err))
		return 0
	}
	return int32(v)
}

func ParseUint32(s string) uint32 {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		glog.Error(errors.WithStack(err))
		return 0
	}
	return uint32(v)
}

func ParseInt16(s string) int16 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		glog.Error(errors.WithStack(err))
		return 0
	}
	return int16(v)
}

func ParseUint16(s string) uint16 {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		glog.Error(errors.WithStack(err))
		return 0
	}
	return uint16(v)
}

func FormatInt(i int) string {
	return FormatInt64(int64(i))
}

func FormatInt32(i int32) string {
	return FormatInt64(int64(i))
}

func FormatInt64(i int64) string {
	return strconv.FormatInt(int64(i), 10)
}

func FormatUint(i uint) string {
	return FormatUint64(uint64(i))
}

func FormatUint32(i uint32) string {
	return FormatUint64(uint64(i))
}

func FormatUint64(i uint64) string {
	return strconv.FormatUint(uint64(i), 10)
}

func FormatInt16(i int16) string {
	return FormatInt64(int64(i))
}

func FormatUint16(i uint16) string {
	return FormatInt64(int64(i))
}
func FormatUint8(i uint8) string {
	return FormatInt64(int64(i))
}
