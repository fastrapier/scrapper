package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"scrapper/internal/parser"
	"scrapper/internal/utils"
)

func main() {
	cfg, err := ini.Load("./config/config.ini")

	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	euro, err := cfg.Section("").Key("euro").Float64()

	if err != nil {
		fmt.Printf("Fail to read euro key: %v", err)
		os.Exit(1)
	}

	configuratorsMap := parser.ParseConfigurator(euro)

	utils.WriteToJson(configuratorsMap)
}
