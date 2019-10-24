package log

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

const (
	MacOS   string = "darwin"
	Windows string = "windows"
	Linux   string = "linux"
)

func CreateDir(name string) bool {
	if IsDir(name) {
		fmt.Printf("%s is already a directory.\n", name)
		return true
	}

	if createDirImpl(name) {
		fmt.Println("Create directory successfully.")
		return true
	}

	return false
}

// 防止Unix系统出现权限问题，logs统一可执行
func IsDir(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		fmt.Println("Error: ", err)
		return false
	}

	curSystem := getExeDirAccording2System()
	if curSystem == Windows {
		return fi.IsDir()
	}

	cmd := exec.Command("chmod", "-R", "777", logDir)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return fi.IsDir()
}

//判断当前系统
func getExeDirAccording2System() string {
	curSystem := runtime.GOOS
	return curSystem
}

// windows 0666权限创建文件夹，所有者具备rwx权限
// MacOS 或Unix permission 0777创建，所有者才会具备rwx权限，否则只具备rw-权限
// 无法操作logs文件夹，日志写入文件失败
func createDirImpl(name string) bool {
	var err error

	//right
	curSystem := getExeDirAccording2System()
	if curSystem == Windows {
		err = os.Mkdir(name, 0666)
	} else {
		err = os.Mkdir(name, os.ModePerm)
	}
	if err == nil {
		return true
	} else {
		fmt.Println("Error: ", err)
		return false
	}
}
