package Common

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
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
	Id          string
}

func (api *ApiClientBase) Get() *goquery.Document {
	client := &http.Client{}

	req, err := http.NewRequest("GET", api.Url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	ua := getUserAgent()
	req.Header.Set("User-Agent", ua)
	req.Header.Add("Accept-Language", "ja,en-US;q=0.9,en;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	reader := changeTextCode(resp)
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

func getUserAgent() string {
	slice := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3864.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:62.0) Gecko/20100101 Firefox/62.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:67.0) Gecko/20100101 Firefox/67.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:68.0) Gecko/20100101 Firefox/68.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:61.0) Gecko/20100101 Firefox/61.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/17.17134",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Safari/605.1.15",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/1.6.5b18.09.26.16 Mobile/16A366 Safari/605.1.15 _id/000002",
	}
	i := rand.Intn(15)
	return slice[i]
}
