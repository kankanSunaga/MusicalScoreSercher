package Piascore

import (
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
	return *getTableInfo(doc)
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

func initPiascore() *Piascore {
	psc := Piascore{}
	var cbi *Common.BasicInfo
	cbi = &Common.BasicInfo{
		ServiceName: piascore,
	}
	psc.BasicInfo = cbi
	return &psc
}

func getTableInfo(gd *goquery.Document) *[]Piascore {
	var respondArray []Piascore
	gd.Find(".displayed-score").EachWithBreak(func(i int, div *goquery.Selection) bool {
		pms := initPiascore()
		pms.setPiascoreInfo(div)
		pms.setBasicInfo(div)
		respondArray = append(respondArray, *pms)
		return true
	})
	return &respondArray
}

func (pms *Piascore) setBasicInfo(gs *goquery.Selection) *Piascore {
	path, _ := gs.Find(".top-score-item").Attr("href")
	url := piascoreDomain + path
	pms.BasicInfo.Url = url
	pms.BasicInfo.Composer = "-"
	pms.setPrice(gs)
	gs.Find(".score_list_inst").Remove()
	gs.Find(".score_list_price").Remove()
	pms.BasicInfo.MusicName = gs.Find(".score_list_title").Text()
	return pms
}

func (pms *Piascore) setPiascoreInfo(gs *goquery.Selection) *Piascore {
	pms.Difficulty = Common.RemoveBlankStrings(gs.Find(".score_list_info").Find(".score_list_inst").Text())
	return pms
}

func (pms *Piascore) setPrice(gs *goquery.Selection) {
	//priceがなぜか"¥100¥100"と言う謎の形で入ってくるため、仕方なくスプリット
	prices := strings.Split(gs.Find(".score_list_price").Text(), "¥")
	price, _ := strconv.Atoi(prices[1])
	pms.BasicInfo.Price = price
}
