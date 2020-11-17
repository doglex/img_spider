package main

import (
	"fmt"
	"time"
)



func init() {
	if FileExists(DbHistoryStorage) {
		DbHistory.LoadFile(DbHistoryStorage)
	}
}

func HistoryContains(key string) (ok bool) {
	defer  DeferRecover()
	ok = false
	_, ok = DbHistory.Get(key)
	return
}

func HistoryAdd(key string) {
	DbHistory.SetDefault(key, "1")
}

func HistoryCheckPointInterval() {
	for {
		time.Sleep(30*time.Second)
		err := DbHistory.SaveFile(DbHistoryStorage)
		if err != nil {
			fmt.Println("Warn",err)
			recover()
		}
	}
}