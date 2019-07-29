package queue

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"look-and-like-web-scrapper/config"
	"look-and-like-web-scrapper/logger"
	"os"
)

var producer *sarama.SyncProducer

func init() {
	configureKafka()
}

func configureKafka(){
	log.Println("Kafka init...")

	var err error
	producer, err = initProducer()
	if err != nil {
		log.Println("Error producer: ", err)
		log.Println("Kafka init failed. You will have to run scissors separately")
	} else {
		log.Println("Kafka init done")
	}
}

func initProducer() (*sarama.SyncProducer, error) {

	log.Println("Init producer")

	kafkaSettings := config.GetConfig().KafkaConfig
	serverUrl := os.Getenv("KAFKA_SERVER_URL")
	if serverUrl == "" {
		panic("KAFKA_SERVER_URL environment variable must be provided")
	}
	serverUrls := []string{serverUrl}
	retryMax := kafkaSettings.RetryMax

	sarama.Logger = log.New(logger.GetOrCreateLogFile("kafka"), "kafka", log.Ltime)
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Retry.Max = retryMax
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true

	prd, err := sarama.NewSyncProducer(serverUrls, saramaConfig)

	return &prd, err

}

func PublishKey(key interface{}) {
	value := fmt.Sprintf("%v", key)
	publish(value, *producer)
}

func publish(message string, producer sarama.SyncProducer) {
	msg := &sarama.ProducerMessage{
		Topic: config.GetConfig().KafkaConfig.Topic,
		Value: sarama.StringEncoder(message),
	}

	if producer != nil {
		_, _, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Println("Error publish: ", err.Error())
		}

		log.Println("Message ", message, " published")
	} else {
		log.Println("No active queue connection. No message will be published; Reconnecting...")
		configureKafka()
	}

}
