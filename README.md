# PushD
Prometheus push acceptor for ephemeral and batch jobs. 

## Configuration flags
```text
Prometheus push acceptor for ephemeral and batch jobs

Usage:
  pushd [flags]

Flags:
      --address string            gateway server address (default ":6379")
      --default-buckets strings   default histogram buckets (default [0.005,0.01,0.025,0.05,0.1,0.25,0.5,1,2.5,5,10])
  -h, --help                      help for pushd
      --metrics-address string    metrics server address (default ":9100")
      --metrics-path string       metrics path (default "/metrics")
      --profiling                 enable profiling
      --threads int               number of operating system threads
```

## Getting Started

### Docker
```bash
$ docker run -p 6379:6379 -p 9100:9100 kismia/pushd
```

### Helm
```bash
$ helm install --name=pushd deploy/pushd
```

## Playing with PushD

### Connect from CLI
```bash
$ telnet localhost 6379
# or
$ redis-cli
> PING

# Counter
> CADD metric_name 1 label1 val1 label2 val2 ...
> CINC metric_name label1 val1 label2 val2 ...

# Gauge
> GADD metric_name 1 label1 val1 label2 val2 ...
> GSET metric_name 1 label1 val1 label2 val2 ...
> GSUB metric_name 1 label1 val1 label2 val2 ...
> GINC metric_name label1 val1 label2 val2 ...
> GDEC metric_name label1 val1 label2 val2 ...

# Histogram
> HIST metric_name 100 label1 val1 label2 val2 ...

# Summary
> SUMM metric_name 100 label1 val1 label2 val2 ...

> QUIT
```

### Get metrics
```
$ curl localhost:9100/metrics
```

## Benchmarks
```bash
$ redis-benchmark -r 1000000 -n 1000000 CADD my_counter __rand_int__

====== CADD my_counter __rand_int__ ======
  1000000 requests completed in 10.82 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1

99.99% <= 1 milliseconds
100.00% <= 1 milliseconds
92429.98 requests per second

```

_Running on a MacBook Pro 15, 2018" 2.2 GHz Intel Core i7 using Go 1.12.1_

## Kubernetes integration

#### Kubernetes Service annotations 
```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    prometheus.io/port: "9100"
    prometheus.io/scrape: "pushd"
```

#### Prometheus job example
```yaml
  - job_name: 'kubernetes-pushd-endpoints'
    kubernetes_sd_configs:
      - role: endpoints
    relabel_configs:
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scrape]
        action: keep
        regex: pushd
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scheme]
        action: replace
        target_label: __scheme__
        regex: (https?)
      - source_labels: [__meta_kubernetes_service_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_service_annotation_prometheus_io_port]
        action: replace
        target_label: __address__
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
    metric_relabel_configs:
      - regex: 'client_uuid'
        action: labeldrop
```

## Client Libraries

- C: [hiredis](https://github.com/redis/hiredis)
- C#: [StackExchange.Redis](https://github.com/StackExchange/StackExchange.Redis)
- C++: [redox](https://github.com/hmartiro/redox)
- Clojure: [carmine](https://github.com/ptaoussanis/carmine)
- Common Lisp: [CL-Redis](https://github.com/vseloved/cl-redis)
- Erlang: [Eredis](https://github.com/wooga/eredis)
- Go: [go-redis](https://github.com/go-redis/redis) ([example code](https://github.com/tidwall/tile38/wiki/Go-example-(go-redis)))
- Go: [redigo](https://github.com/gomodule/redigo) ([example code](https://github.com/tidwall/tile38/wiki/Go-example-(redigo)))
- Haskell: [hedis](https://github.com/informatikr/hedis)
- Java: [lettuce](https://github.com/mp911de/lettuce) ([example code](https://github.com/tidwall/tile38/wiki/Java-example-(lettuce)))
- Node.js: [node_redis](https://github.com/NodeRedis/node_redis) ([example code](https://github.com/tidwall/tile38/wiki/Node.js-example-(node-redis)))
- Perl: [perl-redis](https://github.com/PerlRedis/perl-redis)
- PHP: [tinyredisclient](https://github.com/ptrofimov/tinyredisclient) ([example code](https://github.com/tidwall/tile38/wiki/PHP-example-(tinyredisclient)))
- PHP: [phpredis](https://github.com/phpredis/phpredis)
- Python: [redis-py](https://github.com/andymccurdy/redis-py) ([example code](https://github.com/tidwall/tile38/wiki/Python-example))
- Ruby: [redic](https://github.com/amakawa/redic) ([example code](https://github.com/tidwall/tile38/wiki/Ruby-example-(redic)))
- Ruby: [redis-rb](https://github.com/redis/redis-rb) ([example code](https://github.com/tidwall/tile38/wiki/Ruby-example-(redis-rb)))
- Rust: [redis-rs](https://github.com/mitsuhiko/redis-rs)
- Scala: [scala-redis](https://github.com/debasishg/scala-redis)
- Swift: [Redbird](https://github.com/czechboy0/Redbird)