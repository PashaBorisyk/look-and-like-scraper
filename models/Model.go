package models

import (
	"time"
)

type Product struct {
	Domain      string        `json:"domain" bson:"domain"`
	LocaleLCID  string        `json:"localeLCID" bson:"localeLCID"`
	Alpha3Code  string        `json:"alpha3Code" bson:"alpha3Code"`
	ShopName    string        `json:"shopName" bson:"shopName"`
	BaseURL     string        `json:"baseURL" bson:"baseURL"`
	Url         string        `json:"url" bson:"url"`
	Name        string        `json:"name" bson:"name"`
	Color       string        `json:"color" bson:"color"`
	Sizes       []string      `json:"sizes" bson:"sizes"`
	Price       float64       `json:"price" bson:"price"`
	Currency    string        `json:"currency" bson:"currency"`
	Description string        `json:"description" bson:"description"`
	Article     string        `json:"article" bson:"article"`
	Images      []string      `json:"images" bson:"images"`
	InsertDate  time.Time     `json:"insertDate" bson:"insertDate"`
	Category    string        `json:"category" bson:"category"`
	Sex         string        `json:"sex" bson:"sex"`
	Composition []Composition `json:"composition" bson:"composition"`
}

type Composition struct {
	Part     string
	Material string
	Percent  string
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
