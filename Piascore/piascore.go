package Piascore

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	Common "../common"
	"github.com/PuerkitoBio/goquery"
)

const piascore string = "Piascore"
const piascoreDomain = "https://store.piascore.com"
const piascoreUrlBase string = piascoreDomain + "/search?"

type piascoreClient struct {
	Common.ApiClientBase
}

type Piascore struct {
	*Common.BasicInfo
	Difficulty string
}

func GetInfo(instrument string, music string) []Piascore {
	piascoreClient := initPiascoreClient(instrument, music)
	doc := piascoreClient.Get()
	return *getTableInfo(doc, instrument)
}

func initPiascoreClient(instrument string, music string) *Common.ApiClientBase {
	psc := piascoreClient{}
	var acb *Common.ApiClientBase
	acb = &Common.ApiClientBase{
		ServiceName: piascore,
	}
	psc.ApiClientBase = *acb
	psc.setUrl(instrument, music)
	return &psc.ApiClientBase

}

func (pmc *piascoreClient) setUrl(instrument string, music string) {
	var url string
	switch Common.WhichInstrumentType(instrument) {
	case "Saxophone":
		url = setSaxUrl(instrument, music)
	default:
		url = piascoreUrlBase
	}
	pmc.ApiClientBase.Url = url
}

func setSaxUrl(instrument string, music string) string {
	saxType := map[string]string{
		"sopranoSaxophone":  "i=326&",
		"altoSaxophone":     "i=320&",
		"tenorSaxophone":    "i=322&",
		"baritoneSaxophone": "i=324&",
	}
	music = "n=" + music
	return piascoreUrlBase + saxType[instrument] + music
}

func initPiascore(instrument string) *Piascore {
	psc := Piascore{}
	var cbi *Common.BasicInfo
	cbi = &Common.BasicInfo{
		ServiceName: piascore,
		Instrument:  instrument,
	}
	psc.BasicInfo = cbi
	return &psc
}

func getTableInfo(gd *goquery.Document, instrument string) *[]Piascore {
	var respondArray []Piascore
	gd.Find(".displayed-score").EachWithBreak(func(i int, div *goquery.Selection) bool {
		psc := initPiascore(instrument)
		psc.setPiascoreInfo(div)
		psc.setBasicInfo(div)
		psc.Output()
		return true
	})
	return &respondArray
}

func (psc *Piascore) setBasicInfo(gs *goquery.Selection) *Piascore {
	path, _ := gs.Find(".top-score-item").Attr("href")
	url := piascoreDomain + path
	psc.BasicInfo.Url = url
	psc.BasicInfo.Composer = "-"
	psc.setPrice(gs)
	gs.Find(".score_list_inst").Remove()
	gs.Find(".score_list_price").Remove()
	psc.BasicInfo.MusicName = gs.Find(".score_list_title").Text()
	return psc
}

func (psc *Piascore) setPiascoreInfo(gs *goquery.Selection) *Piascore {
	psc.Difficulty = Common.RemoveBlankStrings(gs.Find(".score_list_info").Find(".score_list_inst").Text())
	return psc
}

func (psc *Piascore) setPrice(gs *goquery.Selection) {
	//priceがなぜか"¥100¥100"と言う謎の形で入ってくるため、仕方なくスプリット
	price := gs.Find(".score_list_price").Text()
	if price == "無料無料" {
		psc.BasicInfo.Price = 0
		return
	}
	prices := strings.Split(price, "¥")
	pri, _ := strconv.Atoi(prices[1])
	psc.BasicInfo.Price = pri

}

func Main() {
	fmt.Println("start")
	insts := Common.Instruments()
	for _, installment := range insts {
		GetAllMusicScore(installment)
	}
}

func GetAllMusicScore(instrument string) {
	client := setUrl(instrument)
	doc := client.Get()
	getAll(doc, client, client.Url, instrument)
}

func setUrl(instrument string) *Common.ApiClientBase {
	url := instalmentUrl(instrument)
	var cab *Common.ApiClientBase
	cab = &Common.ApiClientBase{
		Url: url,
	}
	return cab
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

func getAll(doc *goquery.Document, client *Common.ApiClientBase, url string, instrument string) {
	fmt.Println("getAll start")
	path, exist := doc.Find(".pagination").Children().Last().Children().Attr("href")
	count := 1
	if exist {
		count = getMaxCount(path)
	}
	for i := 1; i <= count; i++ {
		client.Url = url + "page=" + strconv.Itoa(i)
		fmt.Println(client.Url)
		doc = client.Get()
		getTableInfo(doc, instrument)

	}
	fmt.Println("getAll end")
}

func (psc *Piascore) Output() {
	file, err := os.OpenFile("test.csv", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		//エラー処理
		log.Fatal(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	var s = []string{
		psc.BasicInfo.ServiceName,
		psc.BasicInfo.MusicName,
		psc.BasicInfo.Composer,
		strconv.Itoa(psc.BasicInfo.Price),
		psc.BasicInfo.Url,
		psc.BasicInfo.Instrument,
		psc.Difficulty,
	}
	writer.Write(s)
	writer.Flush()
}

func getMaxCount(path string) int {
	pages := strings.Split(path, "page=")
	fmt.Println(pages[1])
	c := pages[1]
	count := 1
	if c != "" {
		count, _ = strconv.Atoi(c)
	}
	return count
}
