package main

import (
	"fmt"
	"time"

	//"./Elise"
	"./Piascore"
	//PrintMusicalScore "./PrintMusicScore"
)

func main() {
	start := time.Now()

	//returnMap := make(map[string]interface{})
	//returnMap["PrintMusicScore"] = PrintMusicalScore.GetInfo("tenorSaxophone", "ルパン")
	//returnMap["Piascore"] = Piascore.GetInfo("tenorSaxophone", "Pretender")
	//returnMap["Elise"] = Elise.GetInfo("tenorSaxophone", "ルパン")
	//json, _ := json.Marshal(returnMap)
	//fmt.Println(string(json))
	//PrintMusicalScore.Main()
	//Elise.Main()
	Piascore.Main()
	end := time.Now()
	fmt.Println("%f秒\n", (end.Sub(start)).Seconds())
}
