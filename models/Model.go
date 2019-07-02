package models

type MetaInformation struct {
	Domain     string `json:"domain" bson:"domain"`
	LocaleLCID string `json:"localeLCID" bson:"localeLCID"`
	Alpha3Code string `json:"alpha3Code" bson:"alpha3Code"`
	ShopName   string `json:"shopName" bson:"shopName"`
	BaseURL    string `json:"baseURL" bson:"baseURL"`
	Url        string `json:"url" bson:"url"`
	InsertDate string `json:"insertDate" bson:"insertDate"`
}

type Price struct {
	Value    float64 `json:"price" bson:"price"`
	Currency string  `json:"currency" bson:"currency"`
}

type Images struct {
	NoBackgroundImageUrl string   `json:"noBackgroundImageUrl" bson:"noBackgroundImageUrl"`
	StockImageUrls       []string `json:"stockImageUrls" bson:"stockImageUrls"`
}

type Composition struct {
	Part     string `json:"part"`
	Material string `json:"material"`
	Percent  string `json:"percent"`
}

type Data struct {
	Name        string        `json:"name" bson:"name"`
	Color       string        `json:"color" bson:"color"`
	Sizes       []string      `json:"sizes" bson:"sizes"`
	Description string        `json:"description" bson:"description"`
	Article     string        `json:"article" bson:"article"`
	Sex         string        `json:"sex" bson:"sex"`
	Category    string        `json:"category" bson:"category"`
	Composition []Composition `json:"composition" bson:"composition"`
	Price       Price         `json:"price" bson:"price"`
	Images      Images        `json:"images" bson:"images"`
}

type Product struct {
	ID              interface{}     `json:"id" bson:"_id"`
	MetaInformation MetaInformation `json:"metaInformation" bson:"metaInformation"`
	Data            Data            `json:"data" bson:"data"`
}

type ProductRep struct {
	Offers struct {
		PriceCurrency string `json:"priceCurrency"`
		Price         string `json:"price"`
	} `json:"offers"`
}

type Size struct {
	Name string `json:"name"`
}
