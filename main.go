package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"strconv"
	"strings"
)

type ConfiguratorProduct struct {
	Id    int
	Title string
	Price int
}
type Configurator struct {
	Title                string
	Link                 string
	Price                int
	ConfiguratorProducts map[string][]ConfiguratorProduct
}

func main() {

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"

	configuratorProducts := c.Clone()

	var configurators = collectConfigurators(*c)

	for i, configurator := range configurators {
		fmt.Printf("Index: %d | %+v\n", i, configurator)
	}

	if len(configurators) == 0 {
		log.Fatal("Empty configurators!")
		return
	}

	configurators = collectConfiguratorProducts(*configuratorProducts, configurators)
}

func collectConfiguratorProducts(collector colly.Collector, configurators []Configurator) []Configurator {

	fmt.Println("Start collecting configurator products!")

	collector.OnHTML("table.product-specifications tbody tr", func(element *colly.HTMLElement) {
		specification := strings.TrimSpace(element.ChildText("td:eq(0)"))
		// Todo: Нужно продумать логику извлечения данных из столбцов
		fmt.Printf("Start parsing specification %s\n", specification)

		var configuratorProducts []ConfiguratorProduct

		element.ForEach("td:eq(1) select option", func(i int, element *colly.HTMLElement) {

			id, err := strconv.Atoi(element.Attr("value"))

			if err != nil {
				fmt.Printf("Link:%s | err:%s\n", element.Request.URL.String(), err.Error())
			}

			title := element.Text

			price, err := strconv.Atoi(element.Attr("data-price"))

			if err != nil {
				fmt.Printf("Link:%s | err:%s\n", element.Request.URL.String(), err.Error())
			}

			configuratorProduct := ConfiguratorProduct{
				Id:    id,
				Title: title,
				Price: price,
			}
			fmt.Printf("configurator product: %+v\n", configuratorProduct)
			configuratorProducts = append(configuratorProducts, configuratorProduct)
		})

		for _, configurator := range configurators {
			configurator.ConfiguratorProducts[specification] = configuratorProducts
		}
	})

	collector.OnError(func(response *colly.Response, err error) {
		fmt.Printf("Link:%s -> err:%s", response.Request.URL.String(), err.Error())
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Printf("Visiting:%s\n", request.URL.String())
	})

	for i, configurator := range configurators {
		fmt.Printf("Start parsing #%d | %+v\n", i, configurator)
		collector.Visit(configurator.Link)
	}

	return configurators
}

func collectConfigurators(collector colly.Collector) []Configurator {
	var configurators []Configurator

	fmt.Println("Start collecting configurators!")

	collector.OnHTML("div.products figure.product-item", func(element *colly.HTMLElement) {
		title := strings.TrimSpace(element.ChildAttr("figcaption.info h3 a", "title"))

		link := element.ChildAttr("figcaption.info h3 a", "href")

		priceText := element.ChildText("div.prices")

		pricesArr := strings.Split(priceText, "₽")

		priceText = pricesArr[0]
		priceText = strings.Replace(priceText, " ", "", -1)
		priceText = strings.Replace(priceText, "₽", "", -1)

		price, err := strconv.Atoi(priceText)

		if err != nil {
			fmt.Printf("Price not found - %s\n", title)
		}

		configurator := Configurator{
			Title: title,
			Link:  link,
			Price: price,
		}

		configurators = append(configurators, configurator)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.OnError(func(response *colly.Response, err error) {
		log.Fatalf("%s", err.Error())
	})

	collector.Visit("https://sale-server.ru/konfigurator/products")

	return configurators
}
