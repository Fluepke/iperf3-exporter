package collector

type Iperf3Results struct {
	Start     *Iperf3Start      `json:"start"`
	Intervals []*Iperf3Interval `json:"intervals"`
	End       *Iperf3End        `json:"end"`
}

type Iperf3Start struct {
	Connected         []*Iperf3Connected `json:"connected"`
	Version           string             `json:"version"`
	SystemInfo        string             `json:"system_info"`
	TcpMSS            int                `json:"tcp_mss"`
	SocketBufferSize  int                `json:"sock_bufsize"`
	SendBufferSize    int                `json:"sndbuf_actual"`
	ReceiveBufferSize int                `json:"rcvbuf_actual"`
	TestStart         *Iperf3TestStart   `json:"test_start"`
}

type Iperf3Connected struct {
	Socket     int    `json:"socket"`
	LocalHost  string `json:"local_host"`
	LocalPort  int    `json:"local_port"`
	RemoteHost string `json:"remote_host"`
	RemotePort int    `json:"remote_port"`
}

type Iperf3TestStart struct {
	Protocol   string `json:"protocol"`
	NumStreams int    `json:"num_streams"`
	BlockSize  int    `json:"blksize"`
	Omit       int    `json:"omit"`
	Duration   int    `json:"duration"`
	Bytes      int    `json:"bytes"`
	Blocks     int    `json:"blocks"`
	Reverse    int    `json:"reverse"`
	TOS        int    `json:"tos"`
}

type Iperf3Interval struct {
	Streams []*Iperf3IntervalStream `json:"streams"`
	Summary *Iperf3IntervalSummary  `json:"sum"`
}

type Iperf3IntervalStream struct {
	Socket                   int     `json:"socket"`
	Start                    float64 `json:"start"`
	End                      float64 `json:"end"`
	Seconds                  float64 `json:"seconds"`
	Bytes                    int     `json:"bytes"`
	BitsPerSecond            float64 `json:"bits_per_second"`
	Retransmits              int     `json:"retransmits"`
	SendCongestionWindowSize int     `json:"snd_cwnd"`
	RoundTripTime            float64 `json:"rtt"`
	RoundTripTimeVariance    float64 `json:"rttvar"`
	PathMTU                  int     `json:"pmtu"`
	Omitted                  bool    `json:"omitted"`
	Sender                   bool    `json:"sender"`
}

type Iperf3IntervalSummary struct {
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         int     `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Retransmits   int     `json:"retransmits"`
	Omitted       bool    `json:"omitted"`
	Sender        bool    `json:"sender"`
}

type Iperf3End struct {
	Streams               []*Iperf3EndStream     `json:"streams"`
	SummarySent           *Iperf3SummarySent     `json:"sum_sent"`
	SummaryReceived       *Iperf3SummaryReceived `json:"sum_received"`
	CpuUsage              *Iperf3CpuUsage        `json:"cpu_utilization_percent"`
	SenderTcpCongestion   string                 `json:"sender_tcp_congestion"`
	ReceiverTcpCongestion string                 `json:"receiver_tcp_congestion"`
}

type Iperf3EndStream struct {
	Sender   *Iperf3Sender   `json:"sender"`
	Receiver *Iperf3Receiver `json:"receiver"`
}

type Iperf3Sender struct {
	Socket                      int     `json:"socket"`
	Start                       float64 `json:"start"`
	End                         float64 `json:"end"`
	Seconds                     float64 `json:"seconds"`
	Bytes                       int     `json:"bytes"`
	BitsPerSecond               float64 `json:"bits_per_second"`
	Retransmits                 int     `json:"retransmits"`
	MaxSendCongestionWindowSize int     `json:"max_snd_cwnd"`
	MaxRoundTripTime            float64 `json:"max_rtt"`
	MinRoundTripTime            float64 `json:"min_rtt"`
	MeanRoundTripTime           float64 `json:"mean_rtt"`
	Sender                      bool    `json:"sender"`
}

type Iperf3Receiver struct {
	Socket        int     `json:"socket"`
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         int     `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Sender        bool    `json:"sender"`
}

type Iperf3SummarySent struct {
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         int     `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Retransmits   int     `json:"retransmits"`
	Sender        bool    `json:"sender"`
}

type Iperf3SummaryReceived struct {
	Start         float64 `json:"start"`
	End           float64 `json:"end"`
	Seconds       float64 `json:"seconds"`
	Bytes         int     `json:"bytes"`
	BitsPerSecond float64 `json:"bits_per_second"`
	Sender        bool    `json:"sender"`
}

type Iperf3CpuUsage struct {
	HostTotal    float64 `json:"host_total"`
	HostUser     float64 `json:"host_user"`
	HostSystem   float64 `json:"host_system"`
	RemoteTotal  float64 `json:"remote_total"`
	RemoteUser   float64 `json:"remote_user"`
	RemoteSystem float64 `json:"remote_system"`
}
