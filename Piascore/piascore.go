package Piascore

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"

	"github.com/PuerkitoBio/goquery"
)

const piascore string = "Piascore"
const piascoreDomain = "https://store.piascore.com"

type Piascore struct {
	Difficulty  string
	ServiceName string
	MusicName   string
	Composer    string
	Price       int
	Url         string
	Instrument  string
	Id          string
}

func initPiascore(instrument string) *Piascore {
	psc := Piascore{
		ServiceName: piascore,
		Instrument:  instrument,
		Url:         instalmentUrl(instrument),
	}
	return &psc
}

func (psc *Piascore) setInfo(gd *goquery.Document) []Piascore {
	pscs := make([]Piascore, 0)
	gd.Find(".displayed-score").EachWithBreak(func(i int, div *goquery.Selection) bool {
		dp := psc.goToDetailPAge(div)
		psc.setInfoInner(dp)
		pscs = append(pscs, *psc)
		return true
	})
	return pscs
}

func (psc *Piascore) setInfoInner(gs *goquery.Document) *Piascore {
	psc.setId()
	psc.Composer = gs.Find(".score-persons").Children().First().Find(".person_name").Text()
	psc.setPrice(gs)
	psc.MusicName = gs.Find(".score_title").Text()
	psc.Difficulty = gs.Find(".score-breadcrumb").Last().Text()
	return psc
}

func (psc *Piascore) setPrice(gs *goquery.Document) {
	price := gs.Find(".score_price").Text()
	if price == "" {
		psc.Price = 0
		return
	}
	price = strings.Replace(price, "¥", "", -1)
	price = strings.Replace(price, "（税込）", "", -1)
	price = strings.Replace(price, ",", "", -1)
	pri, _ := strconv.Atoi(price)
	psc.Price = pri

}

func Main() {
	lambda.Start(start)
}

func start() {
	fmt.Println("start")
	for _, installment := range instruments() {
		fmt.Println(installment)
		GetAllMusicScore(installment)
	}
	fmt.Println("end")
}

func GetAllMusicScore(instrument string) {
	psc := initPiascore(instrument)
	doc := psc.Get()
	psc.bringAllData(doc)
}

func (psc *Piascore) setUrl(instrument string) {
	url := instalmentUrl(instrument)
	psc.Url = url
}

func instalmentUrl(instrument string) string {
	instalments := map[string]string{
		"sopranoSaxophone":  "https://store.piascore.com/search?i=326&",
		"altoSaxophone":     "https://store.piascore.com/search?i=320&",
		"tenorSaxophone":    "https://store.piascore.com/search?i=322&",
		"baritoneSaxophone": "https://store.piascore.com/search?i=324&",
	}
	return instalments[instrument]
}

func (psc *Piascore) bringAllData(doc *goquery.Document) {
	path, exist := doc.Find(".pagination").Children().Last().Children().Attr("href")
	count := 1
	if exist {
		count = getMaxPage(path)
	}
	pscs := make([]Piascore, 0)
	for i := 1; i <= count; i++ {
		psc.Url = instalmentUrl(psc.Instrument) + "page=" + strconv.Itoa(i)
		fmt.Println(psc.Url)
		doc = psc.Get()
		ps := psc.setInfo(doc)
		pscs = append(pscs, ps...)
	}
	output(pscs)
}

func output(pscs []Piascore) {
	file, err := os.OpenFile("test.csv", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for _, psc := range pscs {
		writer := csv.NewWriter(file)
		var s = []string{
			psc.ServiceName,
			psc.MusicName,
			psc.Composer,
			strconv.Itoa(psc.Price),
			psc.Url,
			psc.Instrument,
			psc.Difficulty,
			psc.Id,
		}
		writer.Write(s)
		writer.Flush()
	}
}

func getMaxPage(path string) int {
	pages := strings.Split(path, "page=")
	c := pages[1]
	count := 1
	if c != "" {
		count, _ = strconv.Atoi(c)
	}
	return count
}

func (psc *Piascore) goToDetailPAge(div *goquery.Selection) (gs *goquery.Document) {
	path, _ := div.Find(".top-score-item").Attr("href")
	psc.Url = piascoreDomain + path
	gs = psc.Get()
	return gs

}

func (psc *Piascore) setId() {
	slice := strings.Split(psc.Url, "/")
	psc.Id = slice[len(slice)-1]
}

func (psc *Piascore) Get() *goquery.Document {
	client := &http.Client{}

	req, err := http.NewRequest("GET", psc.Url, nil)
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

func changeTextCode(res *http.Response) io.Reader {
	buf, _ := ioutil.ReadAll(res.Body)

	det := chardet.NewTextDetector()
	detRslt, _ := det.DetectBest(buf)
	bReader := bytes.NewReader(buf)
	reader, _ := charset.NewReaderLabel(detRslt.Charset, bReader)
	return reader
}

func instruments() []string {
	slice := []string{
		"sopranoSaxophone",
		"altoSaxophone",
		"tenorSaxophone",
		"baritoneSaxophone",
	}
	return slice
}
