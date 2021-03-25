package collector

import (
	"errors"
	"strconv"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
)

type ThreadState uint16

const (
	UNKNOWN  ThreadState = iota
	Runnable 
	Blocked
	Waiting
	Timed_Waiting
)

func parseThreadState(s string) (ThreadState, error) {
	s = strings.ToLower(s)
	switch s {
	case "runnable":
		return Runnable, nil
	case "blocked":
		return Blocked, nil
	case "waiting":
		return Waiting, nil
	case "timed_waiting":
		return Timed_Waiting, nil
	default:
		err := errors.New("wrong state")
		return UNKNOWN, err
	}

}

type HotThreadsCollector struct {
	endpoint string

	BusiestThreads          *prometheus.Desc
	ThreadspercentOfCpuTime *prometheus.Desc
	ThreadsState            *prometheus.Desc
}

// NewHotThreadsCollector function
func NewHotThreadsCollector(logstashEndpoint string) (Collector, error) {
	const subsystem = "hotThraeds"

	return &HotThreadsCollector{
		endpoint: logstashEndpoint,

		BusiestThreads: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "Busiest"),
			"BusiestThreads",
			nil,
			nil,
		),

		ThreadspercentOfCpuTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "percentOfCpuTime"),
			"percentOfCpuTime",
			[]string{"name", "thread_id"},
			nil,
		),

		ThreadsState: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "state"),
			"state",
			[]string{"name", "thread_id"},
			nil,
		),
	}, nil
}

// Collect function implements HotThreadsCollector collector
func (c *HotThreadsCollector) Collect(ch chan<- prometheus.Metric) error {
	response, err := HotThreads(c.endpoint)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.BusiestThreads,
		prometheus.GaugeValue,
		float64(response.Hot_threads.Busiest_threads),
	)

	for _, thread := range response.Hot_threads.Threads {
		ThreadState,err:=parseThreadState(thread.State)
		if err != nil {
			return  err
		}
		ch <- prometheus.MustNewConstMetric(
			c.ThreadsState,
			prometheus.GaugeValue,
			float64(ThreadState),
			thread.Name,
			strconv.Itoa(thread.Thread_id),
		)

		ch <- prometheus.MustNewConstMetric(
			c.ThreadspercentOfCpuTime,
			prometheus.GaugeValue,
			thread.Percent_of_cpu_time,
			thread.Name,
			strconv.Itoa(thread.Thread_id),
		)

	}
	//threads := [...]string{}

	//threads = response.Hot_threads.Threads

	//c.CollectHotThreads(threads, ch)

	return nil
}
