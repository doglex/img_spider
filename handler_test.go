package main

import "testing"

func TestWriteTxt(t *testing.T) {
	text := []string {"1","line2"}
	WriteTxt("FolderMadeByTest","1.txt",text,false)
}