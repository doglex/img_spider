package main

import (
	"log"
	"os"
	"regexp"
	"runtime"
)

// 可以被嵌入到Defer中的Recover
func DeferRecover (){
	if err := recover(); err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Println("Recover", file, line, err)
	}
}

// 地址转换为合法文件名
func LegalPathName(uri string) (fileName string) {
	pattern :=regexp.MustCompile("[^a-z0-9A-Z.]+")
	fileName = pattern.ReplaceAllString(uri,"_")
	return
}

// 判断文件是否存在
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
