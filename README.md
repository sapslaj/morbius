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
* `StdoutDestination` - Outputs the flow to stdout in JSON format. Useful for testing and debugging.
* `ElasticsearchDestination` - Indexes the flow into an [Elasticsearch](https://www.elastic.co/elasticsearch/) index
* `LokiDestination` - Pushes the flow to [Loki](https://grafana.com/oss/loki/)

## Use Case

I wanted something that could process NetFlow records on a small-ish scale, like for a homelab (<1Gbps-ish). I wanted it to be as self-contained as possible and with relatively minimal resource utilization (so no Kafka, like is used in the original Cloudflare project). I also wanted something more targetted to NetFlow processing and not a general purpose log/event pipeline (e.g. Logstash, Filebeat) because I find those can be very cumbersome to use with NetFlow and also really limit the amount of enrichment you can do.

## How to use it

A bunch of stuff is still hardcoded but it's mostly controllable from the top level `main.go`. If you want you can implement your own `main` and use this as a library as well. Actual configuration system and deployment method are WIP.
