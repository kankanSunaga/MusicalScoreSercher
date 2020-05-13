package Common

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

type ApiClientBase struct {
	Url         string
	ServiceName string
}

type BasicInfo struct {
	ServiceName string
	MusicName   string
	Composer    string
	Price       int
	Url         string
	Instrument  string
}

func (api *ApiClientBase) Get() *goquery.Document {
	res, _ := http.Get(api.Url)
	reader := changeTextCode(res)
	doc, _ := goquery.NewDocumentFromReader(reader)
	return doc
}

func changeTextCode(res *http.Response) io.Reader {
	buf, _ := ioutil.ReadAll(res.Body)

	det := chardet.NewTextDetector()
	detRslt, _ := det.DetectBest(buf)
	bReader := bytes.NewReader(buf)
	reader, _ := charset.NewReaderLabel(detRslt.Charset, bReader)
	return reader
}

func RemoveBlankStrings(str string) string {
	noTabText := strings.Replace(str, "\t", "", -1)
	return strings.Replace(noTabText, "\n", "", -1)
}

func GetPrice(str string) int {
	noYen := strings.Replace(str, "¥", "", 1)
	noMark := strings.Replace(noYen, "円", "", 1)
	priceStr := RemoveBlankStrings(noMark)
	price, _ := strconv.Atoi(priceStr)
	return price
}

func WhichInstrumentType(instrument string) string {
	var itmType string
	if strings.Contains(instrument, "Saxophone") {
		itmType = "Saxophone"
	}
	return itmType
}

func Instruments() []string {
	slice := []string{
		"sopranoSaxophone",
		"altoSaxophone",
		"tenorSaxophone",
		"baritoneSaxophone",
	}
	return slice
}
