package web

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/zolamk/colly-mongo-storage/colly/mongo"
	"look-and-like-web-scrapper/config"
	"look-and-like-web-scrapper/models"
	"strings"
	"time"
)

const UserAgent = "look-and-like-scrapper"

type Scraper struct {
	Locales       []config.Locale
	CurrentLocale config.Locale
	Category      string
	Sex           string
}

var scrappingTime = time.Now().Format(time.RFC3339)
var storage *mongo.Storage

func init() {
	configureStorage()
}

func configureStorage() {
	storage = &mongo.Storage{
		Database: config.GetConfig().MongoConfig.CollyDatabaseName,
		URI:      config.GetConfig().MongoConfig.Uri,
	}
}

func (scraper Scraper) isGenderTranslation(value string) bool {
	return strings.Contains(value, scraper.CurrentLocale.MaleTranslation) ||
		strings.Contains(value, scraper.CurrentLocale.FemaleTranslation) ||
		strings.Contains(value, scraper.CurrentLocale.KidsTranslation)
}

func foreachSelection(selection *goquery.Selection, foreach func(int, *goquery.Selection)) {
	currentElement := selection.First()
	for i := 0; i < selection.Size(); i++ {
		foreach(i, currentElement)
		currentElement = currentElement.Next()
	}
}

func normalizeString(s string) string {
	return strings.ToLower(strings.Trim(s, " \n\t"))
}

func uniqueComposition(intSlice []models.Composition) []models.Composition {
	keys := make(map[models.Composition]bool)
	var list []models.Composition
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
