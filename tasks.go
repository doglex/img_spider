package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/levigross/grequests"
	"log"
	"strings"
	"time"
)

func RunListSexy() {
	set_galleray := make(map[string]int, 0)
	template := "https://www.tujigu.com/s/35/index_%v.html"
	idx_0 := 1
	idx_1 := 44
	fmt.Println("开始获取所有专辑。。。")
	for i := idx_0; i <= idx_1; i++ {
		mainPage := fmt.Sprintf(template, i)
		ro := grequests.RequestOptions{
			RequestTimeout: time.Second * 180,
			Headers: map[string]string{
				"Referer":    "https://www.baidu.com",
				"User-Agent": GetUA(),
			}}
		resp, err := grequests.Get(mainPage, &ro)
		if err != nil {
			log.Fatalln("Unable to make request: ", err)
		}
		respString := resp.String()
		doc := soup.HTMLParse(respString)
		galleries := doc.Find("div", "class", "hezi").FindAll("a")
		for _, gallery := range galleries {
			gallery_url := gallery.Attrs()["href"]
			if !strings.Contains(gallery_url, "/a/") {
				continue
			}
			_, ok := set_galleray[gallery_url]
			if ok {
				continue
			}
			set_galleray[gallery_url] = 1
			fmt.Println(gallery_url)
		}
	}
	n := len(set_galleray)
	fmt.Println("专辑总量", n, "\n\n\n")
	i := 0
	for k, _ := range set_galleray {
		i += 1
		k := k
		to_print := fmt.Sprintf("[%v / %v]开始抓取-> %v", i, n, k)
		fmt.Println(to_print)
		go WalkGallery(k)
		//time.Sleep(time.Second)
		time.Sleep(2 * time.Microsecond) // 保持可用
	}
	time.Sleep(2 * time.Second) // 保持可用
}

func RunListCN() {
	set_galleray := make(map[string]int, 0)
	template := "https://www.tujigu.com/zhongguo/%v.html"
	idx_0 := 1
	idx_1 := 50
	fmt.Println("开始获取所有专辑。。。")
	for i := idx_0; i <= idx_1; i++ {
		mainPage := fmt.Sprintf(template, i)
		if i == 1 {
			mainPage = "https://www.tujigu.com/zhongguo/"
		}
		ro := grequests.RequestOptions{
			RequestTimeout: time.Second * 180,
			Headers: map[string]string{
				"Referer":    "https://www.baidu.com",
				"User-Agent": GetUA(),
			}}
		resp, err := grequests.Get(mainPage, &ro)
		if err != nil {
			log.Fatalln("Unable to make request: ", err)
		}
		respString := resp.String()
		doc := soup.HTMLParse(respString)
		galleries := doc.Find("div", "class", "hezi").FindAll("a")
		for _, gallery := range galleries {
			gallery_url := gallery.Attrs()["href"]
			if !strings.Contains(gallery_url, "/a/") {
				continue
			}
			_, ok := set_galleray[gallery_url]
			if ok {
				continue
			}
			set_galleray[gallery_url] = 1
			fmt.Println(gallery_url)
		}
	}
	n := len(set_galleray)
	fmt.Println("专辑总量", n, "\n\n\n")
	i := 0
	for k, _ := range set_galleray {
		i += 1
		k := k
		to_print := fmt.Sprintf("[%v / %v]开始抓取-> %v", i, n, k)
		fmt.Println(to_print)
		go WalkGallery(k)
		//time.Sleep(time.Second)
		time.Sleep(100 * time.Microsecond) // 保持可用
	}
	time.Sleep(2 * time.Second) // 保持可用
}

func RunListCompany() {

	companies := make([]string, 0)
	tempX := "https://www.tujigu.com/x/%v/"
	tempY := "https://www.tujigu.com/x/%v/index_%v.html"
	for i := 1; i <= 101; i++ {
		url_company := fmt.Sprintf(tempX, i)
		companies = append(companies, url_company)
		for j := 1; j < 200; j++ {
			url_company := fmt.Sprintf(tempY, i, j)
			ro := grequests.RequestOptions{
				RequestTimeout: time.Second * 180,
				Headers: map[string]string{
					"Referer":    "https://www.tujigu.com/",
					"User-Agent": GetUA(),
				}}
			resp, err := grequests.Get(url_company, &ro)
			if err != nil {
				log.Fatalln("Unable to make request: ", err)
			}
			respString := resp.String()
			if strings.Contains(respString, "网站改版、请从首页重新访问！ 返回首页") {
				fmt.Println("超出了", i, j)
				break
			} else {
				companies = append(companies, url_company)
				fmt.Println(url_company)
			}
		}
	}
	fmt.Println(companies)

	set_galleray := make(map[string]int, 0)
	fmt.Println("开始获取所有专辑。。。")
	for _, mainPage := range companies {
		ro := grequests.RequestOptions{
			RequestTimeout: time.Second * 180,
			Headers: map[string]string{
				"Referer":    "https://www.tujigu.com/",
				"User-Agent": GetUA(),
			}}
		resp, err := grequests.Get(mainPage, &ro)
		if err != nil {
			log.Fatalln("Unable to make request: ", err)
		}
		respString := resp.String()
		doc := soup.HTMLParse(respString)
		galleries := doc.Find("div", "class", "hezi").FindAll("a")
		for _, gallery := range galleries {
			gallery_url := gallery.Attrs()["href"]
			if !strings.Contains(gallery_url, "/a/") {
				continue
			}
			_, ok := set_galleray[gallery_url]
			if ok {
				continue
			}
			set_galleray[gallery_url] = 1
			fmt.Println(gallery_url)
		}
	}
	n := len(set_galleray)
	fmt.Println("专辑总量", n, "\n\n\n")
	i := 0
	for k, _ := range set_galleray {
		i += 1
		k := k
		to_print := fmt.Sprintf("[%v / %v]开始抓取-> %v", i, n, k)
		fmt.Println(to_print)
		go WalkGallery(k)
		//time.Sleep(time.Second)
		time.Sleep(2 * time.Microsecond) // 保持可用
	}
	time.Sleep(2 * time.Second) // 保持可用
}
