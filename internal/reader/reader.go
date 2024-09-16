// TODO: Abstract reader sub service

// package reader

// import (
// 	"log"
// 	"regexp"
// 	"time"

// 	"net/http"

// 	"github.com/confluentinc/confluent-kafka-go/kafka"
// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// )

// // Define Prometheus metrics
// var (
// 	messageLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
// 		Name:    "log_message_latency_seconds",
// 		Help:    "Time taken to consume log messages from Kafka",
// 		Buckets: prometheus.DefBuckets,
// 	})
// )

// func init() {
// 	prometheus.MustRegister(messageLatency)
// }

// func main() {
// 	// Start Prometheus HTTP server
// 	http.Handle("/metrics", promhttp.Handler())
// 	go func() {
// 		log.Fatal(http.ListenAndServe(":2112", nil))
// 	}()

// 	// Kafka consumer configuration
// 	c, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers": "localhost:9092",
// 		"group.id":          "log_consumer_group",
// 		"auto.offset.reset": "earliest",
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed to create consumer: %s", err)
// 	}
// 	defer c.Close()

// 	// Subscribe to the log topic
// 	c.SubscribeTopics([]string{"log_topic"}, nil)

// 	// Regular expression to extract timestamp from log message
// 	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`)

// 	// Consume messages and measure latency
// 	for {
// 		msg, err := c.ReadMessage(-1)
// 		if err == nil {
// 			log.Printf("Consumed message: %s", string(msg.Value))

// 			// Extract timestamp from message and calculate latency
// 			timestampStr := re.FindString(string(msg.Value))
// 			if timestampStr != "" {
// 				timestamp, err := time.Parse(time.RFC3339, timestampStr)
// 				if err == nil {
// 					latency := time.Since(timestamp).Seconds()
// 					messageLatency.Observe(latency)
// 					log.Printf("Message latency: %.4f seconds", latency)
// 				}
// 			}
// 		} else {
// 			log.Printf("Consumer error: %v (%v)", err, msg)
// 		}
// 	}
// }
