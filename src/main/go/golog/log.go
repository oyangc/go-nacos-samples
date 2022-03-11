package golog

import (
	"io"
	"os"
	"time"
)

const (
	//LOGPATH  LOGPATH/time.Now().Format(FORMAT)/*.log
	LOGPATH = "logs/"
	//FORMAT .
	YFORMAT = "2006"
	DFORMAT = "20060102"
	//LineFeed 换行
	LineFeed = "\r\n"
)

//以天为基准,存日志
var logPath = LOGPATH + time.Now().Format(DFORMAT) + "/"

//WriteLog return error
func WriteLog(fileName, msg string) error {
	if !IsExist(logPath) {
		CreateDir(logPath)
	}
	var (
		err error
		f   *os.File
	)
	f, err = os.OpenFile(logPath+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	fd_time := time.Now().Format("2006-01-02 15:04:05")
	_, err = io.WriteString(f, LineFeed+fd_time+" "+msg)

	defer f.Close()
	return err
}

//CreateDir  文件夹创建
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

//IsExist  判断文件夹/文件是否存在  存在返回 true
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
