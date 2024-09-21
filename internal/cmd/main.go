package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/debman/blacklog-exporter/internal/config"
	"github.com/debman/blacklog-exporter/internal/kafkaclient"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger) {
	kc := kafkaclient.NewKafkaClient(cfg.Kafka, logger)
	kc.StartBlackboxTest()
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan
	kc.Dispose()
	logger.Info("Received termination signal. Exiting...")
	os.Exit(0)
}
