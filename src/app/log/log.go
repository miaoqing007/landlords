package log

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
)

const (
	logDir         = "serverLog"
	MaxLogFileSize = 1024 * 1024 * 5
)

func InitLog() {
	if !CreateDir(logDir) {
		fmt.Println("CreateDir() failed, logDir = ", logDir)
		return
	}
	flag.Set("alsologtostderr", "true") // 日志写入文件的同时，输出到stderr
	flag.Set("log_dir", logDir)         // 日志文件保存目录
	//flag.Set("v", conv.FormatInt(Unc))  // 配置V输出的等级。
	flag.Parse()
	glog.MaxSize = MaxLogFileSize
}
