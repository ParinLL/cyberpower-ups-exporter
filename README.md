# CyberPower UPS SNMP Exporter for Prometheus

This exporter collects SNMP data from CyberPower PowerPanel Business server and exposes it as Prometheus metrics. It's designed to work with CyberPower PowerPanel Business on Linux with SNMPv1 enabled.

## Prerequisites

- Go 1.15 or higher
- CyberPower PowerPanel Business installed on a Linux system
- SNMPv1 enabled on the CyberPower PowerPanel Business server
- At least one UPS connected to the PowerPanel Business server via USB or serial console

## Installation

1. Clone this repository:
git clone https://github.com/yourusername/cyberpower-ups-exporter.git
Copy
2. Navigate to the project directory:
cd cyberpower-ups-exporter
Copy
3. Build the exporter:
go build -o cyberpower_exporter
Copy
## Configuration

Set the following environment variables:

- `SNMP_TARGET`: IP address of your CyberPower PowerPanel Business server
- `SNMP_PORT`: SNMP port (default is 161)
- `SNMP_COMMUNITY`: SNMP community string (default is "public")

## Usage

You can run the exporter in two ways:

1. Set environment variables before running:
```bash
export SNMP_TARGET=192.168.1.1
export SNMP_PORT=161
export SNMP_COMMUNITY=public
go run cyberpower_exporter.go
```
2. Set environment variables inline:
```bash
SNMP_TARGET=192.168.1.1 SNMP_PORT=161 SNMP_COMMUNITY=public go run cyberpower_exporter.go
```
The exporter will start serving metrics on `http://localhost:9100/metrics`.

## Metrics

This exporter provides the following metrics:

1. `ups_battery_status`: UPS Battery Status
- 1: Unknown
- 2: Normal
- 3: Low

2. `ups_battery_replace_indicator`: UPS Battery Replace Indicator
- 1: No
- 2: Replace

3. `ups_input_line_fail_cause`: UPS Input Line Fail Cause
- 1: No Transfer
- 2: High Voltage
- 3: Brown Out
- 4: Self Test

4. `ups_input_status`: UPS Input Status
- 1: Normal
- 2: Over Voltage
- 3: Under Voltage
- 4: Frequency Failure
- 5: Blackout

5. `ups_output_status`: UPS Output Status
- 1: Unknown
- 2: Online
- 3: On Battery
- 4: On Boost
- 5: On Sleep
- 6: Off
- 7: Rebooting

6. `ups_battery_capacity`: UPS Battery Capacity (percentage)

7. `ups_output_current`: UPS Output Current (amperes)

8. `ups_battery_voltage`: UPS Battery Voltage (volts)

9. `ups_input_line_voltage`: UPS Input Line Voltage (volts)

10. `ups_output_voltage`: UPS Output Voltage (volts)

11. `ups_output_load`: UPS Output Load (percentage)

12. `ups_battery_temperature`: UPS Battery Temperature (degrees Celsius)

13. `ups_environment_sensor_temperature`: UPS Environment Sensor Temperature (degrees Celsius)

14. `ups_battery_runtime`: UPS Battery Runtime (minutes)

15. `ups_input_frequency`: UPS Input Frequency (Hz)

16. `ups_output_frequency`: UPS Output Frequency (Hz)

17. `ups_environment_sensor_humidity`: UPS Environment Sensor Humidity (percentage)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.