// TODO: Add writer sub service

// package writer

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"time"
// )

// // Define the log structure
// type LogEntry struct {
// 	Timestamp string `json:"timestamp"`
// 	Message   string `json:"message"`
// 	Host      string `json:"host"`
// }

// // Generate a random sentence with random words
// func generateRandomWords() string {
// 	wordList := []string{
// 		"cloud", "server", "data", "microservice", "streaming", "network", "container",
// 		"kubernetes", "Go", "Kafka", "log", "process", "API", "distributed", "scale",
// 		"performance", "load", "metric", "resilience", "scalability",
// 	}

// 	// Randomize number of words (3 to 8 words)
// 	numWords := rand.Intn(6) + 3
// 	sentence := ""
// 	for i := 0; i < numWords; i++ {
// 		randomWord := wordList[rand.Intn(len(wordList))]
// 		sentence += randomWord
// 		if i < numWords-1 {
// 			sentence += " "
// 		}
// 	}
// 	return sentence
// }

// // Get current hostname of logger
// func getHostName() string {
// 	return "dummy"
// }

// func main() {
// 	rand.Seed(time.Now().UnixNano()) // Seed random number generator

// 	// Loop to generate logs continuously
// 	for i := 0; i < 100; i++ { // You can adjust this number or make it infinite
// 		logEntry := LogEntry{
// 			Timestamp: time.Now().Format(time.RFC3339), // Current timestamp in RFC3339 format
// 			Message:   generateRandomWords(),           // Generate a random sentence
// 			Host:      getHostName(),                   // Current hostname
// 		}

// 		// Convert the log entry to JSON format
// 		logEntryJSON, err := json.Marshal(logEntry)
// 		if err != nil {
// 			log.Fatalf("Failed to marshal log entry to JSON: %s", err)
// 		}

// 		// Print log to stdout
// 		fmt.Println(string(logEntryJSON))

// 		// Sleep to simulate time between log generation (adjust as needed)
// 		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
// 	}
// }
