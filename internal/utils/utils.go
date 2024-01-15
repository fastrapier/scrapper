package utils

import (
	"encoding/json"
	"log"
	"os"
	"scrapper/internal/parser"
	"time"
)

func WriteToJson(configurators parser.ConfiguratorsMap) {
	jsonData, err := json.MarshalIndent(configurators, "", "\t")

	if err != nil {
		log.Fatalf("Err:%s", err)
	}

	now := time.Now().Format(time.DateTime)

	if _, err := os.Stat("./tmp"); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("tmp", 0777)

			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	jsonFile, err := os.Create("./tmp/configurators-" + now + ".json")

	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
}
