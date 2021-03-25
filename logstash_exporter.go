package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
	"bufio"
	"os"
	"path"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/wth12138/logstash_exporter/collector"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"	
)

var (
	scrapeDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: collector.Namespace,
			Subsystem: "exporter",
			Name:      "scrape_duration_seconds",
			Help:      "logstash_exporter: Duration of a scrape job.",
		},
		[]string{"collector", "result"},
	)
)

var configfileName = "/etc/logstash_exporter/conf.yaml"

// LogstashCollector collector type
type LogstashCollector struct {
	collectors map[string]collector.Collector
}

func Exists(path string) bool {
	_, err := os.Stat(path)    
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func CreateAndWriteConfigFile(configfileName string){
	os.MkdirAll(path.Dir(configfileName),os.ModePerm)
	 os.Create(configfileName)

    file, err := os.OpenFile(configfileName, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
		panic(err)
    }
    
    defer file.Close()
	write := bufio.NewWriter(file)

    write.WriteString("endpoint: http://localhost:9600 \n")
	write.WriteString("bindaddress: 9198 \n")
    write.Flush()
}


type BaseConf struct {
	Endpoint    string `yaml:"endpoint"`
	BindAddress string `yaml:"bindaddress"`
}

func (c *BaseConf) GetConf() *BaseConf {
	yamlFile, err := ioutil.ReadFile(configfileName)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

// NewLogstashCollector register a logstash collector
func NewLogstashCollector(logstashEndpoint string) (*LogstashCollector, error) {
	nodeStatsCollector, err := collector.NewNodeStatsCollector(logstashEndpoint)
	if err != nil {
		log.Fatalf("Cannot register a new collector: %v", err)
	}

	nodeInfoCollector, err := collector.NewNodeInfoCollector(logstashEndpoint)
	if err != nil {
		log.Fatalf("Cannot register a new collector: %v", err)
	}

	HotThreadsCollector, err := collector.NewHotThreadsCollector(logstashEndpoint)
	if err != nil {
		log.Fatalf("Cannot register a new collector: %v", err)
	}	

	return &LogstashCollector{
		collectors: map[string]collector.Collector{
			"node": nodeStatsCollector,
			"info": nodeInfoCollector,
			"HotThreads": HotThreadsCollector,
		},
	}, nil
}

func listen(exporterBindAddress string) {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics", http.StatusMovedPermanently)
	})

	log.Infoln("Starting server on", exporterBindAddress)
	if err := http.ListenAndServe(exporterBindAddress, nil); err != nil {
		log.Fatalf("Cannot start Logstash exporter: %s", err)
	}
}

// Describe logstash metrics
func (coll LogstashCollector) Describe(ch chan<- *prometheus.Desc) {
	scrapeDurations.Describe(ch)
}

// Collect logstash metrics
func (coll LogstashCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(coll.collectors))
	for name, c := range coll.collectors {
		go func(name string, c collector.Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
	scrapeDurations.Collect(ch)
}

func execute(name string, c collector.Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Collect(ch)
	duration := time.Since(begin)
	var result string

	if err != nil {
		log.Debugf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		result = "error"
	} else {
		log.Debugf("OK: %s collector succeeded after %fs.", name, duration.Seconds())
		result = "success"
	}
	scrapeDurations.WithLabelValues(name, result).Observe(duration.Seconds())
}

func init() {
	prometheus.MustRegister(version.NewCollector("logstash_exporter"))
}

func main() {
	var (
		logstashEndpoint    = kingpin.Flag("logstash.endpoint", "The protocol, host and port on which logstash metrics API listens").Default("http://localhost:9600").String()
		exporterBindAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9198").String()
	)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("logstash_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	if Exists(configfileName) != true {
		CreateAndWriteConfigFile(configfileName)
	}
	
	if *logstashEndpoint == "http://localhost:9600" && *exporterBindAddress == ":9198" {
		var c BaseConf
		c.GetConf()
		*logstashEndpoint = c.Endpoint
		*exporterBindAddress = fmt.Sprintf(":%v", c.BindAddress)
	}

	logstashCollector, err := NewLogstashCollector(*logstashEndpoint)
	if err != nil {
		log.Fatalf("Cannot register a new Logstash Collector: %v", err)
	}

	prometheus.MustRegister(logstashCollector)

	log.Infoln("Starting Logstash exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())
	listen(*exporterBindAddress)
}
