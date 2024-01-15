package main

import (
	"scrapper/internal/parser"
	"scrapper/internal/utils"
)

func main() {
	euro := 105.12

	configuratorsMap := parser.ParseConfigurator(euro)
	utils.WriteToJson(configuratorsMap)
}
