package kafkaclient

import (
	"encoding/json"
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

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Hostname  string `json:"hostname"`
	// OtherLabel string `json:"OTHER_LABEL"` // You can extend this struct as needed
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

	return kc, err
}

func (kc *KafkaClient) Dispose() {
	kc.consumer.Close()
}

func (kc *KafkaClient) StartBlackboxTest() {
	output_type := "kafka"
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`)

	if kc.config.Topics == nil {
		kc.logger.Panic("at least one topic is required.")
	}

	kc.ConsumerInitialize()
	defer kc.consumer.Close()

	err := kc.consumer.SubscribeTopics(kc.config.Topics, nil)
	if err != nil {
		kc.logger.Panic("Failed to subscribe to topics", zap.Error(err))
	}
	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err == nil {
			kc.logger.Debug(fmt.Sprintf("Consumed message: %s", string(msg.Value)))

			var logEntry LogEntry
			// Unmarshal the Kafka message into the log entry struct
			err = json.Unmarshal(msg.Value, &logEntry)
			if err != nil {
				kc.logger.Error(fmt.Sprintf("Error unmarshalling message: %v", err), zap.Error(err))
				continue
			}

			// Extract timestamp from message and calculate latency
			timestampStr := re.FindString(logEntry.Timestamp)
			if timestampStr != "" {
				kc.logger.Debug(fmt.Sprintf("timestampStr: %s", string(timestampStr)))
				timestampStr += "Z"
				timestamp, err := time.Parse(time.RFC3339, timestampStr)
				if err == nil {
					kc.logger.Debug(fmt.Sprintf("timestamp: %s", (timestamp)))
					latency := time.Since(timestamp).Seconds()
					kc.metrics.Latency.With(prometheus.Labels{
						"hostname":    logEntry.Hostname,
						"output_type": output_type,
					}).Observe(latency)
					kc.logger.Debug(fmt.Sprintf("Message latency: %.4f seconds", latency))
				} else {
					kc.logger.Debug(fmt.Sprintf("cant get timestamp for: %s", (err)))
				}
			}
		} else {
			kc.logger.Panic(fmt.Sprintf("Consumer error: %v (%v)", err, msg), zap.Error(err))
		}
	}
}
