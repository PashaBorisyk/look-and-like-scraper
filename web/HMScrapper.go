package web

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/mantyr/pricer"
	"log"
	"look-and-like-web-scrapper/db"
	"look-and-like-web-scrapper/models"
	"regexp"
	"strings"
	"time"
)

const HMDomain = "www2.hm.com"
const HMShopName = "H&M"

const filterPagination = "?sort=stock&image-size=small&image=model&offset=0&page-size=10000"
const imagePrefixLen = len("'image': isDesktop ? '")

type HMScraper struct {
	Scraper
	mainPageCollector     *colly.Collector
	productsListCollector *colly.Collector
	productPageCollector  *colly.Collector
}

var imageRegex *regexp.Regexp

func init() {
	var err error
	imageRegex, err = regexp.Compile("'image': isDesktop \\? '(.+?)'")
	if err != nil {
		log.Println("Can not compile regex for HM images: ")
		panic(err)
	}

}

func NewHMScraper(locales []Locale) *HMScraper {
	scrapper := HMScraper{}
	scrapper.Locales = locales
	return &scrapper
}

func (scraper *HMScraper) Init() {
	log.Println("HMScraper Init")

	scraper.configureMainPageCollector()
	scraper.configureProductListCollector()
	scraper.configureProductPageCollector()
}

func (scraper *HMScraper) Scrap() {
	log.Println("Starting HM scraper")

	for _, locale := range scraper.Locales {
		scraper.CurrentLocale = locale
		err := scraper.mainPageCollector.Visit(locale.BaseURL)
		if err != nil {
			log.Println("Error opening site: ", err)
		}
	}

}

func (scraper *HMScraper) configureMainPageCollector() {
	log.Println("Configuring MainPageCollector")

	scraper.mainPageCollector = colly.NewCollector(
		colly.AllowedDomains(HMDomain),
		colly.UserAgent(UserAgent),
		colly.AllowURLRevisit(),
	)

	scraper.mainPageCollector.OnHTML("a.menu__sub-link[href]", func(element *colly.HTMLElement) {
		productListLink := element.Attr("href")
		scraper.Category = normalizeString(element.Text)
		sex := normalizeString(scraper.getGender(element))
		if scraper.isGenderTranslation(sex) {
			scraper.Sex = sex
		} else {
			scraper.Category = sex
		}
		link := element.Request.AbsoluteURL(productListLink)
		log.Println("Extracted category URL: ", link)
		linkWithFilter := getLinkWithFilter(link)
		err := scraper.productsListCollector.Visit(linkWithFilter)
		if err != nil {
			log.Println("Will not visit product list page: ", err)
		}
	})

}

func (scraper *HMScraper) getGender(element *colly.HTMLElement) string {
	sex := element.DOM.Parent().Parent().Parent().Parent().Find("button.menu__title-button").First()
	return normalizeString(sex.Text())
}

func (scraper *HMScraper) configureProductListCollector() {
	log.Println("Configuring ProductListCollector")

	scraper.productsListCollector = colly.NewCollector(
		colly.AllowedDomains(HMDomain),
		colly.UserAgent(UserAgent),
		colly.AllowURLRevisit(),
	)
	scraper.productsListCollector.SetRequestTimeout(100 * time.Second)

	scraper.productsListCollector.OnHTML("a.item-link[href]", func(element *colly.HTMLElement) {
		productLink := element.Request.AbsoluteURL(element.Attr("href"))
		err := scraper.productPageCollector.Visit(productLink)
		if err != nil {
			log.Println("Will not go to the product page: ", err)
		}
	})
}

func (scraper *HMScraper) configureProductPageCollector() {
	log.Println("Configuring ProductPageCollector")

	productsCollection := db.GetCollection("products")

	scraper.productPageCollector = colly.NewCollector(
		colly.AllowedDomains(HMDomain),
		colly.UserAgent(UserAgent),
	)

	if err := scraper.productPageCollector.SetStorage(storage); err != nil {
		panic(err)
	}

	scraper.productPageCollector.OnHTML("div.module.product-description.sticky-wrapper", func(element *colly.HTMLElement) {
		product := scraper.createProduct(element)
		err := productsCollection.Insert(product)
		if err != nil {
			log.Println("Error inserting H&M product in database: ", err)
		} else {
			log.Println("H&M product with url:' ", product.Url, "' inserted")
		}
	})

}

func (scraper *HMScraper) createProduct(element *colly.HTMLElement) *models.Product {

	productName := scraper.getProductName(element)
	productPrice, productPriceCurrency := scraper.getProductPriceAndCurrency(element)
	productColor, productArticle := scraper.getProductColorAndArticle(element)
	productImages := scraper.getProductImages(element)
	productDescription := scraper.getProductDescription(element)
	productComposition := scraper.getProductComposition(element)
	productSizes := scraper.getSizes(element)

	return &models.Product{
		Name:        normalizeString(productName),
		Price:       productPrice,
		Currency:    productPriceCurrency,
		Color:       normalizeString(productColor),
		Article:     normalizeString(productArticle),
		Images:      productImages,
		Description: normalizeString(productDescription),
		Composition: productComposition,
		InsertDate:  time.Now(),
		Sizes:       productSizes,
		Category:    normalizeString(scraper.Category),
		Sex:         scraper.Sex,
		Domain:      HMDomain,
		ShopName:    HMShopName,
		LocaleLCID:  scraper.CurrentLocale.LocaleLCID,
		Alpha3Code:  scraper.CurrentLocale.Alpha3Code,
		BaseURL:     scraper.CurrentLocale.BaseURL,
		Url:         element.Request.URL.String(),
	}
}

func (scraper *HMScraper) getProductName(element *colly.HTMLElement) string {
	return strings.Trim(element.DOM.Find("h1.primary.product-item-headline").Text(), "\t\n ")
}

func (scraper *HMScraper) getProductPriceAndCurrency(element *colly.HTMLElement) (float64, string) {
	priceToken := strings.Trim(element.DOM.Find("span.price-value").Text(), "\n ")
	price := pricer.NewPrice()
	price.Parse(priceToken)
	return price.GetFloat64(), price.GetType()
}

func (scraper *HMScraper) getProductColorAndArticle(element *colly.HTMLElement) (color string, article string) {
	colorAndArticleElement := element.DOM.Find("a.filter-option.miniature.active[title]")
	color, _ = colorAndArticleElement.Attr("data-color")
	article, _ = colorAndArticleElement.Attr("data-articlecode")
	return color, article
}

func (scraper *HMScraper) getProductImages(element *colly.HTMLElement) (imageURLs []string) {

	parsedImagesTokens := imageRegex.FindAll(element.Response.Body, 5)
	imageURLs = make([]string, len(parsedImagesTokens))
	for index, url := range parsedImagesTokens {
		imageURLs[index] = string(url[imagePrefixLen : len(url)-1])
	}
	return imageURLs
}

func (scraper *HMScraper) getProductDescription(element *colly.HTMLElement) (description string) {
	description = element.DOM.Find("p.pdp-description-text").Text()
	return description
}

func (scraper *HMScraper) getProductComposition(element *colly.HTMLElement) []models.Composition {
	composition := element.DOM.Find("li.article-composition.pdp-description-list-item").Find("li")
	var compositions []models.Composition
	foreachSelection(composition, func(i int, compositionItem *goquery.Selection) {

		compositionText := compositionItem.Text()
		partNameAndComposition := strings.Split(compositionText, ":")

		if len(partNameAndComposition) == 2 {
			part := partNameAndComposition[0]
			compositionsParts := strings.Split(partNameAndComposition[1], ",")
			compositions = append(compositions, scraper.getCompositions(part, compositionsParts, compositions)...)
		} else {
			compositionsParts := strings.Split(partNameAndComposition[0], ",")
			compositions = append(compositions, scraper.getCompositions("", compositionsParts, compositions)...)
		}

	})
	return uniqueComposition(compositions)
}

func (scraper *HMScraper) getCompositions(part string, compositionParts []string, compositions []models.Composition) []models.Composition {

	for _, compositionPart := range compositionParts {
		compositionWithPercent := strings.Split(strings.Trim(compositionPart, " \t\n"), " ")
		if len(compositionWithPercent) == 2 {
			material := compositionWithPercent[0]
			percent := compositionWithPercent[1]
			compositions = append(compositions, models.Composition{
				Part:     part,
				Material: material,
				Percent:  percent,
			})
		}
	}
	return compositions
}

func (scraper *HMScraper) getSizes(element *colly.HTMLElement) []string {
	sizesJson := scraper.getSizesJson(element.Response.Body)
	var sizes []models.Size
	err := json.Unmarshal(sizesJson, &sizes)
	if err != nil {
		log.Println("Error while unmarshal HM sizes: ", err)
		return nil
	}
	return sizeToStringArray(sizes)
}

func (scraper *HMScraper) getSizesJson(body []byte) []byte {
	bodyString := string(body)
	startIndex := strings.Index(bodyString, "'sizes':") + len("'sizes':")
	endIndex := strings.Index(bodyString[startIndex:], "]")
	return body[startIndex : startIndex+endIndex+1]
}

func sizeToStringArray(sizes []models.Size) (stringSizes []string) {

	for _, size := range sizes {
		stringSizes = append(stringSizes, size.Name)
	}
	return stringSizes
}

func getLinkWithFilter(link string) string {
	return link + filterPagination
}
