package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type KafkaConfig struct {
	CGroup     string   `json:"cGroup"`
	Topic      string   `json:"topic"`
	ServerUrls []string `json:"serverUrls"`
	RetryMax   int      `json:"retryMax"`
}

type MongoConfig struct {
	Uri                  string `json:"uri"`
	ProductsDatabaseName string `json:"productsDatabaseName"`
	CollyDatabaseName    string `json:"collyDatabaseName"`
}

type Locale struct {
	LocaleLCID        string `json:"localeLCID"`
	Alpha3Code        string `json:"alpha3Code"`
	BaseURL           string `json:"baseUrl"`
	MaleTranslation   string `json:"maleTranslation"`
	FemaleTranslation string `json:"femaleTranslation"`
	KidsTranslation   string `json:"kidsTranslation"`
}

type ScraperConfig struct {
	Locales map[string][]Locale `json:"locales"`
}

type Config struct {
	KafkaConfig   KafkaConfig   `json:"kafkaConfig"`
	ScraperConfig ScraperConfig `json:"scraperConfig"`
	MongoConfig   MongoConfig   `json:"mongoConfig"`
}

var config *Config

func init() {

	log.Println("Creating configuration...")

	log.Println("Reading config.json file...")
	fileRaw, err := ioutil.ReadFile("./resources/config.json")

	if err != nil {
		log.Fatal("Can not read from config.json file", err)
	} else {
		log.Println("Reading successful")
	}

	log.Println("Validating config.json file...")
	isValid := json.Valid(fileRaw)
	if !isValid {
		log.Fatal("Invalid config.json file format. JSON structure expected expected")
	} else {
		log.Println("Validating successful")
	}

	log.Println("Decoding config.json file... ")
	err = json.Unmarshal(fileRaw, &config)
	if err != nil {
		log.Fatal("Decoding failed ", err)
	} else {
		log.Println("Config decoded successful: ")
		log.Println(*config)
	}

}

func GetConfig() Config {
	return *config
}
