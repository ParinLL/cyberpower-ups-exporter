package collector

import (
	"errors"
	"strconv"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var (
	ErrNoSNMPTarget = errors.New("SNMP_TARGET environment variable is not set")
)

type Config struct {
	SNMPTarget string
	SNMPPort   string
	Community  string
}

type UPSCollector struct {
	config *Config
	logger *zap.Logger

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
	upsBatteryStatus             *prometheus.Desc
	upsInputLineVoltage          *prometheus.Desc
	upsOutputPercentLoad         *prometheus.Desc
	upsConfigOutputVoltage       *prometheus.Desc
	upsConfigOutputFreq          *prometheus.Desc
}

func NewUPSCollector(config *Config, logger *zap.Logger) (*UPSCollector, error) {
	return &UPSCollector{
		config: config,
		logger: logger,
		upsBatteryStatus: prometheus.NewDesc(
			"ups_battery_status",
			"The present battery status",
			nil, nil,
		),
		upsInputLineVoltage: prometheus.NewDesc(
			"ups_input_line_voltage",
			"The magnitude of the present input voltage",
			nil, nil,
		),
		upsOutputPercentLoad: prometheus.NewDesc(
			"ups_output_percent_load",
			"The percentage of the UPS power capacity presently being used",
			nil, nil,
		),
		upsConfigOutputVoltage: prometheus.NewDesc(
			"ups_config_output_voltage",
			"The nominal output voltage",
			nil, nil,
		),
		upsConfigOutputFreq: prometheus.NewDesc(
			"ups_config_output_freq",
			"The nominal output frequency",
			nil, nil,
		),
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
	}, nil
}

func (c *UPSCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.batteryStatus
	ch <- c.batteryReplaceIndicator
	ch <- c.batteryStatus
	ch <- c.batteryReplaceIndicator
	ch <- c.inputLineFailCause
	ch <- c.inputStatus
	ch <- c.outputStatus
	ch <- c.batteryCapacity
	ch <- c.outputCurrent
	ch <- c.batteryVoltage
	ch <- c.inputLineVoltage
	ch <- c.outputVoltage
	ch <- c.outputLoad
	ch <- c.batteryTemperature
	ch <- c.environmentSensorTemperature
	ch <- c.batteryRuntime
	ch <- c.inputFrequency
	ch <- c.outputFrequency
	ch <- c.environmentSensorHumidity
}

func must(i int, err error) int {
	if err != nil {
		panic(err)
	}
	return i
}

func (c *UPSCollector) Collect(ch chan<- prometheus.Metric) {
	c.logger.Info("Starting metrics collection")
	start := time.Now()
	snmp := &gosnmp.GoSNMP{
		Target:    c.config.SNMPTarget,
		Port:      uint16(must(strconv.Atoi(c.config.SNMPPort))),
		Community: c.config.Community,
		Version:   gosnmp.Version1,
		Timeout:   time.Duration(2) * time.Second,
	}

	err := snmp.Connect()
	if err != nil {
		c.logger.Error("Connect error", zap.Error(err))
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
		c.logger.Error("Get error", zap.Error(err))
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
		case ".1.3.6.1.2.1.33.1.2.1.0":
			ch <- prometheus.MustNewConstMetric(c.upsBatteryStatus, prometheus.GaugeValue, value)
		case ".1.3.6.1.2.1.33.1.3.3.1.2.1":
			ch <- prometheus.MustNewConstMetric(c.upsInputLineVoltage, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.2.1.33.1.4.4.1.5.1":
			ch <- prometheus.MustNewConstMetric(c.upsOutputPercentLoad, prometheus.GaugeValue, value)
		case ".1.3.6.1.2.1.33.1.9.9.0":
			ch <- prometheus.MustNewConstMetric(c.upsConfigOutputVoltage, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.2.1.33.1.9.10.0":
			ch <- prometheus.MustNewConstMetric(c.upsConfigOutputFreq, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.1.2.1.1.0":
			ch <- prometheus.MustNewConstMetric(c.batteryStatus, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.2.2.5.0":
			ch <- prometheus.MustNewConstMetric(c.batteryReplaceIndicator, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.3.2.5.0":
			ch <- prometheus.MustNewConstMetric(c.inputLineFailCause, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.3.2.6.0":
			ch <- prometheus.MustNewConstMetric(c.inputStatus, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.4.1.1.0":
			ch <- prometheus.MustNewConstMetric(c.outputStatus, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.2.2.1.0":
			ch <- prometheus.MustNewConstMetric(c.batteryCapacity, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.4.2.4.0":
			ch <- prometheus.MustNewConstMetric(c.outputCurrent, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.1.2.2.2.0":
			ch <- prometheus.MustNewConstMetric(c.batteryVoltage, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.1.3.2.1.0":
			ch <- prometheus.MustNewConstMetric(c.inputLineVoltage, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.1.4.2.1.0":
			ch <- prometheus.MustNewConstMetric(c.outputVoltage, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.1.4.2.3.0":
			ch <- prometheus.MustNewConstMetric(c.outputLoad, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.1.2.2.3.0":
			ch <- prometheus.MustNewConstMetric(c.batteryTemperature, prometheus.GaugeValue, value)
		case ".1.3.6.1.4.1.3808.1.1.4.2.1.0":
			ch <- prometheus.MustNewConstMetric(c.environmentSensorTemperature, prometheus.GaugeValue, (value-32)*5/9/10)
		case ".1.3.6.1.4.1.3808.1.1.1.2.2.4.0":
			ch <- prometheus.MustNewConstMetric(c.batteryRuntime, prometheus.GaugeValue, value/6000)
		case ".1.3.6.1.4.1.3808.1.1.1.3.2.4.0":
			ch <- prometheus.MustNewConstMetric(c.inputFrequency, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.1.4.2.2.0":
			ch <- prometheus.MustNewConstMetric(c.outputFrequency, prometheus.GaugeValue, value/10)
		case ".1.3.6.1.4.1.3808.1.1.4.3.1.0":
			ch <- prometheus.MustNewConstMetric(c.environmentSensorHumidity, prometheus.GaugeValue, value)
		}
	}
	duration := time.Since(start)
	c.logger.Info("Finished metrics collection", zap.Duration("duration", duration))
}
