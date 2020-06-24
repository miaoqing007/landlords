package log

import (
	"flag"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"time"
)

const (
	MaxLogFileSize = 1024 * 1024 * 5
	FileSaveTime   = 12 * time.Hour
)

var logDir = "serverLog"

func InitLog() {
	if logEnv := os.Getenv("LogDir"); logEnv != "" {
		logDir = logEnv
	}
	if !CreateDir(logDir) {
		return
	}
	flag.Set("alsologtostderr", "true") // 日志写入文件的同时，输出到stderr
	flag.Set("log_dir", logDir)         // 日志文件保存目录
	//flag.Set("v", conv.FormatInt(Unc))  // 配置V输出的等级。
	flag.Parse()
	glog.MaxSize = MaxLogFileSize

	go findAndRemoveOutTimeLogFile()
	glog.Info("初始化日志完成")
}

func findAndRemoveOutTimeLogFile() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			if err := removeOutTimeLogFile(); err != nil {
				glog.Errorf("clean outTime LogFile Failed err(%v)", err)
			}
		}
	}
}

func removeOutTimeLogFile() error {
	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if time.Since(file.ModTime()).Seconds() > time.Duration(FileSaveTime).Seconds() {
			if err := os.Remove(logDir + "\\" + file.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}
