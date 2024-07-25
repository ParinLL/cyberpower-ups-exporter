package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	snmpTarget string
	snmpPort   = uint16(161)
	community  string
	logger     *zap.Logger
)

type upsCollector struct {
	batteryStatus                *prometheus.Desc
	batteryReplaceIndicator      *prometheus.Desc
	inputLineFailCause           *prometheus.Desc
	inputStatus                  *prometheus.Desc
	outputStatus                 *prometheus.Desc
	batteryCapacity              *prometheus.Desc
	outputCurrent                *prometheus.Desc
	batteryVoltage               *prometheus.Desc
	inputLineVoltage             *prometheus.Desc
	outputVoltage                *prometheus.Desc
	outputLoad                   *prometheus.Desc
	batteryTemperature           *prometheus.Desc
	environmentSensorTemperature *prometheus.Desc
	batteryRuntime               *prometheus.Desc
	inputFrequency               *prometheus.Desc
	outputFrequency              *prometheus.Desc
	environmentSensorHumidity    *prometheus.Desc
}

func newUPSCollector() *upsCollector {
	return &upsCollector{
		batteryStatus: prometheus.NewDesc(
			"ups_battery_status",
			"UPS Battery Status",
			nil, nil,
		),
		batteryReplaceIndicator: prometheus.NewDesc(
			"ups_battery_replace_indicator",
			"UPS Battery Replace Indicator",
			nil, nil,
		),
		inputLineFailCause: prometheus.NewDesc(
			"ups_input_line_fail_cause",
			"UPS Input Line Fail Cause",
			nil, nil,
		),
		inputStatus: prometheus.NewDesc(
			"ups_input_status",
			"UPS Input Status",
			nil, nil,
		),
		outputStatus: prometheus.NewDesc(
			"ups_output_status",
			"UPS Output Status",
			nil, nil,
		),
		batteryCapacity: prometheus.NewDesc(
			"ups_battery_capacity",
			"UPS Battery Capacity",
			nil, nil,
		),
		outputCurrent: prometheus.NewDesc(
			"ups_output_current",
			"UPS Output Current",
			nil, nil,
		),
		batteryVoltage: prometheus.NewDesc(
			"ups_battery_voltage",
			"UPS Battery Voltage",
			nil, nil,
		),
		inputLineVoltage: prometheus.NewDesc(
			"ups_input_line_voltage",
			"UPS Input Line Voltage",
			nil, nil,
		),
		outputVoltage: prometheus.NewDesc(
			"ups_output_voltage",
			"UPS Output Voltage",
			nil, nil,
		),
		outputLoad: prometheus.NewDesc(
			"ups_output_load",
			"UPS Output Load",
			nil, nil,
		),
		batteryTemperature: prometheus.NewDesc(
			"ups_battery_temperature",
			"UPS Battery Temperature",
			nil, nil,
		),
		environmentSensorTemperature: prometheus.NewDesc(
			"ups_environment_sensor_temperature",
			"UPS Environment Sensor Temperature",
			nil, nil,
		),
		batteryRuntime: prometheus.NewDesc(
			"ups_battery_runtime",
			"UPS Battery Runtime",
			nil, nil,
		),
		inputFrequency: prometheus.NewDesc(
			"ups_input_frequency",
			"UPS Input Frequency",
			nil, nil,
		),
		outputFrequency: prometheus.NewDesc(
			"ups_output_frequency",
			"UPS Output Frequency",
			nil, nil,
		),
		environmentSensorHumidity: prometheus.NewDesc(
			"ups_environment_sensor_humidity",
			"UPS Environment Sensor Humidity",
			nil, nil,
		),
	}
}

func (collector *upsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.batteryStatus
	ch <- collector.batteryReplaceIndicator
	ch <- collector.inputLineFailCause
	ch <- collector.inputStatus
	ch <- collector.outputStatus
	ch <- collector.batteryCapacity
	ch <- collector.outputCurrent
	ch <- collector.batteryVoltage
	ch <- collector.inputLineVoltage
	ch <- collector.outputVoltage
	ch <- collector.outputLoad
	ch <- collector.batteryTemperature
	ch <- collector.environmentSensorTemperature
	ch <- collector.batteryRuntime
	ch <- collector.inputFrequency
	ch <- collector.outputFrequency
	ch <- collector.environmentSensorHumidity
}

func (collector *upsCollector) Collect(ch chan<- prometheus.Metric) {
	snmp := &gosnmp.GoSNMP{
		Target:    snmpTarget,
		Port:      snmpPort,
		Community: community,
		Version:   gosnmp.Version1,
		Timeout:   time.Duration(2) * time.Second,
		Retries:   1,
	}
	err := snmp.Connect()
	if err != nil {
		log.Println("Connect error:", err)
		return
	}
	defer snmp.Conn.Close()

	oids := []string{
		".1.3.6.1.4.1.3808.1.1.1.2.1.1.0", // upsBaseBatteryStatus
		".1.3.6.1.4.1.3808.1.1.1.2.2.5.0", // upsAdvanceBatteryReplaceIndicator
		".1.3.6.1.4.1.3808.1.1.1.3.2.5.0", // upsAdvanceInputLineFailCause
		".1.3.6.1.4.1.3808.1.1.1.3.2.6.0", // upsAdvanceInputStatus
		".1.3.6.1.4.1.3808.1.1.1.4.1.1.0", // upsBaseOutputStatus
		".1.3.6.1.4.1.3808.1.1.1.2.2.1.0", // upsAdvanceBatteryCapacity
		".1.3.6.1.4.1.3808.1.1.1.4.2.4.0", // upsAdvanceOutputCurrent
		".1.3.6.1.4.1.3808.1.1.1.2.2.2.0", // upsAdvanceBatteryVoltage
		".1.3.6.1.4.1.3808.1.1.1.3.2.1.0", // upsAdvanceInputLineVoltage
		".1.3.6.1.4.1.3808.1.1.1.4.2.1.0", // upsAdvanceOutputVoltage
		".1.3.6.1.4.1.3808.1.1.1.4.2.3.0", // upsAdvanceOutputLoad
		".1.3.6.1.4.1.3808.1.1.1.2.2.3.0", // upsAdvanceBatteryTemperature
		".1.3.6.1.4.1.3808.1.1.4.2.1.0",   // envirTemperature
		".1.3.6.1.4.1.3808.1.1.1.2.2.4.0", // upsAdvanceBatteryRunTimeRemaining
		".1.3.6.1.4.1.3808.1.1.1.3.2.4.0", // upsAdvanceInputFrequency
		".1.3.6.1.4.1.3808.1.1.1.4.2.2.0", // upsAdvanceOutputFrequency
		".1.3.6.1.4.1.3808.1.1.4.3.1.0",   // envirHumidity
	}

	result, err := snmp.Get(oids)
	if err != nil {
		log.Println("Get error:", err)
		return
	}

	for _, variable := range result.Variables {
		var value float64
		switch variable.Type {
		case gosnmp.OctetString:
			value, _ = strconv.ParseFloat(string(variable.Value.([]byte)), 64)
		case gosnmp.Integer:
			value = float64(variable.Value.(int))
		case gosnmp.Gauge32:
			value = float64(variable.Value.(uint))
		}

		switch variable.Name {
		case ".1.3.6.1.4.1.3808.1.1.1.2.1.1.0":
			ch <- prometheus.MustNewConstMetric(collector.batteryStatus, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.2.2.5.0":
			ch <- prometheus.MustNewConstMetric(collector.batteryReplaceIndicator, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.3.2.5.0":
			ch <- prometheus.MustNewConstMetric(collector.inputLineFailCause, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.3.2.6.0":
			ch <- prometheus.MustNewConstMetric(collector.inputStatus, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.4.1.1.0":
			ch <- prometheus.MustNewConstMetric(collector.outputStatus, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.2.2.1.0":
			ch <- prometheus.MustNewConstMetric(collector.batteryCapacity, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.4.2.4.0":
			ch <- prometheus.MustNewConstMetric(collector.outputCurrent, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.1.2.2.2.0":
			ch <- prometheus.MustNewConstMetric(collector.batteryVoltage, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.1.3.2.1.0":
			ch <- prometheus.MustNewConstMetric(collector.inputLineVoltage, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.1.4.2.1.0":
			ch <- prometheus.MustNewConstMetric(collector.outputVoltage, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.1.4.2.3.0":
			ch <- prometheus.MustNewConstMetric(collector.outputLoad, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.2.2.3.0":
			ch <- prometheus.MustNewConstMetric(collector.batteryTemperature, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.4.2.1.0":
			ch <- prometheus.MustNewConstMetric(collector.environmentSensorTemperature, prometheus.GaugeValue, (value-32)*5/9/10)
		case ".1.3.6.1.4.1.3808.1.1.1.2.2.4.0":
			ch <- prometheus.MustNewConstMetric(collector.batteryRuntime, prometheus.GaugeValue, value/6000)
		case ".1.3.6.1.4.1.3808.1.1.1.3.2.4.0":
			ch <- prometheus.MustNewConstMetric(collector.inputFrequency, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.1.4.2.2.0":
			ch <- prometheus.MustNewConstMetric(collector.outputFrequency, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.4.3.1.0":
			ch <- prometheus.MustNewConstMetric(collector.environmentSensorHumidity, prometheus.GaugeValue, value)
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom ResponseWriter to capture the status code
		crw := &customResponseWriter{ResponseWriter: w}
		next.ServeHTTP(crw, r)

		duration := time.Since(start)

		// Log the request details
		logger.Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", crw.status),
			zap.Duration("duration", duration),
			zap.String("ip", r.RemoteAddr),
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

func main() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	snmpTarget = os.Getenv("SNMP_TARGET")
	if snmpTarget == "" {
		logger.Fatal("SNMP_TARGET environment variable is not set")
	}

	if port := os.Getenv("SNMP_PORT"); port != "" {
		if p, err := strconv.ParseUint(port, 10, 16); err == nil {
			snmpPort = uint16(p)
		} else {
			logger.Warn("Invalid SNMP_PORT, using default", zap.Uint16("port", snmpPort))
		}
	}

	community = os.Getenv("SNMP_COMMUNITY")
	if community == "" {
		community = "public"
	}

	collector := newUPSCollector()
	prometheus.MustRegister(collector)

	http.Handle("/metrics", loggingMiddleware(promhttp.Handler()))

	addr := ":9100"
	logger.Info("Beginning to serve on port " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Fatal("Error starting HTTP server", zap.Error(err))
	}
}
