package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"strings"
)

type BasicInfo struct {
	ServiceName string
	MusicName string
	Composer  string
	Price     int
	Url       string
}

type PrintGakufu struct {
	BasicInfo
	Difficulty string
}

func InitPrintGakufu() *PrintGakufu {
	printGakufu := PrintGakufu{}
	printGakufu.BasicInfo.ServiceName = "ぷりんと楽譜"
	return &printGakufu
}

func main() {
	url := "https://www.print-gakufu.com/search/result/song__ルパン__is_part___inst__brass-ts/"

	// Getリクエスト
	res, _ := http.Get(url)
	defer res.Body.Close()

	// 読み取り
	buf, _ := ioutil.ReadAll(res.Body)


	// 文字コード判定
	det := chardet.NewTextDetector()
	detRslt, _ := det.DetectBest(buf)
	fmt.Println(detRslt.Charset)
	// => EUC-JP

	// 文字コード変換
	bReader := bytes.NewReader(buf)
	reader, _ := charset.NewReaderLabel(detRslt.Charset, bReader)

	// HTMLパース
	doc, _ := goquery.NewDocumentFromReader(reader)
	doc.Find("#footerContainer").Remove()
	musicTable := doc.Find("tbody")
	// titleを抜き出し
	musicTable.Find(".thead").Remove()
	var respondArray []PrintGakufu
	musicTable.Find("tr").EachWithBreak(func(i int, tr *goquery.Selection) bool {
		if tr.Is("td") {
			return false
		}
		printGakuhu := InitPrintGakufu()
		path, _ :=tr.Find(".title").Find("a").Attr("href")
		url := "https://www.print-gakufu.com" + path
		printGakuhu.BasicInfo.Url = url
		noblank:= RemoveBlankStrings(tr.Find(".title").Text())
		musicAndComposer := strings.SplitN(noblank, "/", 2 )
		printGakuhu.BasicInfo.MusicName = musicAndComposer[0]
		printGakuhu.BasicInfo.Composer = musicAndComposer[1]
		printGakuhu.Difficulty = RemoveBlankStrings(tr.Find(".status").Text())
		printGakuhu.BasicInfo.Price = getPrice(tr.Find(".price").Text())
		respondArray = append(respondArray, *printGakuhu)
		return true
	})
	data, _ := json.Marshal(respondArray)
	fmt.Println(string(data))
}
