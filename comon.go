package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ApiClientBase struct {
	Url               string
	ServiceName       string
}

func (api *ApiClientBase) get() *goquery.Document {
	res, _ := http.Get(api.Url)
	defer res.Body.Close()

	reader := changeTextCode(res)
	doc, _ := goquery.NewDocumentFromReader(reader)
	return doc
}

func changeTextCode (res *http.Response) io.Reader {
	buf, _ := ioutil.ReadAll(res.Body)

	det := chardet.NewTextDetector()
	detRslt, _ := det.DetectBest(buf)
	bReader := bytes.NewReader(buf)
	reader, _ := charset.NewReaderLabel(detRslt.Charset, bReader)
	return reader
}

func  RemoveBlankStrings(str string) string{
	noTabText := strings.Replace(str, "\t", "", -1)
	return strings.Replace(noTabText, "\n", "", -1)
}

func getPrice(str string) int {
	noMark := strings.Replace(str, "¥", "", 1)
	priceStr := RemoveBlankStrings(noMark)
	price, _ := strconv.Atoi(priceStr)
	return price
}