package parser

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type ConfiguratorProduct struct {
	ID       int32
	Title    string
	Price    int32
	PriceEur float64
}

type Configuration struct {
	Name                 string
	ConfiguratorProducts []ConfiguratorProduct
}

type Configurator struct {
	Title          string
	Price          int32
	PriceEur       float64
	Link           string
	Configurations []Configuration
}

type ConfiguratorsMap map[string]*Configurator

func ParseConfigurator(euro float64) ConfiguratorsMap {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"

	c.OnError(func(response *colly.Response, err error) {
		log.Fatalf("Link: %s -> %s", response.Request.URL.String(), err)
	})

	c.OnRequest(func(request *colly.Request) {
		log.Printf("Visiting: %s", request.URL.String())
	})

	configurators := make(ConfiguratorsMap)

	configurators.fill(*c, euro)

	return configurators
}

func (configurators ConfiguratorsMap) fill(c colly.Collector, euro float64) {

	c.OnHTML("div.products figure.product-item", func(element *colly.HTMLElement) {

		title := strings.TrimSpace(element.ChildAttr("figcaption.info h3 a", "title"))

		link := element.ChildAttr("figcaption.info h3 a", "href")

		priceText := element.ChildText("div.prices")

		pricesArr := strings.Split(priceText, "₽")

		priceText = pricesArr[0]
		priceText = strings.Replace(priceText, " ", "", -1)
		priceText = strings.Replace(priceText, "₽", "", -1)

		price, err := strconv.Atoi(priceText)

		if err != nil {
			fmt.Printf("Price not found - %s %s\n", title, link)
		}

		configurators[link] = &Configurator{
			Title:          title,
			Price:          int32(price),
			PriceEur:       float64(price) / euro,
			Link:           link,
			Configurations: make([]Configuration, 0),
		}

		c.Visit(link)
	})

	c.OnHTML("table.product-specifications > tbody > tr", func(element *colly.HTMLElement) {
		html, err := element.DOM.Find("td:nth-child(1)").Html()

		if err != nil {
			log.Fatalf("%s", err)
			return
		}

		configuration := Configuration{
			Name:                 strings.TrimSpace(html),
			ConfiguratorProducts: make([]ConfiguratorProduct, 0),
		}

		element.ForEach("td:nth-of-type(2) select option", func(i int, element *colly.HTMLElement) {

			if element.Attr("value") == "0" || strings.Contains(element.Text, "Нет в наличии") {
				return
			}

			id, err := strconv.Atoi(element.Attr("value"))

			if err != nil {
				log.Printf("ID - Link:%s | err:%s\n", element.Request.URL.String(), err)
			}

			title := strings.TrimSpace(element.Text)

			priceRegex, err := regexp.Compile("-\\s+\\d+.*₽")

			if err != nil {
				log.Fatalf("%s", err)
			}

			title = strings.TrimSpace(priceRegex.ReplaceAllString(title, ""))

			price, err := strconv.Atoi(element.Attr("data-price"))

			if err != nil {
				log.Printf("PRICE - Link:%s | err:%s\n", element.Request.URL.String(), err)
			}

			configuratorProduct := ConfiguratorProduct{
				ID:       int32(id),
				Title:    title,
				Price:    int32(price),
				PriceEur: float64(price) / euro,
			}

			configuration.ConfiguratorProducts = append(configuration.ConfiguratorProducts, configuratorProduct)
		})
		configurators[element.Request.URL.String()].Configurations = append(configurators[element.Request.URL.String()].Configurations, configuration)
	})
	c.Visit("https://sale-server.ru/konfigurator/products")
}
