package uuid

import (
	"landlords/helperv"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"sync"
	"time"
)

var _idwork *IdWorker

type IdWorker struct {
	startTime             int64
	workerIdBits          uint
	datacenterIdBits      uint
	maxWorkerId           int64
	maxDatacenterId       int64
	sequenceBits          uint
	workerIdLeftShift     uint
	datacenterIdLeftShift uint
	timestampLeftShift    uint
	sequenceMask          int64
	workerId              int64
	datacenterId          int64
	sequence              int64
	lastTimestamp         int64
	signMask              int64
	idLock                *sync.Mutex
}

func InitUUID() {
	_idwork = &IdWorker{}
	if err := _idwork.initIdWorker(100, 1); err != nil {
		return
	}
	glog.Info("初始化uuid成功")
}

func GetUUID() string {
	id, _ := _idwork.nextId()
	return conv.FormatInt64(id)
}

func (i *IdWorker) initIdWorker(workerId, datacenterId int64) error {
	var baseValue int64 = -1
	i.startTime = 1463834116272
	i.workerIdBits = 5
	i.datacenterIdBits = 5
	i.maxWorkerId = baseValue ^ (baseValue << i.workerIdBits)
	i.maxDatacenterId = baseValue ^ (baseValue << i.datacenterIdBits)
	i.sequenceBits = 12
	i.workerIdLeftShift = i.sequenceBits
	i.datacenterIdLeftShift = i.workerIdBits + i.workerIdLeftShift
	i.timestampLeftShift = i.datacenterIdBits + i.datacenterIdLeftShift
	i.sequenceMask = baseValue ^ (baseValue << i.sequenceBits)
	i.sequence = 0
	i.lastTimestamp = -1
	i.signMask = ^baseValue + 1

	i.idLock = &sync.Mutex{}

	if i.workerId < 0 || i.workerId > i.maxWorkerId {
		return errors.New(fmt.Sprintf("workerId[%v] is less than 0 or greater than maxWorkerId[%v].", workerId, datacenterId))
	}
	if i.datacenterId < 0 || i.datacenterId > i.maxDatacenterId {
		return errors.New(fmt.Sprintf("datacenterId[%d] is less than 0 or greater than maxDatacenterId[%d].", workerId, datacenterId))
	}
	i.workerId = workerId
	i.datacenterId = datacenterId
	return nil
}

func (i *IdWorker) nextId() (int64, error) {
	i.idLock.Lock()
	timestamp := time.Now().UnixNano()
	if timestamp < i.lastTimestamp {
		return -1, errors.New(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds", i.lastTimestamp-timestamp))
	}

	if timestamp == i.lastTimestamp {
		i.sequence = (i.sequence + 1) & i.sequenceMask
		if i.sequence == 0 {
			timestamp = i.tilNextMillis()
			i.sequence = 0
		}
	} else {
		i.sequence = 0
	}

	i.lastTimestamp = timestamp

	i.idLock.Unlock()

	id := ((timestamp - i.startTime) << i.timestampLeftShift) |
		(i.datacenterId << i.datacenterIdLeftShift) |
		(i.workerId << i.workerIdLeftShift) |
		i.sequence

	if id < 0 {
		id = -id
	}

	return id, nil
}

func (i *IdWorker) tilNextMillis() int64 {
	timestamp := time.Now().UnixNano()
	if timestamp <= i.lastTimestamp {
		timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	}
	return timestamp
}
