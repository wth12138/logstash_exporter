package collector

import (
	"net/url"
	"strconv"
)

const (
	threads             = 32
	stacktrace_size     = 0
	ignore_idle_threads = false
)

type Hot_threads struct {
	Time            interface{} `json:"time"`
	Busiest_threads int         `json:"busiest_threads"`
	Threads         []struct {
		Name                string   `json:"name"`
		Thread_id           int      `json:"thread_id"`
		Percent_of_cpu_time float64  `json:"percent_of_cpu_time"`
		State               string   `json:"state"`
		Traces              []string `json:"traces"`
	}
}

// HotThreadsResponse type
type HotThreadsResponse struct {
	Host         string `json:"host"`
	Version      string `json:"version"`
	HTTPAddress  string `json:"http_address"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	Ephemeral_id string `json:"ephemeral_id"`
	Status       string `json:"status"`
	Snapshot     bool   `json:"snapshot"`
	Pipeline     struct {
		Workers    int `json:"workers"`
		BatchSize  int `json:"batch_size"`
		BatchDelay int `json:"batch_delay"`
	} `json:"pipeline"`
	Monitoring struct {
		Hosts    []string `json:"hosts"`
		Username string   `json:"username"`
	} `json:"monitoring"`
	Hot_threads Hot_threads `json:"hot_threads"`
}

func ParamsEndpoint(path string) string {
	params := url.Values{}
	Url, err := url.Parse(path + "/_node/hot_threads")
	if err != nil {
		panic(err)
	}
	params.Set("threads", strconv.Itoa(threads))
	params.Set("stacktrace_size", strconv.Itoa(stacktrace_size))
	params.Set("ignore_idle_threads", strconv.FormatBool(ignore_idle_threads))

	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	return urlPath

}

//HotThreads function
func HotThreads(endpoint string) (HotThreadsResponse, error) {
	var response HotThreadsResponse

	urlPath := ParamsEndpoint(endpoint)
	handler := &HTTPHandler{
		Endpoint: urlPath,
	}

	err := getMetrics(handler, &response)

	return response, err
}
