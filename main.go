package main

import (
	"encoding/json"
	"fmt"

	Piascore "./Piascore"
	PrintMusicalScore "./PrintMusicScore"
)

func main() {
	returnMap := make(map[string]interface{})
	returnMap["PrintMusicScore"] = PrintMusicalScore.GetInfo("tenorSaxophone", "ルパン")
	returnMap["Piascore"] = Piascore.GetInfo("tenorSaxophone", "Pretender")
	json, _ := json.Marshal(returnMap)
	fmt.Println(string(json))
}
