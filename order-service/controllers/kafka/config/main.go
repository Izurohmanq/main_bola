package kafka

import (
	"order-service/config"
	"order-service/controllers/kafka"
	kafkaPayment "order-service/controllers/kafka/payment"

	"golang.org/x/exp/slices"
)

type Kafka struct {
	consumer *ConsumerGroup
	kafka    kafka.IKafkaRegistry
}

type IKafka interface {
	Register()
}

func NewKafkaConsumer(consumer *ConsumerGroup, kafka kafka.IKafkaRegistry) IKafka {
	return &Kafka{consumer: consumer, kafka: kafka}
}

func (k *Kafka) paymentHandler() {
	if slices.Contains(config.Config.Kafka.Topics, kafkaPayment.PaymentTopic) {
		k.consumer.RegisterHandler(kafkaPayment.PaymentTopic, k.kafka.GetPayment().HandlePayment)
	} // mengecek untuk apakah topik ini mengandung atau ada topi dengan nama "kafkaPayment.PaymentTopic", maka dia akan melakukan RegisterHandler
}

func (k *Kafka) Register() {
	k.paymentHandler()
}
