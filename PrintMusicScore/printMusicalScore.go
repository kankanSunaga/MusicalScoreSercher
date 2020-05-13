package PrintMusicalScore

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

//func GetInfo(instrument string, music string) []PrintMusicScore {
//	printMusicScoreClient := initPrintMusicScoreClient(instrument, music)
//	doc := printMusicScoreClient.Get()
//	musicTable := GetPrintMusicScoreTable(doc)
//	return *getTableInfo(musicTable)
//}
//
//func initPrintMusicScoreClient(instrument string, music string) *Common.ApiClientBase {
//	pmc := printMusicScoreClient{}
//	var acb *Common.ApiClientBase
//	acb = &Common.ApiClientBase{
//		ServiceName: printMusicScore,
//	}
//	pmc.ApiClientBase = *acb
//	pmc = setUrl(instrument, music)
//	return &pmc.ApiClientBase
//
//}

func setSaxUrl(instrument string, music string) string {
	saxType := map[string]string{
		"sopranoSaxophone":  "ss",
		"altoSaxophone":     "as",
		"tenorSaxophone":    "ts",
		"baritoneSaxophone": "bs",
	}
	return printMusicScoreUrlBase + music + saxophoneUrl + saxType[instrument]
}

func GetPrintMusicScoreTable(gd *goquery.Document) *goquery.Selection {
	gd.Find("#footerContainer").Remove()
	musicTable := gd.Find("tbody")
	musicTable.Find(".thead").Remove()
	return musicTable
}

func initPrintMusicScore(instrument string) *PrintMusicScore {
	printGakufu := PrintMusicScore{}
	var cbi *Common.BasicInfo
	cbi = &Common.BasicInfo{
		ServiceName: printMusicScore,
		Instrument:  instrument,
	}
	printGakufu.BasicInfo = cbi
	return &printGakufu
}

func getTableInfo(gd *goquery.Selection, instrument string) {
	gd.Find("tr").EachWithBreak(func(i int, tr *goquery.Selection) bool {
		if tr.Is("td") {
			return false
		}
		pms := initPrintMusicScore(instrument)
		pms = pms.setInfo(tr)
		return true
	})
}

func (pms *PrintMusicScore) setInfo(gs *goquery.Selection) *PrintMusicScore {
	path, _ := gs.Find(".title").Find("a").Attr("href")
	pms.BasicInfo.Url = printMusicScoreDomain + path
	noBlank := Common.RemoveBlankStrings(gs.Find(".title").Text())
	musicAndComposer := strings.SplitN(noBlank, "/", 2)
	pms.BasicInfo.MusicName = musicAndComposer[0]
	pms.BasicInfo.Composer = musicAndComposer[1]
	pms.BasicInfo.Price = Common.GetPrice(gs.Find(".price").Text())
	pms.Difficulty = Common.RemoveBlankStrings(gs.Find(".status").Text())
	pms.Output()
	return pms
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
	fmt.Println(" setUrl done")
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
		"sopranoSaxophone":  "https://www.print-gakufu.com/search/result/inst__brass-ss/?",
		"altoSaxophone":     "https://www.print-gakufu.com/search/result/inst__brass-as/?",
		"tenorSaxophone":    "https://www.print-gakufu.com/search/result/inst__brass-ts/?",
		"baritoneSaxophone": "https://www.print-gakufu.com/search/result/inst__brass-bs/?",
	}
	return instalments[instrument]
}

func getAll(doc *goquery.Document, client *Common.ApiClientBase, url string, instrument string) {
	fmt.Println("getAll start")
	c := doc.Find(".headInfo").Find(".pageFeed-item-last").Text()
	count := 1
	if c != "" {
		count, _ = strconv.Atoi(c)
	}

	for i := 1; i <= count; i++ {
		client.Url = url + "p=" + strconv.Itoa(i)
		fmt.Println(client.Url)
		doc = client.Get()
		musicTable := getPrintMusicScoreTable(doc)
		getTableInfo(musicTable, instrument)

	}
	fmt.Println("getAll end")
}

func getPrintMusicScoreTable(gd *goquery.Document) *goquery.Selection {
	gd.Find("#footerContainer").Remove()
	musicTable := gd.Find("tbody")
	musicTable.Find(".thead").Remove()
	return musicTable
}
func (pms *PrintMusicScore) Output() {
	file, err := os.OpenFile("test.csv", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		//エラー処理
		log.Fatal(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	var s = []string{
		pms.BasicInfo.ServiceName,
		pms.BasicInfo.MusicName,
		pms.BasicInfo.Composer,
		strconv.Itoa(pms.BasicInfo.Price),
		pms.BasicInfo.Url,
		pms.BasicInfo.Instrument,
		pms.Difficulty,
	}
	writer.Write(s)
	writer.Flush()
}
