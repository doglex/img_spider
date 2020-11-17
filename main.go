package main

import (
	"github.com/patrickmn/go-cache"
	"log"
)

const DbHistoryStorage = "history.db" // 断点续传文件
const NumWorker = 32                  // 下载协程数量
var DbHistory = cache.New(cache.NoExpiration, cache.NoExpiration)
var CountSuccess int64 = 0                  // 已成功数量
var CountFail int64 = 0                     // 已失败数量
var BusDownload = make(chan DownLoad, 1024) // 任务队列
var DirDownImg = ""

func init() {
	if FileExists(DbHistoryStorage) {
		DbHistory.LoadFile(DbHistoryStorage)
	}
	DirDownImg = "IMAGE"
	log.Println("Starting APP: Image path:", DirDownImg)
}

func main() {
	go HistoryCheckPointInterval()
	go PrintCounterInterval()
	for i := 0; i < NumWorker; i++ {
		go DownloadFileWorker(BusDownload)
	}
	//RunListSexy()
	//RunListCN()
	//RunListCompany()
	RunListSets()
}
