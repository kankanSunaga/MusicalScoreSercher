package Elise

import (
	Common "../common"
	"github.com/PuerkitoBio/goquery"
)

const elise string = "at-elise"
const eliseDomain = "https://www.at-elise.com"
const eliseUrlBase string = eliseDomain + "/elise/Services.SvSession?method=GakufuSearch&F_TEMPLATE=search_ret.htx&F_CATEGORY1=Gd_Instrument"

type eliseClient struct {
	Common.ApiClientBase
}

type Elise struct {
	*Common.BasicInfo
	Difficulty string
}

func GetInfo(instrument string, music string) []Elise {
	eliseClient := initEliseClient(instrument, music)
	doc := eliseClient.Get()
	musicTable := getEliseTable(doc)
	return *getTableInfo(musicTable)
}

func initEliseClient(instrument string, music string) *Common.ApiClientBase {
	elc := eliseClient{}
	var acb *Common.ApiClientBase
	acb = &Common.ApiClientBase{
		ServiceName: elise,
	}
	elc.ApiClientBase = *acb
	elc.setUrl(instrument, music)
	return &elc.ApiClientBase

}

func (pmc *eliseClient) setUrl(instrument string, music string) {
	var url string
	switch Common.WhichInstrumentType(instrument) {
	case "Saxophone":
		url = setSaxUrl(instrument, music)
	default:
		url = eliseUrlBase
	}
	pmc.ApiClientBase.Url = url
}

func setSaxUrl(instrument string, music string) string {
	saxType := map[string]string{
		"sopranoSaxophone":  "&F_GENRE1=27",
		"altoSaxophone":     "&F_GENRE1=28",
		"tenorSaxophone":    "&F_GENRE1=29",
		"baritoneSaxophone": "&F_GENRE1=30",
	}
	music = "&F_KEYWORD=" + music
	return eliseUrlBase + music + saxType[instrument]
}

func getEliseTable(gd *goquery.Document) *goquery.Selection {

	table := gd.Find("#s_result")
	table.Find("#s_result_head").Remove()
	return table
}

func initElise() *Elise {
	eli := Elise{}
	var cbi *Common.BasicInfo
	cbi = &Common.BasicInfo{
		ServiceName: elise,
	}
	eli.BasicInfo = cbi
	return &eli
}

func getTableInfo(gd *goquery.Selection) *[]Elise {
	var respondArray []Elise
	gd.Find(".s_result_content").EachWithBreak(func(i int, tr *goquery.Selection) bool {
		eli := initElise()
		eli.setBasicInfo(tr)
		eli.setEliseInfo(tr)
		respondArray = append(respondArray, *eli)
		return true
	})
	return &respondArray
}

func (eli *Elise) setBasicInfo(gs *goquery.Selection) *Elise {
	path, _ := gs.Find(".title").Find("a").Attr("href")
	url := eliseDomain + path
	eli.BasicInfo.Url = url
	eli.BasicInfo.MusicName = gs.Find(".title").Text()
	eli.BasicInfo.Composer = gs.Find(".artist").Text()
	eli.BasicInfo.Price = Common.GetPrice(gs.Find(".costcv").Text())
	return eli
}

func (pms *Elise) setEliseInfo(gs *goquery.Selection) *Elise {
	pms.Difficulty = Common.RemoveBlankStrings(gs.Find(".inst").Text())
	return pms
}
