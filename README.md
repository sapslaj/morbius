# morbius

By far one of the NetFlow collectors of all time.

## How it works

It uses [GoFlow from Cloudflare](https://github.com/cloudflare/goflow) to collect NetFlow and implements a custom transports to do stuff with the data. First a set of Enrichers take the flows and do extra processing to add information (such as converting protocol numbers to protocol names and such). Then one or more Destinations are responsible for putting the flow information somewhere.

### Enrichers

* `ProtonamesEnricher` - adds protocol and etype names based on a lookup table
* `RDNSEnricher` - adds rDNS hostname based on IP address fields
* `MaxmindDBEnricher` - adds IP address information from a [MaxMind DB](https://github.com/maxmind/MaxMind-DB)

### Destinations

* `DiscardDestination` - A dummy destination that simply does a JSON marshall and then throws the result away. Used mainly in development.
* `StdoutDestination` - Outputs the flow to stdout in JSON or logfmt format. Useful for testing and debugging.
* `ElasticsearchDestination` - Indexes the flow into an [Elasticsearch](https://www.elastic.co/elasticsearch/) index.
* `LokiDestination` - Pushes the flow to [Loki](https://grafana.com/oss/loki/).
* `PrometheusDestination` - Aggregates flow information info metrics and exposes those in the `:http/metrics` endpoint.

## Use Case

I wanted something that could process NetFlow records on a small-ish scale, like for a homelab (<1Gbps-ish). I wanted it to be as self-contained as possible and with relatively minimal resource utilization (so no Kafka, like is used in the original Cloudflare project). I also wanted something more targetted to NetFlow processing and not a general purpose log/event pipeline (e.g. Logstash, Filebeat) because I find those can be very cumbersome to use with NetFlow and also really limit the amount of enrichment you can do.

## How to use it

Configuration is done via a YAML file. The path to file can be passed with the flag `-config-file`. It defaults to reading `./config.yaml` in the current working directory. An annotated example config file is present at `config.example.yaml`. Using the `-print-config` flag can help debug configuration issues.

It's probably a good idea to create a new config from scratch and only use the example as reference. Here's a decent minimal config to build on with Loki and Prometheus destinations enabled:

```yaml
server:
  netflowv5:
    enable: true
  netflowv9:
    enable: true
  sflow:
    enable: true
  http:
    enable: true
enrichers:
  proto_names:
    enable: true
  rdns:
    enable_cache: true
    cache_size: 2048
  maxmind_db:
    enable_cache: true
    cache_size: 128
    database_paths:
      - /opt/MaxmindDB/GeoLite2-ASN.mmdb
      - /opt/MaxmindDB/GeoLite2-City.mmdb
    enabled_field_groups:
      - asn
      - city
destinations:
  loki:
    push_url: http://loki.monitoring.svc.cluster.local:3100/loki/api/v1/push
    static_labels:
      job: netflow
    dynamic_labels:
      - dst_addr
      - src_addr
  prometheus:
    count_bytes: true
    count_packets: true
    metric_labels:
      - dst_addr
      - dst_port
      - src_port
      - src_addr
      - protocol_name
    export_ip_info: true
    ip_info_labels:
      - addr
      - hostname
      - asn_org
      - asn
      - city_name
      - connection_type
      - continent_name
      - country_name
```
