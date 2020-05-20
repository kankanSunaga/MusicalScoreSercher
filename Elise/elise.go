package Elise

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

	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"

	Common "../common"
	"github.com/PuerkitoBio/goquery"
)

const elise string = "at-elise"
const eliseDomain = "https://www.at-elise.com"
const eliseUrlBase string = eliseDomain + "/elise/Services.SvSession?method=GakufuSearch&F_TEMPLATE=search_ret.htx&F_CATEGORY1=Gd_Instrument"

type Elise struct {
	Difficulty  string
	ServiceName string
	MusicName   string
	Composer    string
	Price       int
	Url         string
	Instrument  string
	Id          string
}

func Main() {
	fmt.Println("start")
	insts := instruments()
	for _, installment := range insts {
		fmt.Println("start", installment)
		GetInfo(installment)
	}
}

func GetInfo(instrument string) {
	eli := initElise(instrument)
	doc := eli.Get()
	musicTable := getEliseTable(doc)
	eli.getTableInfo(musicTable)
}

func instalmentUrl(instrument string) string {
	saxType := map[string]string{
		"sopranoSaxophone":  "&F_GENRE1=27",
		"altoSaxophone":     "&F_GENRE1=28",
		"tenorSaxophone":    "&F_GENRE1=29",
		"baritoneSaxophone": "&F_GENRE1=30",
	}
	return eliseUrlBase + saxType[instrument]
}

func getEliseTable(gd *goquery.Document) *goquery.Selection {

	table := gd.Find("#s_result")
	table.Find("#s_result_head").Remove()
	return table
}

func initElise(instrument string) *Elise {
	psc := Elise{
		ServiceName: elise,
		Instrument:  instrument,
		Url:         instalmentUrl(instrument),
	}
	return &psc
}

func (eli *Elise) getTableInfo(gd *goquery.Selection) {
	elis := make([]Elise, 0)
	gd.Find(".s_result_content").EachWithBreak(func(i int, tr *goquery.Selection) bool {
		eli.setInfo(tr)
		elis = append(elis, *eli)
		return true
	})
	output(elis)

}

func (eli *Elise) setInfo(gs *goquery.Selection) *Elise {
	eli.Difficulty = Common.RemoveBlankStrings(gs.Find(".inst").Text())
	path, _ := gs.Find(".title").Find("a").Attr("href")
	url := eliseDomain + path
	eli.Url = url
	eli.Id = getId(path)
	eli.MusicName = gs.Find(".title").Text()
	eli.Composer = gs.Find(".artist").Text()
	eli.Price = Common.GetPrice(gs.Find(".costcv").Text())
	return eli
}

func output(elis []Elise) {
	file, err := os.OpenFile("test.csv", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		//エラー処理
		log.Fatal(err)
	}
	defer file.Close()
	for _, eli := range elis {
		writer := csv.NewWriter(file)
		var s = []string{
			eli.ServiceName,
			eli.MusicName,
			eli.Composer,
			strconv.Itoa(eli.Price),
			eli.Url,
			eli.Instrument,
			eli.Difficulty,
			eli.Id,
		}
		writer.Write(s)
		writer.Flush()
	}
}

func getId(path string) string {
	path = strings.TrimRight(path, "/")
	slice := strings.Split(path, "/")
	return slice[len(slice)-1]
}

func (eli *Elise) Get() *goquery.Document {
	client := &http.Client{}

	req, err := http.NewRequest("GET", eli.Url, nil)
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
