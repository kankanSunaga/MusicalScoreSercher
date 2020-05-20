package PrintMusicalScore

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

const printMusicScore string = "ぷりんと楽譜"
const printMusicScoreDomain = "https://www.print-gakufu.com"

type PrintMusicScore struct {
	Difficulty  string
	ServiceName string
	MusicName   string
	Composer    string
	Price       int
	Url         string
	Instrument  string
	Id          string
}

func initPrintMusicScore(instrument string) *PrintMusicScore {
	printGakufu := PrintMusicScore{
		ServiceName: printMusicScore,
		Instrument:  instrument,
		Url:         instalmentUrl(instrument),
	}
	return &printGakufu
}

func (pms *PrintMusicScore) getTableInfo(gd *goquery.Selection) []PrintMusicScore {
	pmss := make([]PrintMusicScore, 0)
	gd.Find("tr").EachWithBreak(func(i int, tr *goquery.Selection) bool {
		if tr.Is("td") {
			return false
		}
		pms.setInfo(tr)
		pmss = append(pmss, *pms)
		return true
	})
	return pmss
}

func (pms *PrintMusicScore) setInfo(gs *goquery.Selection) *PrintMusicScore {
	path, _ := gs.Find(".title").Find("a").Attr("href")
	pms.Url = printMusicScoreDomain + path
	noBlank := Common.RemoveBlankStrings(gs.Find(".title").Text())
	musicAndComposer := strings.SplitN(noBlank, "/", 2)
	pms.MusicName = musicAndComposer[0]
	pms.Composer = musicAndComposer[1]
	pms.Price = Common.GetPrice(gs.Find(".price").Text())
	pms.Difficulty = Common.RemoveBlankStrings(gs.Find(".status").Text())
	pms.Id = getId(path)
	return pms
}

func Main() {
	insts := Common.Instruments()
	for _, installment := range insts {
		fmt.Println("start", "-", installment)
		GetAllMusicScore(installment)
	}
	fmt.Println("end")

}

func GetAllMusicScore(instrument string) {
	pms := initPrintMusicScore(instrument)
	doc := pms.Get()
	pms.getAll(doc)
}

func instalmentUrl(instrument string) string {
	instalments := map[string]string{
		"sopranoSaxophone":  "https://www.print-gakufu.com/search/result/inst__brass-ss/?",
		"altoSaxophone":     "https://www.print-gakufu.com/search/result/inst__brass-as/?",
		"tenorSaxophone":    "https://www.print-gakufu.com/search/result/inst__brass-ts/?",
		"baritoneSaxophone": "https://www.print-gakufu.com/search/result/inst__brass-bs/?",
	}
	return instalments[instrument]
}

func (pms *PrintMusicScore) getAll(doc *goquery.Document) {
	c := doc.Find(".headInfo").Find(".pageFeed-item-last").Text()
	count := 1
	if c != "" {
		count, _ = strconv.Atoi(c)
	}
	pmss := make([]PrintMusicScore, 0)
	for i := 1; i <= count; i++ {
		pms.Url = instalmentUrl(pms.Instrument) + "p=" + strconv.Itoa(i)
		fmt.Println(pms.Url)
		doc = pms.Get()
		musicTable := getPrintMusicScoreTable(doc)
		ps := pms.getTableInfo(musicTable)
		pmss = append(pmss, ps...)

	}
	output(pmss)
}

func getPrintMusicScoreTable(gd *goquery.Document) *goquery.Selection {
	gd.Find("#footerContainer").Remove()
	musicTable := gd.Find("tbody")
	musicTable.Find(".thead").Remove()
	return musicTable
}
func output(pmss []PrintMusicScore) {
	file, err := os.OpenFile("test.csv", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		//エラー処理
		log.Fatal(err)
	}
	defer file.Close()
	for _, pms := range pmss {
		writer := csv.NewWriter(file)
		var s = []string{
			pms.ServiceName,
			pms.MusicName,
			pms.Composer,
			strconv.Itoa(pms.Price),
			pms.Url,
			pms.Instrument,
			pms.Difficulty,
			pms.Id,
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

func (pms *PrintMusicScore) Get() *goquery.Document {
	client := &http.Client{}

	req, err := http.NewRequest("GET", pms.Url, nil)
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
