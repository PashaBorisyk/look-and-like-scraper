package web

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/zolamk/colly-mongo-storage/colly/mongo"
	"look-and-like-web-scrapper/models"
	"strings"
)

const UserAgent = "look-and-like-scrapper"

var storage *mongo.Storage

type Locale struct {
	LocaleLCID        string
	Alpha3Code        string
	BaseURL           string
	MaleTranslation   string
	FemaleTranslation string
	KidsTranslation   string
}

type Scraper struct {
	Locales       []Locale
	CurrentLocale Locale
	Category      string
	Sex           string
}

func init() {
	configureStorage()
}

func configureStorage() {
	storage = &mongo.Storage{
		Database: "colly",
		URI:      "mongodb://127.0.0.1:27017",
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
