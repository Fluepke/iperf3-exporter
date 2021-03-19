package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
	"time"
)

type Collector struct {
	Timeout      time.Duration
	Iperf3Path   string
	Target       string
	Duration     time.Duration
	OmitDuration time.Duration
	MSS          int
}

var (
	successDesc *prometheus.Desc

	localPortDesc         *prometheus.Desc
	remotePortDesc        *prometheus.Desc
	versionDesc           *prometheus.Desc
	systemInfoDesc        *prometheus.Desc
	tcpMssDesc            *prometheus.Desc
	socketBufferSizeDesc  *prometheus.Desc
	sendBufferSizeDesc    *prometheus.Desc
	receiveBufferSizeDesc *prometheus.Desc

	protocolDesc   *prometheus.Desc
	numStreamsDesc *prometheus.Desc
	omitDesc       *prometheus.Desc
	durationDesc   *prometheus.Desc
	bytesDesc      *prometheus.Desc
	blocksDesc     *prometheus.Desc
	reverseDesc    *prometheus.Desc

	intervalStreamsSecondsDesc               *prometheus.Desc
	intervalStreamsBytesDesc                 *prometheus.Desc
	intervalStreamsRetransmitsDesc           *prometheus.Desc
	intervalStreamsCongestionWindowSizeDesc  *prometheus.Desc
	intervalStreamsRoundTripTimeDesc         *prometheus.Desc
	intervalStreamsRoundTripTimeVarianceDesc *prometheus.Desc
	intervalStreamsPathMTUDesc               *prometheus.Desc

	intervalSummarySecondsDesc       *prometheus.Desc
	intervalSummaryBytesDesc         *prometheus.Desc
	intervalSummaryRetransmittedDesc *prometheus.Desc

	endStreamsSenderSecondsDesc                 *prometheus.Desc
	endStreamsSenderBytesDesc                   *prometheus.Desc
	endStreamsSenderRetransmitsDesc             *prometheus.Desc
	endStreamsSenderMaxSendCongestionWindowDesc *prometheus.Desc
	endStreamsSenderMaxRoundTripTimeDesc        *prometheus.Desc
	endStreamsSenderMinRoundTripTimeDesc        *prometheus.Desc
	endStreamsSenderMeanRoundTripTimeDesc       *prometheus.Desc

	endStreamsReceiverSecondsDesc *prometheus.Desc
	endStreamsReceiverBytesDesc   *prometheus.Desc

	sumSentSecondsDesc     *prometheus.Desc
	sumSentBytesDesc       *prometheus.Desc
	sumReceivedSecondsDesc *prometheus.Desc
	sumReceivedBytesDesc   *prometheus.Desc

	cpuUtilizationPercentHostTotalDesc    *prometheus.Desc
	cpuUtilizationPercentHostUserDesc     *prometheus.Desc
	cpuUtilizationPercentHostSystemDesc   *prometheus.Desc
	cpuUtilizationPercentRemoteTotalDesc  *prometheus.Desc
	cpuUtilizationPercentRemoteUserDesc   *prometheus.Desc
	cpuUtilizationPercentRemoteSystemDesc *prometheus.Desc

	senderTcpCongestionDesc   *prometheus.Desc
	receiverTcpCongestionDesc *prometheus.Desc
)

func init() {
	successDesc = prometheus.NewDesc("iperf3_success", "1 if probe was succesfull", nil, nil)

	localPortDesc = prometheus.NewDesc("iperf3_local_port_info", "Local port", []string{"socket", "local_host"}, nil)
	remotePortDesc = prometheus.NewDesc("iperf3_remote_port_info", "Remote port", []string{"socket", "reote_host"}, nil)
	versionDesc = prometheus.NewDesc("iperf3_version_info", "Iperf3 version information", []string{"version"}, nil)
	systemInfoDesc = prometheus.NewDesc("iperf3_system_info", "System information", []string{"system_info"}, nil)
	tcpMssDesc = prometheus.NewDesc("iperf3_tcp_mss_bytes", "TCPP maximum segment size", nil, nil)
	socketBufferSizeDesc = prometheus.NewDesc("iperf3_socket_buffer_size_bytes", "Socket buffer size", nil, nil)
	sendBufferSizeDesc = prometheus.NewDesc("iperf3_send_buffer_size_bytes", "Send buffer size", nil, nil)
	receiveBufferSizeDesc = prometheus.NewDesc("iperf3_receive_buffer_size_bytes", "Receive buffer size", nil, nil)

	protocolDesc = prometheus.NewDesc("iperf3_protocol_info", "Test protocol", []string{"protocol"}, nil)
	numStreamsDesc = prometheus.NewDesc("iperf3_num_streams_info", "Number of streams", nil, nil)
	omitDesc = prometheus.NewDesc("iperf3_omit_seconds", "Seconds to omit to skip past the TCP slow-start period", nil, nil)
	durationDesc = prometheus.NewDesc("iperf3_duration_seconds", "Test duration", nil, nil)
	bytesDesc = prometheus.NewDesc("iperf3_bytes", "Test bytes to transfer", nil, nil)
	blocksDesc = prometheus.NewDesc("iperf3_blocks_count", "Test blocks to transfer", nil, nil)
	reverseDesc = prometheus.NewDesc("iperf3_reverse_bool", "Wheter to run test in reverse", nil, nil)

	intervalStreamsLabels := []string{"socket", "start", "end", "omitted", "sender"}
	intervalStreamsSecondsDesc = prometheus.NewDesc("iperf3_intervals_streams_seconds", "Duration of the interval in seconds", intervalStreamsLabels, nil)
	intervalStreamsBytesDesc = prometheus.NewDesc("iperf3_intervals_streams_bytes", "Bytes transferred in interval", intervalStreamsLabels, nil)
	intervalStreamsRetransmitsDesc = prometheus.NewDesc("iperf3_intervals_streams_retransmits_count", "Retransmissions in interval", intervalStreamsLabels, nil)
	intervalStreamsCongestionWindowSizeDesc = prometheus.NewDesc("iperf3_intervals_streams_congestion_window_size_byte", "TCP congestion window size in interval", intervalStreamsLabels, nil)
	intervalStreamsRoundTripTimeDesc = prometheus.NewDesc("iperf3_intervals_streams_round_trip_time_seconds", "Round trip time in interval", intervalStreamsLabels, nil)
	intervalStreamsRoundTripTimeVarianceDesc = prometheus.NewDesc("iperf3_intervals_streams_round_trip_time_variance", "Round trip time variance in interval", intervalStreamsLabels, nil)
	intervalStreamsPathMTUDesc = prometheus.NewDesc("iperf3_intervals_streams_path_mtu", "Path MTU discovered in interval", intervalStreamsLabels, nil)

	intervalSummaryLabels := []string{"start", "end", "omitted", "sender"}
	intervalSummarySecondsDesc = prometheus.NewDesc("iperf3_intervals_summary_seconds", "Duration of the interval in seconds", intervalSummaryLabels, nil)
	intervalSummaryBytesDesc = prometheus.NewDesc("iperf3_intervals_summary_bytes", "Total bytes transferred in interval", intervalSummaryLabels, nil)
	intervalSummaryRetransmittedDesc = prometheus.NewDesc("iperf3_intervals_summary_retransmits_count", "Total retransmits in interval", intervalSummaryLabels, nil)

	endStreamsLabels := []string{"socket", "start", "end", "sender"}
	endStreamsSenderSecondsDesc = prometheus.NewDesc("iperf3_end_streams_sender_seconds", "Total send time for stream", endStreamsLabels, nil)
	endStreamsSenderBytesDesc = prometheus.NewDesc("iperf3_end_streams_sender_bytes", "Total bytes send in stream", endStreamsLabels, nil)
	endStreamsSenderRetransmitsDesc = prometheus.NewDesc("iperf3_end_streams_sender_retransmits", "Total retransmit count in stream", endStreamsLabels, nil)
	endStreamsSenderMaxSendCongestionWindowDesc = prometheus.NewDesc("iperf3_end_streams_sender_max_send_congestion_window_bytes", "Maximum send congestion window size", endStreamsLabels, nil)
	endStreamsSenderMaxRoundTripTimeDesc = prometheus.NewDesc("iperf3_end_streams_sender_max_round_trip_time", "Maximum round trip time", endStreamsLabels, nil)
	endStreamsSenderMinRoundTripTimeDesc = prometheus.NewDesc("iperf3_end_streams_sender_min_round_trip_time", "Minimum round trip time", endStreamsLabels, nil)
	endStreamsSenderMeanRoundTripTimeDesc = prometheus.NewDesc("iperf3_end_streams_sender_mean_round_trip_time", "Mean round trip time", endStreamsLabels, nil)

	endStreamsReceiverSecondsDesc = prometheus.NewDesc("iperf3_end_streams_receiver_seconds", "Total receive time for stream", endStreamsLabels, nil)
	endStreamsReceiverBytesDesc = prometheus.NewDesc("iperf3_end_streams_receiver_bytes", "Total received bytes in stream", endStreamsLabels, nil)

	sumSentSecondsDesc = prometheus.NewDesc("iperf3_sum_sent_seconds", "Total send duration", nil, nil)
	sumSentBytesDesc = prometheus.NewDesc("iperf3_sum_sent_bytes", "Total bytes sent", nil, nil)
	sumReceivedSecondsDesc = prometheus.NewDesc("iperf3_sum_received_seconds", "Total receive duration", nil, nil)
	sumReceivedBytesDesc = prometheus.NewDesc("iperf3_sum_received_bytes", "Total received bytes", nil, nil)

	cpuUtilizationPercentHostTotalDesc = prometheus.NewDesc("iperf3_cpu_utilization_host_total_percent", "CPU utilization host total", nil, nil)
	cpuUtilizationPercentHostUserDesc = prometheus.NewDesc("iperf3_cpu_utilization_host_user_percent", "CPU utilization host user", nil, nil)
	cpuUtilizationPercentHostSystemDesc = prometheus.NewDesc("iperf3_cpu_utilization_host_system_percent", "CPU utilization host system", nil, nil)
	cpuUtilizationPercentRemoteTotalDesc = prometheus.NewDesc("iperf3_cpu_utilization_remote_total_percent", "CPU utilization remote total", nil, nil)
	cpuUtilizationPercentRemoteUserDesc = prometheus.NewDesc("iperf3_cpu_utilization_remote_user_percent", "CPU utilization remote user", nil, nil)
	cpuUtilizationPercentRemoteSystemDesc = prometheus.NewDesc("iperf3_cpu_utilization_remote_system_percent", "CPU utilization remote system", nil, nil)

	senderTcpCongestionDesc = prometheus.NewDesc("iperf3_sender_tcp_congestion_control_algorithm_info", "Sender TCP congestion control algorithm", []string{"algorithm"}, nil)
	receiverTcpCongestionDesc = prometheus.NewDesc("iperf3_receiver_tcp_congestion_control_algorithm_info", "Receiver TCP congestion control algorithm", []string{"algorithm"}, nil)
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- successDesc

	ch <- localPortDesc
	ch <- remotePortDesc
	ch <- versionDesc
	ch <- systemInfoDesc
	ch <- tcpMssDesc
	ch <- socketBufferSizeDesc
	ch <- sendBufferSizeDesc
	ch <- receiveBufferSizeDesc

	ch <- protocolDesc
	ch <- numStreamsDesc
	ch <- omitDesc
	ch <- durationDesc
	ch <- bytesDesc
	ch <- blocksDesc
	ch <- reverseDesc

	ch <- intervalStreamsSecondsDesc
	ch <- intervalStreamsBytesDesc
	ch <- intervalStreamsRetransmitsDesc
	ch <- intervalStreamsCongestionWindowSizeDesc
	ch <- intervalStreamsRoundTripTimeDesc
	ch <- intervalStreamsRoundTripTimeVarianceDesc
	ch <- intervalStreamsPathMTUDesc

	ch <- intervalSummarySecondsDesc
	ch <- intervalSummaryBytesDesc
	ch <- intervalSummaryRetransmittedDesc

	ch <- endStreamsSenderSecondsDesc
	ch <- endStreamsSenderBytesDesc
	ch <- endStreamsSenderRetransmitsDesc
	ch <- endStreamsSenderMaxSendCongestionWindowDesc
	ch <- endStreamsSenderMaxRoundTripTimeDesc
	ch <- endStreamsSenderMinRoundTripTimeDesc
	ch <- endStreamsSenderMeanRoundTripTimeDesc

	ch <- endStreamsReceiverSecondsDesc
	ch <- endStreamsReceiverBytesDesc

	ch <- sumSentSecondsDesc
	ch <- sumSentBytesDesc
	ch <- sumReceivedSecondsDesc
	ch <- sumReceivedBytesDesc

	ch <- cpuUtilizationPercentHostTotalDesc
	ch <- cpuUtilizationPercentHostUserDesc
	ch <- cpuUtilizationPercentHostSystemDesc
	ch <- cpuUtilizationPercentRemoteTotalDesc
	ch <- cpuUtilizationPercentRemoteUserDesc
	ch <- cpuUtilizationPercentRemoteSystemDesc

	ch <- senderTcpCongestionDesc
	ch <- receiverTcpCongestionDesc
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	logger := log.WithFields(log.Fields{
		"iperf3_path":   c.Iperf3Path,
		"target":        c.Target,
		"duration":      c.Duration,
		"omit_duration": c.OmitDuration,
		"mss":           c.MSS,
	})

	logger.Debug("Performing iperf3")

	out, err := exec.CommandContext(ctx, c.Iperf3Path, "-J", "-M", strconv.Itoa(c.MSS), "-t", strconv.FormatFloat(c.Duration.Seconds(), 'f', 0, 64), "-O", strconv.FormatFloat(c.OmitDuration.Seconds(), 'f', 0, 64), "-C", "reno", "-c", c.Target).Output()

	logger.Debug("iperf3 done")

	if err != nil {
		logger.WithFields(log.Fields{
			"err": err,
		}).Error("iperf3 probe failed")
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0)
		return
	}

	results := &Iperf3Results{}
	if err := json.Unmarshal(out, &results); err != nil {
		logger.WithFields(log.Fields{
			"err": err,
		}).Error("Deserialize iperf3 results failed")
		ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(successDesc, prometheus.GaugeValue, 1)

	reportMetrics(results, ch)
}

func reportMetrics(r *Iperf3Results, ch chan<- prometheus.Metric) {
	for _, info := range r.Start.Connected {
		ch <- prometheus.MustNewConstMetric(localPortDesc, prometheus.GaugeValue, float64(info.LocalPort), strconv.Itoa(info.Socket), info.LocalHost)
		ch <- prometheus.MustNewConstMetric(remotePortDesc, prometheus.GaugeValue, float64(info.RemotePort), strconv.Itoa(info.Socket), info.RemoteHost)
	}
	ch <- prometheus.MustNewConstMetric(versionDesc, prometheus.GaugeValue, 1, r.Start.Version)
	ch <- prometheus.MustNewConstMetric(systemInfoDesc, prometheus.GaugeValue, 1, r.Start.SystemInfo)

	ch <- prometheus.MustNewConstMetric(tcpMssDesc, prometheus.GaugeValue, float64(r.Start.TcpMSS))
	ch <- prometheus.MustNewConstMetric(socketBufferSizeDesc, prometheus.GaugeValue, float64(r.Start.SocketBufferSize))
	ch <- prometheus.MustNewConstMetric(sendBufferSizeDesc, prometheus.GaugeValue, float64(r.Start.SendBufferSize))
	ch <- prometheus.MustNewConstMetric(receiveBufferSizeDesc, prometheus.GaugeValue, float64(r.Start.ReceiveBufferSize))

	ch <- prometheus.MustNewConstMetric(protocolDesc, prometheus.GaugeValue, 1, r.Start.TestStart.Protocol)
	ch <- prometheus.MustNewConstMetric(numStreamsDesc, prometheus.GaugeValue, float64(r.Start.TestStart.NumStreams))
	ch <- prometheus.MustNewConstMetric(omitDesc, prometheus.GaugeValue, float64(r.Start.TestStart.Omit))
	ch <- prometheus.MustNewConstMetric(durationDesc, prometheus.GaugeValue, float64(r.Start.TestStart.Duration))
	ch <- prometheus.MustNewConstMetric(bytesDesc, prometheus.GaugeValue, float64(r.Start.TestStart.Bytes))
	ch <- prometheus.MustNewConstMetric(blocksDesc, prometheus.GaugeValue, float64(r.Start.TestStart.Blocks))
	ch <- prometheus.MustNewConstMetric(reverseDesc, prometheus.GaugeValue, float64(r.Start.TestStart.Reverse))

	for _, interval := range r.Intervals {
		for _, stream := range interval.Streams {
			labels := []string{
				strconv.Itoa(stream.Socket),
				fmt.Sprintf("%f", stream.Start),
				fmt.Sprintf("%f", stream.End),
				strconv.FormatBool(stream.Omitted),
				strconv.FormatBool(stream.Sender),
			}
			ch <- prometheus.MustNewConstMetric(intervalStreamsSecondsDesc, prometheus.GaugeValue, stream.Seconds, labels...)
			ch <- prometheus.MustNewConstMetric(intervalStreamsBytesDesc, prometheus.GaugeValue, float64(stream.Bytes), labels...)
			ch <- prometheus.MustNewConstMetric(intervalStreamsRetransmitsDesc, prometheus.GaugeValue, float64(stream.Retransmits), labels...)
			ch <- prometheus.MustNewConstMetric(intervalStreamsCongestionWindowSizeDesc, prometheus.GaugeValue, float64(stream.SendCongestionWindowSize), labels...)
			ch <- prometheus.MustNewConstMetric(intervalStreamsRoundTripTimeDesc, prometheus.GaugeValue, stream.RoundTripTime/1000000, labels...)
			ch <- prometheus.MustNewConstMetric(intervalStreamsRoundTripTimeVarianceDesc, prometheus.GaugeValue, stream.RoundTripTimeVariance, labels...)
			ch <- prometheus.MustNewConstMetric(intervalStreamsPathMTUDesc, prometheus.GaugeValue, float64(stream.PathMTU), labels...)
		}

		summary := interval.Summary
		labels := []string{
			fmt.Sprintf("%f", summary.Start),
			fmt.Sprintf("%f", summary.End),
			strconv.FormatBool(summary.Omitted),
			strconv.FormatBool(summary.Sender),
		}
		ch <- prometheus.MustNewConstMetric(intervalSummarySecondsDesc, prometheus.GaugeValue, summary.Seconds, labels...)
		ch <- prometheus.MustNewConstMetric(intervalSummaryBytesDesc, prometheus.GaugeValue, float64(summary.Bytes), labels...)
		ch <- prometheus.MustNewConstMetric(intervalSummaryRetransmittedDesc, prometheus.GaugeValue, float64(summary.Retransmits), labels...)
	}

	for _, stream := range r.End.Streams {
		senderLabels := []string{
			strconv.Itoa(stream.Sender.Socket),
			fmt.Sprintf("%f", stream.Sender.Start),
			fmt.Sprintf("%f", stream.Sender.End),
			strconv.FormatBool(stream.Sender.Sender),
		}
		receiverLabels := []string{
			strconv.Itoa(stream.Receiver.Socket),
			fmt.Sprintf("%f", stream.Receiver.Start),
			fmt.Sprintf("%f", stream.Receiver.End),
			strconv.FormatBool(stream.Receiver.Sender),
		}
		ch <- prometheus.MustNewConstMetric(endStreamsSenderSecondsDesc, prometheus.GaugeValue, stream.Sender.Seconds, senderLabels...)
		ch <- prometheus.MustNewConstMetric(endStreamsSenderBytesDesc, prometheus.GaugeValue, float64(stream.Sender.Bytes), senderLabels...)
		ch <- prometheus.MustNewConstMetric(endStreamsSenderRetransmitsDesc, prometheus.GaugeValue, float64(stream.Sender.Retransmits), senderLabels...)
		ch <- prometheus.MustNewConstMetric(endStreamsSenderMaxSendCongestionWindowDesc, prometheus.GaugeValue, float64(stream.Sender.MaxSendCongestionWindowSize), senderLabels...)
		ch <- prometheus.MustNewConstMetric(endStreamsSenderMaxRoundTripTimeDesc, prometheus.GaugeValue, stream.Sender.MaxRoundTripTime/1000000, senderLabels...)
		ch <- prometheus.MustNewConstMetric(endStreamsSenderMinRoundTripTimeDesc, prometheus.GaugeValue, stream.Sender.MinRoundTripTime/1000000, senderLabels...)
		ch <- prometheus.MustNewConstMetric(endStreamsSenderMeanRoundTripTimeDesc, prometheus.GaugeValue, stream.Sender.MeanRoundTripTime/1000000, senderLabels...)

		ch <- prometheus.MustNewConstMetric(endStreamsReceiverSecondsDesc, prometheus.GaugeValue, stream.Receiver.Seconds, receiverLabels...)
		ch <- prometheus.MustNewConstMetric(endStreamsReceiverBytesDesc, prometheus.GaugeValue, float64(stream.Receiver.Bytes), receiverLabels...)
	}

	ch <- prometheus.MustNewConstMetric(sumSentSecondsDesc, prometheus.GaugeValue, r.End.SummarySent.Seconds)
	ch <- prometheus.MustNewConstMetric(sumSentBytesDesc, prometheus.GaugeValue, float64(r.End.SummarySent.Bytes))
	ch <- prometheus.MustNewConstMetric(sumReceivedSecondsDesc, prometheus.GaugeValue, r.End.SummaryReceived.Seconds)
	ch <- prometheus.MustNewConstMetric(sumReceivedBytesDesc, prometheus.GaugeValue, float64(r.End.SummaryReceived.Bytes))

	ch <- prometheus.MustNewConstMetric(cpuUtilizationPercentHostTotalDesc, prometheus.GaugeValue, r.End.CpuUsage.HostTotal)
	ch <- prometheus.MustNewConstMetric(cpuUtilizationPercentHostUserDesc, prometheus.GaugeValue, r.End.CpuUsage.HostUser)
	ch <- prometheus.MustNewConstMetric(cpuUtilizationPercentHostSystemDesc, prometheus.GaugeValue, r.End.CpuUsage.HostSystem)
	ch <- prometheus.MustNewConstMetric(cpuUtilizationPercentRemoteTotalDesc, prometheus.GaugeValue, r.End.CpuUsage.RemoteTotal)
	ch <- prometheus.MustNewConstMetric(cpuUtilizationPercentRemoteUserDesc, prometheus.GaugeValue, r.End.CpuUsage.RemoteUser)
	ch <- prometheus.MustNewConstMetric(cpuUtilizationPercentRemoteSystemDesc, prometheus.GaugeValue, r.End.CpuUsage.RemoteSystem)

	ch <- prometheus.MustNewConstMetric(senderTcpCongestionDesc, prometheus.GaugeValue, 1, r.End.SenderTcpCongestion)
	ch <- prometheus.MustNewConstMetric(receiverTcpCongestionDesc, prometheus.GaugeValue, 1, r.End.ReceiverTcpCongestion)
}
