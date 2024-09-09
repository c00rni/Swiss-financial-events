package main

import (
	"context"
	"fmt"
	"github.com/c00rni/Swiss-financial-events/internal/database"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
)

type swissCfaEvent struct {
	Date time.Time
}

func (cfg *apiConfig) scrapeCfasocietyTopics() {
	c := colly.NewCollector()
	c.AllowURLRevisit = false
	domain := "https://cfasocietyswitzerland.org"
	topics := map[string]string{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML(".dropdown-menu", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, h *colly.HTMLElement) {
			if link := h.Attr("href"); strings.Contains(link, "topic=") {
				h.Request.Visit(h.Attr("href"))
				id := uuid.New().String()
				linkSplit := strings.Split(link, "=")
				topicName := linkSplit[1]
				topic := database.AddTopicParams{
					ID:   id,
					Name: topicName,
				}
				_, err := cfg.DB.AddTopic(context.Background(), topic)
				if err == nil {
					log.Printf("Topic '%v' has been added.", topicName)
				}
				cat, err := cfg.DB.GetTopicByName(context.Background(), topicName)
				if err != nil {
					log.Println(err)
				}
				topics[topicName] = cat.ID
			}
		})
	})

	c.OnHTML(".events__teaser__link", func(e *colly.HTMLElement) {
		fullURL := e.Request.URL.String()
		if strings.Contains(fullURL, "topic=") {
			linkSplit := strings.Split(fullURL, "=")
			topicName := linkSplit[1]
			event, err := cfg.DB.GetEventsByLink(context.Background(), domain+e.Attr("href"))
			if err != nil {
				log.Println(err)
				return
			}

			topicId, ok := topics[topicName]
			if !ok {
				log.Printf("Topics '%v' not found", topicId)
				return
			}

			params := database.LinkEventToTopicParams{
				EventID: event.ID,
				TopicID: topicId,
			}
			cfg.DB.LinkEventToTopic(context.Background(), params)
		}
	})

	c.Visit(domain + "/events/event-calendar/")
}

func (cfg *apiConfig) scrapeCfasocietyCategories() {
	c := colly.NewCollector()
	c.AllowURLRevisit = false
	domain := "https://cfasocietyswitzerland.org"
	categories := map[string]string{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnHTML(".dropdown-menu", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, h *colly.HTMLElement) {
			if link := h.Attr("href"); strings.Contains(link, "category=") {
				h.Request.Visit(h.Attr("href"))
				id := uuid.New().String()
				linkSplit := strings.Split(link, "=")
				categoryName := linkSplit[1]
				category := database.AddCategoryParams{
					ID:   id,
					Name: categoryName,
				}
				_, err := cfg.DB.AddCategory(context.Background(), category)
				if err == nil {
					log.Printf("Category '%v' has been added.", categoryName)
				}
				cat, err := cfg.DB.GetCategoryByName(context.Background(), categoryName)
				if err != nil {
					log.Println(err)
				}
				categories[categoryName] = cat.ID
			}
		})
	})

	c.OnHTML(".events__teaser__link", func(e *colly.HTMLElement) {
		fullURL := e.Request.URL.String()
		if strings.Contains(fullURL, "category=") {
			linkSplit := strings.Split(fullURL, "=")
			categoryName := linkSplit[1]
			event, err := cfg.DB.GetEventsByLink(context.Background(), domain+e.Attr("href"))
			if err != nil {
				log.Println(err)
				return
			}

			categoryId, ok := categories[categoryName]
			if !ok {
				log.Printf("Category '%v' not found", categoryId)
				return
			}

			params := database.LinkEventToCategoryParams{
				EventID:    event.ID,
				CategoryID: categoryId,
			}
			cfg.DB.LinkEventToCategory(context.Background(), params)
		}
	})

	c.Visit(domain + "/events/event-calendar/")
}

func (cfg *apiConfig) scrapeCfasocietyEvents() {
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

		cfg.DB.CreateEvent(context.Background(), events[fullURL])
	})

	c.Visit(domain + "/events/event-calendar/")
}
