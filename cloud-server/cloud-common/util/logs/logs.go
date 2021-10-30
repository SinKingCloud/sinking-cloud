package logs

type Logs interface {
	Println(data ...interface{})
}

// NewLogstash 实例化NewLogstash
func NewLogstash(hostname string, port int, timeout int) Logs {
	l := &Logstash{
		Hostname:   hostname,
		Port:       port,
		Connection: nil,
		Timeout:    timeout,
	}
	return l
}
