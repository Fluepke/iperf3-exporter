# A better iperf3 prometheus exporter

This software exporters  **all** metrics measured by `iperf3` for use with prometheus.

## Usage
```
Usage of ./iperf3-exporter:
  -iper3.omitTime duration
    	Omit the first  n  seconds  of the test, to skip past the TCP slow-start period (default 5s)
  -iperf3.mss int
    	Set TCP/SCTP maximum segment size (MTU - 40 bytes) (default 1400)
  -iperf3.path string
    	iper3 binary path (default "iperf3")
  -iperf3.time duration
    	time in seconds to transmit for (default 10s)
  -iperf3.timeout duration
    	iperf3 timeout (default 30s)
  -log.level string
    	Logging level (default "info")
  -web.listen-address string
    	Address to listen on for web interface and telemetry (default ":9579")
```

## Prometheus configuration
```yaml
scrape_configs:
  - job_name: 'iperf3'
    static_configs:
      - targets:
        - some.speedtest.server.com
    metrics_path: /probe
    params:
      duration: [ "10s" ]  # overwrite -iperf3.time
      omit-duration: [ "5s" ]  # overwrite -iper3.omitTime
      mss: [ "1400" ]  # overwrite -iperf3.mss
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: localhost:9579
```
