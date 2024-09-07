package main

import (
	"context"
	"fmt"
	"github.com/c00rni/Swiss-financial-events/internal/database"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"log"
	"time"
)

type swissCfaEvent struct {
	Date time.Time
}

func (cfg *apiConfig) scrapeCfasociety() {
	c := colly.NewCollector()
	domain := "https://cfasocietyswitzerland.org"
	events := map[string]database.CreateEventParams{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML(".events__teaser__link", func(e *colly.HTMLElement) {
		id := uuid.New().String()
		relativeLink := e.Attr("href")
		location := e.ChildText(".events__teaser__location")
		fullLink := domain + relativeLink

		event := database.CreateEventParams{
			ID:       id,
			Link:     fullLink,
			Location: location,
		}
		events[fullLink] = event
		e.Request.Visit(relativeLink)
	})

	c.OnHTML(".events-app", func(e *colly.HTMLElement) {
		startDate, endDate, err := cfasocietyDateFormater(e.ChildText(".event_details__info__header-date"))
		if err != nil {
			log.Println("Error:", err)
		}
		fullURL := e.Request.URL.String()
		event := events[fullURL]
		event.Title = e.ChildText(".event_details__title .cfa-title")
		event.Address = e.ChildText(".event_details__info__header-location a")
		event.StartAt = startDate
		event.EndAt = endDate
		event.Description = e.ChildText(".container > p")
		events[fullURL] = event

		_, err = cfg.DB.CreateEvent(context.Background(), events[fullURL])
		if err != nil {
			log.Println("Error:", err)
		}
	})

	c.Visit(domain + "/events/event-calendar/")
}
