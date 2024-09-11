package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(colly.Async())
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 5})

	c.OnHTML("article p, article h1, article h2, article h3, article h4, article h5, article h6", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})

	c.Visit("https://www.engadget.com/2019/08/25/sony-and-yamaha-sc-1-sociable-cart/")
	c.Wait()
}
