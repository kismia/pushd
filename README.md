# PushD
Prometheus push acceptor for ephemeral and batch jobs.

## RESP API

### Client Libraries

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
- Node.js: [node-tile38](https://github.com/phulst/node-tile38) ([example code](https://github.com/tidwall/tile38/wiki/Node.js-example-(node-tile38)))
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

#### Connect from CLI
```bash
telnet {{ PushdHost }} 6379
# or
redis-cli -h {{ PushdHost }}
```

#### Available commands
```bash
PING

# Counter
CADD metric_name 1 label1 val1 label2 val2 ...
CINC metric_name label1 val1 label2 val2 ...

# Gauge
GADD metric_name 1 label1 val1 label2 val2 ...
GSET metric_name 1 label1 val1 label2 val2 ...
GSUB metric_name 1 label1 val1 label2 val2 ...
GINC metric_name label1 val1 label2 val2 ...
GDEC metric_name label1 val1 label2 val2 ...

# Histogram
HIST metric_name 100 label1 val1 label2 val2 ...

# Summary
SUMM metric_name 100 label1 val1 label2 val2 ...

QUIT
```

## Prometheus job example
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
