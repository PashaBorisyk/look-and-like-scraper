package main

import (
	"log"
	"look-and-like-scraper/logger"
	"look-and-like-scraper/web"
	"sync"
	"time"
)

var zaraScraper *web.ZaraScraper
var hmScraper *web.HMScraper

func init() {
	logger.Init()
	initScrappers()
}

func initScrappers() {

	hmScraper = web.NewHMScraper()
	hmScraper.Init()

	zaraScraper = web.NewZaraScrapper()
	zaraScraper.Init()
}

func main() {
	startScrapping()
}

func startScrapping() {

	log.Println("Starting scraping at: ", time.Now())

	start := time.Now()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		hmScraper.Scrap()
		hmTime := time.Since(start)
		log.Println("HM scrap took: ", hmTime.Seconds(), " seconds")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		zaraStart := time.Now()
		zaraScraper.Scrap()
		zaraTime := time.Since(zaraStart)
		log.Println("Zara scrap took: ", zaraTime.Seconds(), " seconds")

	}()

	wg.Wait()
	programTime := time.Since(start)
	log.Println("Full scrap took: ", programTime, " seconds")
}
