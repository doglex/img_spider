package main

import (
	"bufio"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/levigross/grequests"
	"log"
	"os"
	"path"
	"strings"
	"sync/atomic"
	"time"
)

func GetUA() string {
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36 Edg/84.0.522.63"
}

// s是分类
// t是明星
// a是某个专辑
// x是某个全集


type DownLoad struct {
	folder  string
	url     string
	referer string
}


func PrintCounterInterval() {
	for {
		fmt.Println(time.Now(),"Counter=====================", "\t:success:",CountSuccess, "\t:fail:",CountFail)
		time.Sleep(10 * time.Second)
	}
}


// 遍历一个专辑
func WalkGallery(url string) {
	defer DeferRecover()
	if HistoryContains(url) {
		fmt.Println("<skip> In database: ", url)
		return
	} else {
		time.Sleep(time.Second)
		HistoryAdd(url)
		fmt.Println("Walking ", url)
	}
	ro := grequests.RequestOptions{
		RequestTimeout: time.Second * 60,
		Headers: map[string]string{
			"Referer":    "https://www.tujigu.com/",
			"User-Agent": GetUA(),
		}}
	resp, err := grequests.Get(url, &ro)
	if err != nil {
		fmt.Println("Unable to make request, enter sleep: ", err)
		time.Sleep(60*time.Second)
		return
	}
	respString := resp.String()
	doc := soup.HTMLParse(respString)
	// 取到标题和文件夹
	folder, e := ExtractGalleryMeta(doc, url)
	if e {
		return
	}
	// 下载图片
	links := doc.Find("div", "class", "content").FindAll("img")
	for _, link := range links {
		link := link
		BusDownload <- DownLoad{folder,link.Attrs()["src"],url}
		//go DownloadFile(folder, link.Attrs()["src"], url)
	}
	// 找到其他链接
	sisterPages := doc.Find("div", "id", "pages").FindAll("a")
	for _, link := range sisterPages {
		WalkGallery(link.Attrs()["href"])
	}
	time.Sleep(60 * time.Second) // 保持可用
}

func DownloadFileWorker(bus chan DownLoad) {
	for {
		down := <- bus
		DownloadFile(down.folder,down.url,down.referer)
	}
}

// 下载单文件，一般是图片文件
func DownloadFile(folder string, url string, referer string) {
	defer DeferRecover()
	_ = os.MkdirAll(folder, os.ModePerm)
	fileName := LegalPathName(url)
	fullFileName := path.Join(folder, fileName)
	// 不要启用，有些图片没有下载完整
	//if FileExists(fullFileName) {
	//	return
	//}
	ro := grequests.RequestOptions{
		RequestTimeout: time.Second * 90,
		Headers: map[string]string{
			"Referer":    referer,
			"User-Agent": GetUA(),
		}}
	resp, err := grequests.Get(url, &ro)
	if err != nil {
		log.Println("Error", "get file fail:", url)
		atomic.AddInt64(&CountFail,1)
		return
	}
	if err := resp.DownloadToFile(fullFileName); err != nil {
		log.Println("Error", "save file fail:", err)
		atomic.AddInt64(&CountFail,1)
		return
	}
	fi, err := os.Stat(fullFileName)
	if err == nil {
		fileSize := fi.Size() / 1000
		if fileSize < 1 {
			fmt.Println("Warn", " file size too small: ", fullFileName, "|", fileSize, "kB", "|", url)
			atomic.AddInt64(&CountFail,1)
		}
		fmt.Println("GOT <- ", fullFileName, "|", fileSize, "kB", "|", url)
		atomic.AddInt64(&CountSuccess,1)
	} else {
		fmt.Println("Warn", " can not get file size : ", fullFileName, "|", url)
		atomic.AddInt64(&CountFail,1)
	}
}

// 提取图集描述信息
func ExtractGalleryMeta(root soup.Root, url string) (folder string, e bool) {
	defer DeferRecover()
	e = true
	tuji := root.Find("div", "class", "tuji")
	title := LegalPathName(url)
	to_write := make([]string, 0)
	for _, child := range tuji.Children() {
		line := strings.ReplaceAll(child.FullText(), "\n", "")
		if len(line) == 0 {
			continue
		}
		if child.Attrs()["class"] == "weizhi" {
			titleLocation := strings.Split(line, ">")
			title = strings.Trim(titleLocation[len(titleLocation)-1], " ")
			title = strings.Trim(strings.Split(title, "/")[0], " ")
		}
		if len(line) > 0 {
			//fmt.Println(line)
			to_write = append(to_write, line)
		}
	}
	folder = path.Join(DirDownImg, title)
	WriteTxt(folder, "describe.txt", to_write, false)
	e = false
	return
}

// 写文本文件
func WriteTxt(folder string, fileName string, content []string, overwrite bool) {
	defer DeferRecover()
	fullFileName := path.Join(folder, fileName)
	if !overwrite && FileExists(fullFileName) {
		return
	}
	_ = os.MkdirAll(folder, os.ModePerm)
	file, err := os.OpenFile(fullFileName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	defer file.Close()
	if err != nil {
		//log.Println("Error", "file open:", err)
		return
	}
	writer := bufio.NewWriter(file)
	for _, line := range content {
		writer.WriteString(line + "\n")
	}
	//Flush将缓存的文件真正写入到文件中
	writer.Flush()
}
