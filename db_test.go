package main

import (
	"fmt"
	"testing"
	"time"
)

func TestHistoryContains(t *testing.T)  {
	x1 := HistoryContains("6")
	x2 := HistoryContains("62")
	fmt.Println(x1, x2)
	go HistoryCheckPointInterval()
	time.Sleep(200*time.Second)
}