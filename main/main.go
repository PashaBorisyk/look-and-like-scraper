package main

import (
	"log"
	"look-and-like-web-scrapper/logger"
	"look-and-like-web-scrapper/web"
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

	hmScraper.Scrap()
	hmTime := time.Since(start)
	log.Println("HM scrap took: ", hmTime.Seconds(), " seconds")

	zaraStart := time.Now()
	zaraScraper.Scrap()
	zaraTime := time.Since(zaraStart)
	log.Println("Zara scrap took: ", zaraTime.Seconds(), " seconds")

	programTime := time.Since(start)
	log.Println("Full scrap took: ", programTime, " seconds")
}
