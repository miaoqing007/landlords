package idgenerator

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	sf             *Sonyflake
	machineType    uint16
	machineTypeMap = map[string]uint16{
		"login":   1,
		"gateway": 2,
		"online":  3,
		"home":    4,
		"game":    5,
		"gm":      6,
		"upload":  7,
		"robot":   8,
		"unknow":  9,
	}
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {

	processName := os.Args[0]
	for key, value := range machineTypeMap {
		if strings.Contains(processName, key) {
			machineType = value
			break
		}
	}
	if machineType == 0 {
		machineType = machineTypeMap["unknow"]
	}

	var st Settings
	st.StartTime = time.Date(2006, 01, 02, 15, 04, 05, 0, time.UTC)
	st.MachineID = createMachineID

	sf = NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func createMachineID() (uint16, error) {
	if machineType == 0 {
		return 0, errors.New("机器类型没有设置")
	}
	return machineType, nil
}

// FetchID 生成唯一ID
func FetchID() uint64 {
	// sf改造后 按照实例类型类型区分ID
	id, err := sf.NextID()
	if err != nil {
		log.Printf("[id] fetch id error:", err)
		return 0
	}
	return id
}

func ip2Int(ip string) uint32 {
	var ipInt uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &ipInt)
	return ipInt
}

func backtoIP4(ipInt int64) string {

	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipInt & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// FetchUUID ...
func FetchUUID() string {
	id, err := NewV4()
	if err != nil {
		logger.Errorf("[uuid] fetch id:%v error:%v", id, err)
		return randStringBytes(16)
	}
	return id.String()
}

// FetchServerID 生成服务器ID(根据服名称 服进程索引 以及IP机器ID)
func FetchServerID(srvName string, idx int) string {
	ip, err := privateIPv4()
	if err != nil {
		return FetchUUID()
	}
	machineID := ip2Int(ip.String())
	return srvName + "-" + strconv.Itoa(int(machineID)) + "-" + strconv.Itoa(idx)
}

// FetchServiceID 生成服务ID(根据服务名称 服进程索引 以及IP机器ID)
func FetchServiceID(svcName string, svcAddr string) string {
	host, port, err := net.SplitHostPort(svcAddr)
	if err != nil {
		return svcName + "-" + svcAddr
	}
	ipInt := ip2Int(host)
	return fmt.Sprintf("%v-%v", ipInt, port)
}

// GetTimeFromID 通过sonyflake生成的id获取生成时间
func GetTimeFromID(id uint64) time.Time {
	parts := Decompose(id)
	// 过去的时间的纳秒数
	overtime := (sf.startTime + int64(parts["time"])) * sonyflakeTimeUnit
	return time.Unix(0, overtime)
}
