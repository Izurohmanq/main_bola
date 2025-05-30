package kafka

import (
	"context"
	configKafka "order-service/config"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type (
	TopicName string
	Handler   func(ctx context.Context, message *sarama.ConsumerMessage) error
)

type ConsumerGroup struct {
	handler map[TopicName]Handler
}

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (*ConsumerGroup) ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error {
	panic("unimplemented")
}

func NewConsumerGroup() *ConsumerGroup {
	return &ConsumerGroup{handler: make(map[TopicName]Handler)}
}

func (c *ConsumerGroup) Setup(sarama sarama.ConsumerGroupSession) error {
	logrus.Infof("setup consumer group")
	return nil
}

func (c *ConsumerGroup) Cleanup(sarama sarama.ConsumerGroupSession) error {
	logrus.Infof("cleanup consumer group")
	return nil
}

func (c *ConsumerGroup) ConsumerClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	messages := claim.Messages()
	for message := range messages {
		handler, ok := c.handler[TopicName(message.Topic)]
		if !ok {
			logrus.Errorf("handler for topic %s not found", message.Topic)
			continue
		}

		var err error
		maxRetry := configKafka.Config.Kafka.MaxRetry
		for attempt := 1; attempt < maxRetry; attempt++ {
			err = handler(context.Background(), message)

			if err == nil {
				break
			}

			logrus.Errorf("error handling message pm %s, attemp %d: %v", message.Topic, attempt, err)
			if attempt == maxRetry {
				logrus.Errorf("max retry reached, message will be ignored")
			}
		}

		if err != nil {
			logrus.Errorf("error handling message on %s: %v", message.Topic, err)
			break
		}

		session.MarkMessage(message, time.Now().UTC().String()) // menandakan bahwa kafka sudah berhasil mengirim datanya
	}
	return nil
}

func (c *ConsumerGroup) RegisterHandler(topic TopicName, handler Handler) {
	c.handler[topic] = handler
	logrus.Infof("register handler for topic %s", topic)
}
