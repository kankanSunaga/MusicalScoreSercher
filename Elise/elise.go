package Elise

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

func Main() {
	fmt.Println("start")
	insts := Common.Instruments()
	for _, installment := range insts {
		GetInfo(installment)
	}
}

func GetInfo(instrument string) {
	eliseClient := initEliseClient(instrument)
	doc := eliseClient.Get()
	musicTable := getEliseTable(doc)
	getTableInfo(musicTable, instrument)
}

func initEliseClient(instrument string) *Common.ApiClientBase {
	elc := eliseClient{}
	var acb *Common.ApiClientBase
	acb = &Common.ApiClientBase{
		ServiceName: elise,
	}
	elc.ApiClientBase = *acb
	elc.setUrl(instrument)
	return &elc.ApiClientBase

}

func (eli *eliseClient) setUrl(instrument string) {
	var url string
	switch Common.WhichInstrumentType(instrument) {
	case "Saxophone":
		url = setSaxUrl(instrument)
	default:
		url = eliseUrlBase
	}
	eli.ApiClientBase.Url = url
}

func setSaxUrl(instrument string) string {
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
	eli := Elise{}
	var cbi *Common.BasicInfo
	cbi = &Common.BasicInfo{
		ServiceName: elise,
		Instrument:  instrument,
	}
	eli.BasicInfo = cbi
	return &eli
}

func getTableInfo(gd *goquery.Selection, instrument string) {
	gd.Find(".s_result_content").EachWithBreak(func(i int, tr *goquery.Selection) bool {
		eli := initElise(instrument)
		eli.setBasicInfo(tr)
		eli.setEliseInfo(tr)
		eli.Output()
		return true
	})
	return
}

func (eli *Elise) setBasicInfo(gs *goquery.Selection) *Elise {
	path, _ := gs.Find(".title").Find("a").Attr("href")
	url := eliseDomain + path
	eli.BasicInfo.Url = url
	eli.BasicInfo.Id = getId(path)
	eli.BasicInfo.MusicName = gs.Find(".title").Text()
	eli.BasicInfo.Composer = gs.Find(".artist").Text()
	eli.BasicInfo.Price = Common.GetPrice(gs.Find(".costcv").Text())
	return eli
}

func (pms *Elise) setEliseInfo(gs *goquery.Selection) *Elise {
	pms.Difficulty = Common.RemoveBlankStrings(gs.Find(".inst").Text())
	return pms
}

func (eli *Elise) Output() {
	file, err := os.OpenFile("test.csv", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		//エラー処理
		log.Fatal(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	var s = []string{
		eli.BasicInfo.ServiceName,
		eli.BasicInfo.MusicName,
		eli.BasicInfo.Composer,
		strconv.Itoa(eli.BasicInfo.Price),
		eli.BasicInfo.Url,
		eli.BasicInfo.Instrument,
		eli.Difficulty,
		eli.BasicInfo.Id,
	}
	writer.Write(s)
	writer.Flush()
}

func getId(path string) string {
	path = strings.TrimRight(path, "/")
	slice := strings.Split(path, "/")
	return slice[len(slice)-1]
}
