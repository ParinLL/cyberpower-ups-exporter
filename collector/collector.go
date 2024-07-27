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

	upsBatteryStatus             *prometheus.Desc
	upsSecondsOnBattery          *prometheus.Desc
	upsEstimatedMinutesRemaining *prometheus.Desc
	upsEstimatedChargeRemaining  *prometheus.Desc
	upsBatteryVoltage            *prometheus.Desc
	upsInputFrequency            *prometheus.Desc
	upsInputVoltage              *prometheus.Desc
	upsOutputSource              *prometheus.Desc
	upsOutputFrequency           *prometheus.Desc
	upsOutputVoltage             *prometheus.Desc
	upsOutputCurrent             *prometheus.Desc
	upsOutputPower               *prometheus.Desc
	upsOutputPercentLoad         *prometheus.Desc
}

func NewUPSCollector(config *Config, logger *zap.Logger) (*UPSCollector, error) {
	return &UPSCollector{
		config: config,
		logger: logger,
		upsBatteryStatus: prometheus.NewDesc(
			"ups_battery_status",
			"The current battery status",
			[]string{"snmp_target"}, nil,
		),
		upsSecondsOnBattery: prometheus.NewDesc(
			"ups_seconds_on_battery",
			"The number of seconds on battery power",
			[]string{"snmp_target"}, nil,
		),
		upsEstimatedMinutesRemaining: prometheus.NewDesc(
			"ups_estimated_minutes_remaining",
			"The estimated minutes of battery runtime remaining",
			[]string{"snmp_target"}, nil,
		),
		upsEstimatedChargeRemaining: prometheus.NewDesc(
			"ups_estimated_charge_remaining",
			"The estimated battery charge remaining in percent",
			[]string{"snmp_target"}, nil,
		),
		upsBatteryVoltage: prometheus.NewDesc(
			"ups_battery_voltage",
			"The current battery voltage in 0.1 Volt DC",
			[]string{"snmp_target"}, nil,
		),
		upsInputFrequency: prometheus.NewDesc(
			"ups_input_frequency",
			"The current input frequency in 0.1 Hertz",
			[]string{"snmp_target"}, nil,
		),
		upsInputVoltage: prometheus.NewDesc(
			"ups_input_voltage",
			"The current input voltage in Volt AC",
			[]string{"snmp_target"}, nil,
		),
		upsOutputSource: prometheus.NewDesc(
			"ups_output_source",
			"The current output source",
			[]string{"snmp_target"}, nil,
		),
		upsOutputFrequency: prometheus.NewDesc(
			"ups_output_frequency",
			"The current output frequency in 0.1 Hertz",
			[]string{"snmp_target"}, nil,
		),
		upsOutputVoltage: prometheus.NewDesc(
			"ups_output_voltage",
			"The current output voltage in Volt AC",
			[]string{"snmp_target"}, nil,
		),
		upsOutputCurrent: prometheus.NewDesc(
			"ups_output_current",
			"The current output current in 0.1 Ampere",
			[]string{"snmp_target"}, nil,
		),
		upsOutputPower: prometheus.NewDesc(
			"ups_output_power",
			"The current output power in Watt",
			[]string{"snmp_target"}, nil,
		),
		upsOutputPercentLoad: prometheus.NewDesc(
			"ups_output_percent_load",
			"The current output load in percent",
			[]string{"snmp_target"}, nil,
		),
	}, nil
}
func (c *UPSCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upsBatteryStatus
	ch <- c.upsSecondsOnBattery
	ch <- c.upsEstimatedMinutesRemaining
	ch <- c.upsEstimatedChargeRemaining
	ch <- c.upsBatteryVoltage
	ch <- c.upsInputFrequency
	ch <- c.upsInputVoltage
	ch <- c.upsOutputSource
	ch <- c.upsOutputFrequency
	ch <- c.upsOutputVoltage
	ch <- c.upsOutputCurrent
	ch <- c.upsOutputPower
	ch <- c.upsOutputPercentLoad
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
		".1.3.6.1.2.1.33.1.2.1.0",     // upsBatteryStatus
		".1.3.6.1.2.1.33.1.2.2.0",     // upsSecondsOnBattery
		".1.3.6.1.2.1.33.1.2.3.0",     // upsEstimatedMinutesRemaining
		".1.3.6.1.2.1.33.1.2.4.0",     // upsEstimatedChargeRemaining
		".1.3.6.1.2.1.33.1.2.5.0",     // upsBatteryVoltage
		".1.3.6.1.2.1.33.1.3.3.1.2.1", // upsInputFrequency
		".1.3.6.1.2.1.33.1.3.3.1.3.1", // upsInputVoltage
		".1.3.6.1.2.1.33.1.4.1.0",     // upsOutputSource
		".1.3.6.1.2.1.33.1.4.2.0",     // upsOutputFrequency
		".1.3.6.1.2.1.33.1.4.4.1.2.1", // upsOutputVoltage
		".1.3.6.1.2.1.33.1.4.4.1.3.1", // upsOutputCurrent
		".1.3.6.1.2.1.33.1.4.4.1.4.1", // upsOutputPower
		".1.3.6.1.2.1.33.1.4.4.1.5.1", // upsOutputPercentLoad
	}

	result, err := snmp.Get(oids)
	if err != nil {
		c.logger.Error("Get error", zap.Error(err))
		return
	}

	snmpTarget := c.config.SNMPTarget
	for _, variable := range result.Variables {
		var value float64
		switch variable.Type {
		case gosnmp.OctetString:
			value, _ = strconv.ParseFloat(string(variable.Value.([]byte)), 64)
		case gosnmp.Integer:
			value = float64(variable.Value.(int))
		case gosnmp.Gauge32:
			value = float64(variable.Value.(uint))
		case gosnmp.TimeTicks:
			value = float64(variable.Value.(uint))
		}

		switch variable.Name {
		case ".1.3.6.1.2.1.33.1.2.1.0":
			ch <- prometheus.MustNewConstMetric(c.upsBatteryStatus, prometheus.GaugeValue, value, snmpTarget)
		case ".1.3.6.1.2.1.33.1.2.2.0":
			ch <- prometheus.MustNewConstMetric(c.upsSecondsOnBattery, prometheus.GaugeValue, value, snmpTarget)
		case ".1.3.6.1.2.1.33.1.2.3.0":
			ch <- prometheus.MustNewConstMetric(c.upsEstimatedMinutesRemaining, prometheus.GaugeValue, value, snmpTarget)
		case ".1.3.6.1.2.1.33.1.2.4.0":
			ch <- prometheus.MustNewConstMetric(c.upsEstimatedChargeRemaining, prometheus.GaugeValue, value, snmpTarget)
		case ".1.3.6.1.2.1.33.1.2.5.0":
			ch <- prometheus.MustNewConstMetric(c.upsBatteryVoltage, prometheus.GaugeValue, value/10, snmpTarget)
		case ".1.3.6.1.2.1.33.1.3.3.1.2.1":
			ch <- prometheus.MustNewConstMetric(c.upsInputFrequency, prometheus.GaugeValue, value/10, snmpTarget)
		case ".1.3.6.1.2.1.33.1.3.3.1.3.1":
			ch <- prometheus.MustNewConstMetric(c.upsInputVoltage, prometheus.GaugeValue, value, snmpTarget)
		case ".1.3.6.1.2.1.33.1.4.1.0":
			ch <- prometheus.MustNewConstMetric(c.upsOutputSource, prometheus.GaugeValue, value, snmpTarget)
		case ".1.3.6.1.2.1.33.1.4.2.0":
			ch <- prometheus.MustNewConstMetric(c.upsOutputFrequency, prometheus.GaugeValue, value/10, snmpTarget)
		case ".1.3.6.1.2.1.33.1.4.4.1.2.1":
			ch <- prometheus.MustNewConstMetric(c.upsOutputVoltage, prometheus.GaugeValue, value, snmpTarget)
		case ".1.3.6.1.2.1.33.1.4.4.1.3.1":
			ch <- prometheus.MustNewConstMetric(c.upsOutputCurrent, prometheus.GaugeValue, value/10, snmpTarget)
		case ".1.3.6.1.2.1.33.1.4.4.1.4.1":
			ch <- prometheus.MustNewConstMetric(c.upsOutputPower, prometheus.GaugeValue, value, snmpTarget)
		case ".1.3.6.1.2.1.33.1.4.4.1.5.1":
			ch <- prometheus.MustNewConstMetric(c.upsOutputPercentLoad, prometheus.GaugeValue, value, snmpTarget)
		}
	}

	duration := time.Since(start)
	c.logger.Info("Finished metrics collection", zap.Duration("duration", duration))
}
