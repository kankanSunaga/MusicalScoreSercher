package PrintMusicalScore

import (
	"strings"

	Common "../common"
	"github.com/PuerkitoBio/goquery"
)

const printMusicScore string = "ぷりんと楽譜"
const printMusicScoreDomain = "https://www.print-gakufu.com"
const printMusicScoreUrlBase string = printMusicScoreDomain + "/search/result/song__"
const saxophoneUrl string = "__is_part___inst__brass-"

type printMusicScoreClient struct {
	Common.ApiClientBase
}

type PrintMusicScore struct {
	*Common.BasicInfo
	Difficulty string
}

func GetInfo(instrument string, music string) []PrintMusicScore {
	printMusicScoreClient := initPrintMusicScoreClient(instrument, music)
	doc := printMusicScoreClient.Get()
	musicTable := getPrintMusicScoreTable(doc)
	return *getTableInfo(musicTable)
}

func initPrintMusicScoreClient(instrument string, music string) *Common.ApiClientBase {
	pmc := printMusicScoreClient{}
	var acb *Common.ApiClientBase
	acb = &Common.ApiClientBase{
		ServiceName: printMusicScore,
	}
	pmc.ApiClientBase = *acb
	pmc.setUrl(instrument, music)
	return &pmc.ApiClientBase

}

func (pmc *printMusicScoreClient) setUrl(instrument string, music string) {
	var url string
	switch Common.WhichInstrumentType(instrument) {
	case "Saxophone":
		url = setSaxUrl(instrument, music)
	default:
		url = printMusicScoreUrlBase
	}
	pmc.ApiClientBase.Url = url
}

func setSaxUrl(instrument string, music string) string {
	saxType := map[string]string{
		"sopranoSaxophone":  "ss",
		"altoSaxophone":     "as",
		"tenorSaxophone":    "ts",
		"baritoneSaxophone": "bs",
	}
	return printMusicScoreUrlBase + music + saxophoneUrl + saxType[instrument]
}

func getPrintMusicScoreTable(gd *goquery.Document) *goquery.Selection {
	gd.Find("#footerContainer").Remove()
	musicTable := gd.Find("tbody")
	musicTable.Find(".thead").Remove()
	return musicTable
}

func initPrintMusicScore() *PrintMusicScore {
	printGakufu := PrintMusicScore{}
	var cbi *Common.BasicInfo
	cbi = &Common.BasicInfo{
		ServiceName: printMusicScore,
	}
	printGakufu.BasicInfo = cbi
	return &printGakufu
}

func getTableInfo(gd *goquery.Selection) *[]PrintMusicScore {
	var respondArray []PrintMusicScore
	gd.Find("tr").EachWithBreak(func(i int, tr *goquery.Selection) bool {
		if tr.Is("td") {
			return false
		}
		pms := initPrintMusicScore()
		pms.setBasicInfo(tr)
		pms.setPrintMusicScoreInfo(tr)
		respondArray = append(respondArray, *pms)
		return true
	})
	return &respondArray
}

func (pms *PrintMusicScore) setBasicInfo(gs *goquery.Selection) *PrintMusicScore {
	path, _ := gs.Find(".title").Find("a").Attr("href")
	url := printMusicScoreDomain + path
	pms.BasicInfo.Url = url
	noBlank := Common.RemoveBlankStrings(gs.Find(".title").Text())
	musicAndComposer := strings.SplitN(noBlank, "/", 2)
	pms.BasicInfo.MusicName = musicAndComposer[0]
	pms.BasicInfo.Composer = musicAndComposer[1]
	pms.BasicInfo.Price = Common.GetPrice(gs.Find(".price").Text())
	return pms
}

func (pms *PrintMusicScore) setPrintMusicScoreInfo(gs *goquery.Selection) *PrintMusicScore {
	pms.Difficulty = Common.RemoveBlankStrings(gs.Find(".status").Text())
	return pms
}
