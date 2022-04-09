package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var client http.Client//构造客户端

func main() {
	var titles []string
	var contents []string
	title,content,err:=Spider()
	if err != nil {
		fmt.Println("spider err",err)
	}
	titles=append(titles,title)
	contents=append(contents,content)
	SaveJoke2File(1,titles,contents)
}

func Spider()(title, content string, err error) {
	result, err := Get("http://xiaodiaodaya.cn/article/view.aspx?id=174")
	if err != nil {
		fmt.Println("Get err",err)
		return
	}
	//解析、编译正则表达式， 处理 title
	// <a href="https://m.pengfue.com/content/1857797/" title="标题">
	//<h2 class="titleview">那天去公园儿，看见一个...</h2>
	//<h2 class="titleview">=(?s:(.*?))</h2>
	ret1 := regexp.MustCompile(`<h2 class="titleview">=(?s:(.*?))</h2>`)
	if ret1 == nil {
		err = fmt.Errorf("%s", "MustCompile err")
		return
	}
	// 提取 title
	// 有两处，取第一处
	tmpTitle := ret1.FindAllStringSubmatch(result, 1)
	for _, data := range tmpTitle {
		// 存至返回值 title
		title = data[1]
		title = strings.Replace(title, "\t", "", -1)
		// 取一个即可。
		break
	}

	//解析、编译正则表达式， 处理 content
	// `<div class="con-txt">正文内容</div>`
	//<p align="center"></p><!--listS-->(?s:(.*?))<!--listE-->
	ret2 := regexp.MustCompile(`<p align="center"></p><!--listS-->(?s:(.*?))<!--listE-->`)
	if ret2 == nil {
		err = fmt.Errorf("%s", "MustCompile err")
		return
	}
	// 提取 Content
	tmpContent := ret2.FindAllStringSubmatch(result, -1)
	for _, data := range tmpContent {
		// 存至返回值 content
		content = data[1]
		content = strings.Replace(content, "\t", "", -1)
		// 提取一个即可。
		break
	}
	return
}


func Get(url string)(result string,err error)  {
	client:=&http.Client{}
	req,err:=http.NewRequest("GET",url,nil)
	if err != nil {
		fmt.Println("rep err",err)
	}
	//防止爬虫被**网站搞
	req.Header.Set("Connection","keep-alive")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding","gzip, deflate")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("resport err",err)
	}
	defer resp.Body.Close()
	buf := make([]byte, 4096)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}
	return
}

func SaveJoke2File(idx int, fileTitle, fileContent []string) {
	f, err := os.Create(strconv.Itoa(idx) + ".txt")
	if err != nil {
		fmt.Println("Create err:", err)
		return
	}
	defer f.Close()

	n := len(fileTitle)
	for i:=0; i<n; i++ {
		// 写入标题
		f.WriteString(fileTitle[i] + "\n")
		// 写入内容
		f.WriteString(fileContent[i] + "\n")
		// 写一个华丽分割线
		f.WriteString("--------------------------------------------------------------\n")
	}
}


