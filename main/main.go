package main

import (
	"log"
	"look-and-like-web-scrapper/web"
	"time"
)

var zaraScraper *web.ZaraScraper
var hmScraper *web.HMScraper

func init() {
	initScrappers()
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

func initScrappers() {

	hmLocales := [] web.Locale{
		{
			BaseURL:           "https://www2.hm.com/ru_ru/index.html",
			Alpha3Code:        "RUS",
			LocaleLCID:        "ru",
			FemaleTranslation: "женщ",
			MaleTranslation:   "муж",
			KidsTranslation:   "дети",
		},
	}
	hmScraper = web.NewHMScraper(hmLocales)
	hmScraper.Init()

	zaraLocales := [] web.Locale{
		{
			BaseURL:           "https://www.zara.com/by/ru/",
			Alpha3Code:        "BLR",
			LocaleLCID:        "be",
			FemaleTranslation: "женщ",
			MaleTranslation:   "муж",
			KidsTranslation:   "дети",
		},
	}

	zaraScraper = web.NewZaraScrapper(zaraLocales)
	zaraScraper.Init()
}
