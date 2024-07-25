package main

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/ParinLL/cyberpower-ups-exporter/collector"
	"github.com/ParinLL/cyberpower-ups-exporter/logger"
)

func getConfig(logger *zap.Logger) (*collector.Config, error) {
	snmpTarget := os.Getenv("SNMP_TARGET")
	if snmpTarget == "" {
		logger.Error("SNMP_TARGET environment variable is not set")
		return nil, collector.ErrNoSNMPTarget
	}

	snmpPort := os.Getenv("SNMP_PORT")
	if snmpPort == "" {
		snmpPort = "161" // default SNMP port
		logger.Info("SNMP_PORT not set, using default", zap.String("port", snmpPort))
	}

	community := os.Getenv("SNMP_COMMUNITY")
	if community == "" {
		community = "public" // default community string
		logger.Info("SNMP_COMMUNITY not set, using default", zap.String("community", community))
	}

	return &collector.Config{
		SNMPTarget: snmpTarget,
		SNMPPort:   snmpPort,
		Community:  community,
	}, nil
}

func main() {
	log, err := logger.Initialize()
	if err != nil {
		panic("can't initialize zap logger: " + err.Error())
	}
	defer log.Sync()

	config, err := getConfig(log)
	if err != nil {
		log.Fatal("Failed to get configuration", zap.Error(err))
	}

	upsCollector, err := collector.NewUPSCollector(config, log)
	if err != nil {
		log.Fatal("Failed to create UPS collector", zap.Error(err))
	}

	prometheus.MustRegister(upsCollector)

	http.Handle("/metrics", logger.LoggingMiddleware(log, promhttp.Handler()))

	addr := ":9100"
	log.Info("Beginning to serve on port " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Error starting HTTP server", zap.Error(err))
	}
}
