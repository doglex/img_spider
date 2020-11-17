package main

import (
	"fmt"
	"testing"
)

func TestTryRecover(t *testing.T) {
	defer func () {
		defer DeferRecover()
		panic("something wrong")
	}()
}

func TestLegalFileName(t *testing.T) {
	fileName := "67//werx//https"
	fmt.Println(LegalPathName(fileName))
}

