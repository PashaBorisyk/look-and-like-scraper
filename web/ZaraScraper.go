package web

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"look-and-like-web-scrapper/config"
	"look-and-like-web-scrapper/db"
	"look-and-like-web-scrapper/models"
	"look-and-like-web-scrapper/queue"
	"strconv"
	"strings"
)

const zaraDomain = "www.zara.com"
const zaraShopName = "zara"

type ZaraScraper struct {
	Scraper
	mainPageCollector     *colly.Collector
	productsListCollector *colly.Collector
	productPageCollector  *colly.Collector
}

func NewZaraScrapper() *ZaraScraper {
	scrapper := ZaraScraper{}
	scrapper.Locales = config.GetConfig().ScraperConfig.Locales[zaraShopName]
	return &scrapper
}

func (scraper *ZaraScraper) Init() {
	log.Println("ZaraScraper Init")

	scraper.configureMainPageCollector()
	scraper.configureProductsListCollector()
	scraper.configureProductPageCollector()
}

func (scraper *ZaraScraper) Scrap() {
	log.Println("Starting ZaraScraper")

	for _, locale := range scraper.Locales {
		scraper.CurrentLocale = locale
		err := scraper.mainPageCollector.Visit(locale.BaseURL)
		if err != nil {
			log.Println("Error opening site: ", err)
		}
	}

}

func (scraper *ZaraScraper) configureMainPageCollector() {

	scraper.mainPageCollector = colly.NewCollector(
		colly.AllowedDomains(zaraDomain),
		colly.UserAgent(UserAgent),
		colly.AllowURLRevisit(),
	)

	scraper.mainPageCollector.OnHTML("a._category-link[href] ", func(e *colly.HTMLElement) {
		productListLink := e.Attr("href")
		//log.Println("Extracted category URL: ",productListLink)
		value := strings.ToLower(e.Text)
		if scraper.isGenderTranslation(value) {
			scraper.Sex = value
		} else {
			scraper.Category = value
		}
		err := scraper.productsListCollector.Visit(productListLink)
		if err != nil {
			log.Println("Error while visiting Product list page")
		}
	})

}

func (scraper *ZaraScraper) configureProductsListCollector() {

	scraper.productsListCollector = colly.NewCollector(
		colly.AllowedDomains(zaraDomain),
		colly.UserAgent(UserAgent),
		colly.AllowURLRevisit(),
	)

	scraper.productsListCollector.OnHTML("a._item.item[href]", func(element *colly.HTMLElement) {
		//log.Println("Found Product: ", element)
		productLink := element.Attr("href")
		err := scraper.productPageCollector.Visit(productLink)
		if err != nil {
			log.Println("Error while visiting Product page: ", err)
		}
	})
}

func (scraper *ZaraScraper) configureProductPageCollector() {

	productsCollection := db.GetCollection("products")

	scraper.productPageCollector = colly.NewCollector(
		colly.AllowedDomains(zaraDomain),
		colly.UserAgent(UserAgent),
	)
	if err := scraper.productPageCollector.SetStorage(storage); err != nil {
		panic(err)
	}

	scraper.productPageCollector.OnHTML("section.content-main[id=main]", func(element *colly.HTMLElement) {

		product := scraper.createProduct(element)
		insertedKey, err := productsCollection.Insert(*product)
		if err != nil {
			log.Println("Error while inserting Zara product in database: ", err)
		} else {
			log.Println("Zara product with url:' ", product.MetaInformation.Url, "' inserted; Publishing key")
			queue.PublishKey(insertedKey)
		}

	})

}

func (scraper *ZaraScraper) createProduct(element *colly.HTMLElement) *models.Product {

	productUrl := scraper.getProductUrl(element)
	productName := scraper.getProductName(element)
	productColor, productArticle := scraper.getColorAndArticle(element)
	productPriceValue, productPriceCurrency := scraper.getProductPrice(element)
	productDescription := scraper.getProductDescription(element)
	productImages := scraper.getProductImages(element)
	productSizes := scraper.getProductSizes(element)

	metaInformation := models.MetaInformation{
		Url:        productUrl,
		InsertDate: scrappingTime,
		ShopName:   zaraShopName,
		BaseURL:    scraper.CurrentLocale.BaseURL,
		Alpha3Code: scraper.CurrentLocale.Alpha3Code,
		LocaleLCID: scraper.CurrentLocale.LocaleLCID,
		Domain:     zaraDomain,
	}

	price := models.Price{
		Value:    productPriceValue,
		Currency: productPriceCurrency,
	}

	images := models.Images{
		StockImageUrls: productImages,
	}

	data := models.Data{
		Images:      images,
		Price:       price,
		Sex:         scraper.Sex,
		Category:    scraper.Category,
		Color:       productColor,
		Name:        productName,
		Article:     productArticle,
		Sizes:       productSizes,
		Description: normalizeString(productDescription),
	}

	return &models.Product{
		Data:            data,
		MetaInformation: metaInformation,
	}
}

func (scraper ZaraScraper) getProductUrl(element *colly.HTMLElement) string {
	return element.Request.URL.String()
}

func (scraper ZaraScraper) getProductName(element *colly.HTMLElement) string {
	name := element.DOM.Find("h1.product-name").Before("span").Text()
	return normalizeString(name)
}

func (scraper ZaraScraper) getProductDescription(element *colly.HTMLElement) string {
	return normalizeString(element.DOM.Find("p.description").Text())
}

func (scraper ZaraScraper) getColorAndArticle(element *colly.HTMLElement) (color string, article string) {
	color = element.DOM.Find("span._colorName").Text()
	article = element.DOM.Find("span[data-qa-qualifier]").Text()
	return normalizeString(color), article
}

func (scraper ZaraScraper) getProductPrice(element *colly.HTMLElement) (price float64, priceCurrency string) {

	priceJson := element.DOM.Find("script").Text()
	var jsonRep []models.ProductRep
	err := json.Unmarshal([]byte(priceJson), &jsonRep)

	if err != nil {
		log.Println("Error while unmarshal prices json: ", err)
		return 0, ""
	}
	if len(jsonRep) == 0 {
		return 0, ""
	}

	offer := jsonRep[0].Offers
	price, err = strconv.ParseFloat(offer.Price, 32)
	if err != nil {
		log.Println("Error parsing price: ", err)
		return 0, ""
	}
	return price, offer.PriceCurrency
}

func (scraper ZaraScraper) getProductImages(element *colly.HTMLElement) (imageURLs []string) {
	imageContainers := element.DOM.Find("div.media-wrap.image-wrap")
	imageURLs = make([]string, imageContainers.Size())
	foreachSelection(imageContainers, func(i int, imageContainer *goquery.Selection) {
		url, isPresent := imageContainer.Find("a[href]").Attr("href")
		if isPresent {
			imageURLs[i] = url
		}
	})

	return imageURLs
}

func (scraper ZaraScraper) getProductSizes(element *colly.HTMLElement) (sizes []string) {
	sizesContainers := element.DOM.Find("span.size-name")
	sizes = make([]string, sizesContainers.Size())
	foreachSelection(sizesContainers, func(i int, sizesContainer *goquery.Selection) {
		sizes[i] = normalizeString(sizesContainer.Text())
	})

	return sizes
}
