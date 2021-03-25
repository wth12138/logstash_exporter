# Logstash exporter [![Build Status](https://travis-ci.org/sequra/logstash_exporter.svg)]
Prometheus exporter for the metrics available in Logstash since version 5.0.

Continuous integration: [travis](https://travis-ci.org/sequra/logstash_exporter/)

## Version compatibility

As logstash can change API metrics in new versions (which happened on v7.3.0), we decided to change the
version of `logstash_exporter` to adapt to the minimum supported version of logstash using the first three numbers for compatibility and the latest one as own increasing version.

For instance, `logstash_exporter v7.3.0.0` supports a minimum version of logstash 7.3.0, meanwhile `logstash-exporter v5.0.0.0` suppots a minimum version of logstash 5.0.

This change will also be reflected on branch names, existing `master` (the latest supported version), `v7.3.0` (>= logstash 7.3.0), and `v5.0.0` (>= logstash 5.0 and < 7.3.0).

## Usage

```bash
go get github.com/wth12138/logstash_exporter@v0.1.7
cd $GOPATH/src/github.com/sequra/logstash_exporter
go run logstash_exporter.go --web.listen-address=:1234 --logstash.endpoint="http://localhost:1235"
```

### Config
- default path = /etc/logstash_exporter/conf.yaml
```yaml
endpoint: http://localhost:9600
bindaddress: 9198
```


### Flags
- flags has  higher priority than config file,but I used to apply yaml file in k8s environment for more convenient
```sh
logstash_exporter --help
usage: logstash_exporter [<flags>]

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and
                          --help-man).
      --logstash.endpoint="http://localhost:9600"
                          The protocol, host and port on which logstash metrics
                          API listens
      --web.listen-address=":9198"
                          Address on which to expose metrics and web interface.
      --log.level="info"  Only log messages with the given severity or above.
                          Valid levels: [debug, info, warn, error, fatal]
      --log.format="logger:stderr"
                          Set the log target and format. Example:
                          "logger:syslog?appname=bob&local=7" or
                          "logger:stdout?json=true"
      --version           Show application version.
```

## Implemented metrics

* `logstash_exporter_build_info` (gauge)
* `logstash_exporter_scrape_duration_seconds`: logstash_exporter: Duration of a scrape job. (summary)
* `logstash_info_jvm`: A metric with a constant '1' value labeled by name, version and vendor of the JVM running Logstash. (counter)
* `logstash_info_node`: A metric with a constant '1' value labeled by Logstash version. (counter)
* `logstash_info_os`: A metric with a constant '1' value labeled by name, arch, version and available_processors to the OS running Logstash. (counter)
* `logstash_node_gc_collection_duration_seconds_total` (counter)
* `logstash_node_gc_collection_total` (gauge)
* `logstash_node_jvm_threads_count` (gauge)
* `logstash_node_jvm_threads_peak_count` (gauge)
* `logstash_node_mem_heap_committed_bytes` (gauge)
* `logstash_node_mem_heap_max_bytes` (gauge)
* `logstash_node_mem_heap_used_bytes` (gauge)
* `logstash_node_mem_nonheap_committed_bytes` (gauge)
* `logstash_node_mem_nonheap_used_bytes` (gauge)
* `logstash_node_mem_pool_committed_bytes` (gauge)
* `logstash_node_mem_pool_max_bytes` (gauge)
* `logstash_node_mem_pool_peak_max_bytes` (gauge)
* `logstash_node_mem_pool_peak_used_bytes` (gauge)
* `logstash_node_mem_pool_used_bytes` (gauge)
* `logstash_node_pipeline_duration_seconds_total` (counter)
* `logstash_node_pipeline_events_filtered_total` (counter)
* `logstash_node_pipeline_events_in_total` (counter)
* `logstash_node_pipeline_events_out_total` (counter)
* `logstash_node_pipeline_queue_push_duration_seconds_total` (counter)
* `logstash_node_plugin_bulk_requests_failures_total` (counter)
* `logstash_node_plugin_bulk_requests_successes_total` (counter)
* `logstash_node_plugin_bulk_requests_with_errors_total` (counter)
* `logstash_node_plugin_documents_failures_total` (counter)
* `logstash_node_plugin_documents_successes_total` (counter)
* `logstash_node_plugin_duration_seconds_total` (counter
* `logstash_node_plugin_queue_push_duration_seconds_total` (counter)
* `logstash_node_plugin_events_in_total` (counter)
* `logstash_node_plugin_events_out_total` (counter)
* `logstash_node_plugin_current_connections_count` (gauge)
* `logstash_node_plugin_peak_connections_count` (gauge)
* `logstash_node_process_cpu_total_seconds_total` (counter)
* `logstash_node_process_max_filedescriptors` (gauge)
* `logstash_node_process_mem_total_virtual_bytes` (gauge)
* `logstash_node_process_open_filedescriptors` (gauge)
* `logstash_node_queue_events` (counter)
* `logstash_node_queue_size_bytes` (counter)
* `logstash_node_queue_max_size_bytes` (counter)
* `logstash_node_dead_letter_queue_size_bytes` (counter)
* `logstash_node_up`: whether logstash node is up (1) or not (0) (gauge)
* `logstash_hotthreads_BusiestThreads` (gauge)
* `logstash_hotthreads_percentOfCpuTime` (gauge)
* `logstash_hotthreads_state` (gauge)

### logstash_hotthreads_state
- Runnable      1
-	Blocked       2
-	Waiting       3
- Timed_Waiting 4


## Integration tests
In order to execute manual integration tests (to know if certain logstash version is compatible with logstash-exporter), you can follow instructions present on file [integration-tests/README.md](integration-tests/README.md).
