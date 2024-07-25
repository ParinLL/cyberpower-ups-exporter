package main

import (
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/ParinLL/cyberpower-ups-exporter/collector"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("can't initialize zap logger: " + err.Error())
	}
	defer logger.Sync()

	config, err := getConfig(logger)
	if err != nil {
		logger.Fatal("Failed to get configuration", zap.Error(err))
	}

	upsCollector, err := collector.NewUPSCollector(config, logger)
	if err != nil {
		logger.Fatal("Failed to create UPS collector", zap.Error(err))
	}

	prometheus.MustRegister(upsCollector)

	http.Handle("/metrics", loggingMiddleware(logger, promhttp.Handler()))

	addr := ":9100"
	logger.Info("Beginning to serve on port " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Fatal("Error starting HTTP server", zap.Error(err))
	}
}

func loggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom ResponseWriter to capture the status code
		crw := &customResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(crw, r)

		duration := time.Since(start)

		logger.Info("Metrics request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", crw.status),
			zap.Duration("duration", duration),
			zap.String("remote_addr", r.RemoteAddr),
		)
	})
}

type customResponseWriter struct {
	http.ResponseWriter
	status int
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.status = code
	crw.ResponseWriter.WriteHeader(code)
}

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
