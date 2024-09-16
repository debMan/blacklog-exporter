package kafkaclient

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type KafkaClient struct {
	consumer *kafka.Consumer
	config   *Config
	logger   *zap.Logger
	metrics  Metrics
}

type KafkaMessage struct {
	Value     map[string]float64
	Timestamp time.Time
	Offset    kafka.Offset
}

func NewKafkaClient(config Config, logger *zap.Logger) *KafkaClient {

	kc := &KafkaClient{
		config:   &config,
		logger:   logger,
		consumer: nil,
		metrics:  NewMetrics(),
	}
	return kc
}

func (kc *KafkaClient) ConsumerInitialize() (*KafkaClient, error) {
	var err error

	cm := kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(kc.config.BootstrapServers, ","),
		"group.id":           kc.config.GroupId,
		"enable.auto.commit": kc.config.AutoCommit,
		"auto.offset.reset":  kc.config.AutoOffsetReset,
	}

	if kc.config.SecurityProtocol != "" {
		cm.SetKey("security.protocol", kc.config.SecurityProtocol)
	}
	if kc.config.SaslMechanisms != "" {
		cm.SetKey("sasl.mechanisms", kc.config.SaslMechanisms)
	}
	if kc.config.SaslMechanisms != "" {
		cm.SetKey("sasl.username", kc.config.SaslUsername)
	}
	if kc.config.SaslMechanisms != "" {
		cm.SetKey("sasl.password", kc.config.SaslPassword)
	}
	if kc.config.Debug {
		cm.SetKey("debug", kc.config.Debug)
	}

	kc.consumer, err = kafka.NewConsumer(&cm)

	if err != nil {
		// panic(err)
		kc.logger.Panic("could not connect to kafka", zap.Error(err))
	}

	defer kc.consumer.Close()
	return kc, err
}

func (kc *KafkaClient) Dispose() {
	kc.consumer.Close()
}

func (kc *KafkaClient) StartBlackboxTest(hostName string) {
	outputName := "dummy"
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`)

	if kc.config.Topics == nil {
		kc.logger.Panic("at least one topic is required.")
	}
	kc.ConsumerInitialize()
	err := kc.consumer.SubscribeTopics(kc.config.Topics, nil)
	if err != nil {
		kc.logger.Panic("Failed to subscribe to topics", zap.Error(err))
	}

	msg, err := kc.consumer.ReadMessage(-1)
	if err == nil {
		kc.logger.Info(fmt.Sprintf("Consumed message: %s", string(msg.Value)))

		// Extract timestamp from message and calculate latency
		timestampStr := re.FindString(string(msg.Value))
		if timestampStr != "" {
			timestamp, err := time.Parse(time.RFC3339, timestampStr)
			if err == nil {
				latency := time.Since(timestamp).Seconds()
				kc.metrics.Latency.With(prometheus.Labels{
					"host":   hostName,
					"output": outputName,
				}).Observe(latency)
				kc.logger.Info(fmt.Sprintf("Message latency: %.4f seconds", latency))
			}
		}
	} else {
		kc.logger.Panic(fmt.Sprintf("Consumer error: %v (%v)", err, msg), zap.Error(err))
	}
}
