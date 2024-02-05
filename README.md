# morbius

By far one of the NetFlow collectors of all time.

## How it works

It uses [GoFlow from Cloudflare](https://github.com/cloudflare/goflow) to collect NetFlow and implements a custom transports to do stuff with the data. First a set of Enrichers take the flows and do extra processing to add information (such as converting protocol numbers to protocol names and such). Then one or more Destinations are responsible for putting the flow information somewhere.

### Enrichers

* `AddrTypeEnricher` - sets a `_type` field based on the type of IP address (`private`, `global`, etc.)
* `FieldMapperEnricher` - allows arbitrary field additions based on either simple key/value mappings or more complex logic. Useful for setting config-specific friendly names e.g. `{in,out}_interface`, `sampler_address`, etc.
* `MaxmindDBEnricher` - adds IP address information from a [MaxMind DB](https://github.com/maxmind/MaxMind-DB)
* `NetDBEnricher` - adds protocol, service, and EtherType information based on [netdb](https://github.com/thediveo/netdb/)
* `ProtonamesEnricher` *(deprecated - use `NetDBEnricher` instead)* - adds protocol and etype names based on a lookup table
* `RDNSEnricher` - adds rDNS hostname based on IP address fields

#### `ProtonamesEnricher` -> `NetDBEnricher` migration

The built-in database for netdb is based on [Debian's netbase](https://salsa.debian.org/md/netbase) project. Unfortunately, that database doesn't contain all of the entries supported by `ProtonamesEnricher` nor does it present the names in the exact same format. Morbius takes care of the missing entries, however there is no special handling for full backwards compatibility. If you need full backwards compatibility, use the following configuration to enable name aliases for the protocols and EtherTypes that will match that `ProtonamesEnricher` outputs:

<details>
<summary>Show configuration</summary>

```yaml
enrichers:
  netdb:
    protocols:
      built_in: true
      name_aliases:
        ah: IPSEC-AH
        hmp: HMP
        hip: HIP
        ddp: DDP
        xtp: XTP
        vmtp: VMTP
        rspf: RSPF
        tcp: TCP
        dccp: DCCP
        ipv6-frag: IPv6-Frag
        hopopt: HOPOPT
        pim: PIM
        manet: MANET
        rsvp: RSVP
        idpr-cmtp: IDPR-CMTP
        skip: SKIP
        ggp: GGP
        ipencap: IP-ENCAP
        l2tp: L2TP
        ipv6: IPv6
        ipv6-opts: IPv6-Opts
        udp: UDP
        udplite: UDPLite
        mobility-header: Mobility-Header
        igmp: IGMP
        shim6: Shim6
        vrrp: VRRP
        ax.25: AX.25
        sctp: SCTP
        ipv6-nonxt: IPv6-NoNxt
        gre: GRE
        mpls-in-ip: MPLS-in-IP
        ipv6-icmp: IPv6-ICMP
        eigrp: EIGRP
        pup: PUP
        ospf: OSPFIGP
        esp: IPSEC-ESP
        encap: ENCAP
        fc: FC
        ipcomp: IPCOMP
        wesp: WESP
        icmp: ICMP
        egp: EGP
        xns-idp: XNS-IDP
        iso-tp4: ISO-TP4
        st: ST
        igp: IGP
        rohc: ROHC
        isis: ISIS
        ipv6-route: IPv6-Route
        idrp: IDRP
        ipip: IPIP
        rdp: RDP
        etherip: ETHERIP
    ethertypes:
      built_in: true
      name_aliases:
        wake-on-lan: Wake-on-LAN
        PPP_DISC: PPPoE Discovery Stage
        PPP_SES: PPPoE Session Stage
        MACSEC: MACsec
        AARP: AppleTalk AARP
        srp: SRP
        ATALK: AppleTalk
        EAPOL: 802.1X
        loopback: Loopback
        S-TAG: S-Tag
        mikrotik-romon: MikroTik RoMON
        qnx-qnet: QNX Qnet
        slpp: SLPP
        epon: EPON
        MPLS_MULTI: MPLS multicast
        802_1Q: C-Tag
        lacp: LACP
        cobranet: CobraNet
        vlacp: VLACP
        avtp: AVTP
        MPLS: MPLS unicast
```

</details>

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
