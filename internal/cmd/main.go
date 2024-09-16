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
	hostName, err := os.Hostname()
	if err != nil {
		logger.Panic("can't gte hostname", zap.Error(err))
		os.Exit(1)
	}
	kc := kafkaclient.NewKafkaClient(cfg.Kafka, logger)
	kc.StartBlackboxTest(hostName)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan
	kc.Dispose()
	logger.Info("Received termination signal. Exiting...")
	os.Exit(0)
}
