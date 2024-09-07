package main

import (
	"fmt"
	"github.com/gocolly/colly"
	//"io"
	"log"
	//"net/http"
	"time"
)

type swissCfaEvent struct {
	Date time.Time
}

func scrapeCfasociety() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML(".events__teaser__link", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
		log.Println("Location:", e.ChildText(".events__teaser__location"))
	})

	c.OnHTML(".events-app", func(e *colly.HTMLElement) {
		// event := swissCfaEvent{}
		startDate, endDate, err := cfasocietyDateFormater(e.ChildText(".event_details__info__header-date"))
		if err != nil {
			log.Println("Error: ", err)
		}
		log.Println("Title:", e.ChildText(".event_details__title .cfa-title"))
		log.Println("Address:", e.ChildText(".event_details__info__header-location a"))
		log.Println(startDate, "/", endDate)
		log.Println("#######")
	})

	c.Visit("https://cfasocietyswitzerland.org/events/event-calendar/")
}
