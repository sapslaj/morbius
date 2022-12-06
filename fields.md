# Fields

This is a reference for all potential fields in a flow

| Field                                | JSON Type | Origin             | Note/Description                                                  |
| ------------------------------------ | --------- | ------------------ | ----------------------------------------------------------------- |
| `type`                               | string    | Transport          | The type of flow that this comes from (NetFlow v5/v9, sFlow, etc) |
| `time_received`                      | number    | Transport          | UNIX epoch timestamp                                              |
| `sequence_num`                       | number    | Transport          |                                                                   |
| `sampling_rate `                     | number    | Transport          |                                                                   |
| `flow_direction`                     | number    | Transport          |                                                                   |
| `sampler_address`                    | string    | Transport          |                                                                   |
| `time_flow_start`                    | number    | Transport          | UNIX epoch timestamp                                              |
| `time_flow_end`                      | number    | Transport          | UNIX epoch timestamp                                              |
| `bytes`                              | number    | Transport          |                                                                   |
| `packets`                            | number    | Transport          |                                                                   |
| `src_addr`                           | string    | Transport          |                                                                   |
| `dst_addr`                           | string    | Transport          |                                                                   |
| `ethernet_type`                      | number    | Transport          |                                                                   |
| `proto`                              | number    | Transport          |                                                                   |
| `src_port`                           | number    | Transport          |                                                                   |
| `dst_port`                           | number    | Transport          |                                                                   |
| `in_interface`                       | number    | Transport          |                                                                   |
| `out_interface`                      | number    | Transport          |                                                                   |
| `src_mac`                            | string    | Transport          |                                                                   |
| `dst_mac`                            | string    | Transport          |                                                                   |
| `src_vlan`                           | number    | Transport          |                                                                   |
| `dst_vlan`                           | number    | Transport          |                                                                   |
| `vlan_id`                            | number    | Transport          |                                                                   |
| `ingress_vrf_id`                     | number    | Transport          |                                                                   |
| `egress_vrf_id`                      | number    | Transport          |                                                                   |
| `ip_tos`                             | number    | Transport          |                                                                   |
| `forwarding_status`                  | number    | Transport          |                                                                   |
| `ip_ttl`                             | number    | Transport          |                                                                   |
| `tcp_flags`                          | number    | Transport          |                                                                   |
| `icmp_types`                         | number    | Transport          |                                                                   |
| `icmp_code`                          | number    | Transport          |                                                                   |
| `ipv6_flow_label`                    | number    | Transport          |                                                                   |
| `fragment_id`                        | number    | Transport          |                                                                   |
| `fragment_offset`                    | number    | Transport          |                                                                   |
| `bi_flow_direction`                  | number    | Transport          |                                                                   |
| `src_as`                             | number    | Transport          |                                                                   |
| `dst_as`                             | number    | Transport          |                                                                   |
| `next_hop`                           | string    | Transport          |                                                                   |
| `next_hop_as`                        | number    | Transport          |                                                                   |
| `src_net`                            | number    | Transport          |                                                                   |
| `dst_net`                            | number    | Transport          |                                                                   |
| `has_encap`                          | boolean   | Transport          |                                                                   |
| `src_addr_encap`                     | string    | Transport          |                                                                   |
| `dst_addr_encap`                     | string    | Transport          |                                                                   |
| `proto_encap`                        | number    | Transport          |                                                                   |
| `ethernet_type_encap`                | number    | Transport          |                                                                   |
| `ip_tos_encap`                       | number    | Transport          |                                                                   |
| `ip_ttl_encap`                       | number    | Transport          |                                                                   |
| `ipv6_flow_label_encap`              | number    | Transport          |                                                                   |
| `fragment_id_encap`                  | number    | Transport          |                                                                   |
| `fragment_offset_encap`              | number    | Transport          |                                                                   |
| `has_mpls`                           | boolean   | Transport          |                                                                   |
| `mpls_count`                         | number    | Transport          |                                                                   |
| `mpls_1_ttl`                         | number    | Transport          |                                                                   |
| `mpls_1_label`                       | number    | Transport          |                                                                   |
| `mpls_2_ttl`                         | number    | Transport          |                                                                   |
| `mpls_2_label`                       | number    | Transport          |                                                                   |
| `mpls_3_ttl`                         | number    | Transport          |                                                                   |
| `mpls_3_label`                       | number    | Transport          |                                                                   |
| `mpls_last_ttl`                      | number    | Transport          |                                                                   |
| `mpls_last_label`                    | number    | Transport          |                                                                   |
| `has_ppp`                            | boolean   | Transport          |                                                                   |
| `ppp_address_control`                | number    | Transport          |                                                                   |
| `protocol_name`                      | string    | ProtonamesEnricher |                                                                   |
| `protocol_encap_name`                | string    | ProtonamesEnricher |                                                                   |
| `ethernet_type_name`                 | string    | ProtonamesEnricher |                                                                   |
| `ethernet_type_encap_name`           | string    | ProtonamesEnricher |                                                                   |
| `src_hostname`                       | string    | RDNSEnricher       |                                                                   |
| `dst_hostname`                       | string    | RDNSEnricher       |                                                                   |
| `src_hostname_encap`                 | string    | RDNSEnricher       |                                                                   |
| `dst_hostname_encap`                 | string    | RDNSEnricher       |                                                                   |
| `src_asn`                            | number    | MaxmindDBEnricher  |                                                                   |
| `src_asn_org`                        | string    | MaxmindDBEnricher  |                                                                   |
| `src_average_income`                 | number    | MaxmindDBEnricher  |                                                                   |
| `src_city_confidence`                | number    | MaxmindDBEnricher  |                                                                   |
| `src_city_name`                      | string    | MaxmindDBEnricher  |                                                                   |
| `src_connection_type`                | string    | MaxmindDBEnricher  |                                                                   |
| `src_continent_code`                 | string    | MaxmindDBEnricher  |                                                                   |
| `src_continent_name`                 | string    | MaxmindDBEnricher  |                                                                   |
| `src_country_code`                   | string    | MaxmindDBEnricher  |                                                                   |
| `src_country_confidence`             | number    | MaxmindDBEnricher  |                                                                   |
| `src_country_eu`                     | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_country_name`                   | string    | MaxmindDBEnricher  |                                                                   |
| `src_domain`                         | string    | MaxmindDBEnricher  |                                                                   |
| `src_ip_risk`                        | number    | MaxmindDBEnricher  |                                                                   |
| `src_is_anonymous`                   | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_is_anonymous_proxy`             | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_is_anonymous_vpn`               | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_is_hosting_provider`            | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_is_legitimate_proxy`            | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_isp`                            | string    | MaxmindDBEnricher  |                                                                   |
| `src_is_public_proxy`                | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_is_residential_proxy`           | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_is_satellite_provider`          | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_is_tor_exit_node`               | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_loc_accuracy`                   | number    | MaxmindDBEnricher  |                                                                   |
| `src_loc_lat`                        | number    | MaxmindDBEnricher  |                                                                   |
| `src_loc_long`                       | number    | MaxmindDBEnricher  |                                                                   |
| `src_loc_metro_code`                 | number    | MaxmindDBEnricher  |                                                                   |
| `src_loc_postal_code`                | string    | MaxmindDBEnricher  |                                                                   |
| `src_loc_postal_confidence`          | number    | MaxmindDBEnricher  |                                                                   |
| `src_loc_tz`                         | string    | MaxmindDBEnricher  |                                                                   |
| `src_organization`                   | string    | MaxmindDBEnricher  |                                                                   |
| `src_population_density`             | number    | MaxmindDBEnricher  |                                                                   |
| `src_registered_country_code`        | string    | MaxmindDBEnricher  |                                                                   |
| `src_registered_country_eu`          | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_registered_country_name`        | string    | MaxmindDBEnricher  |                                                                   |
| `src_represented_country_code`       | string    | MaxmindDBEnricher  |                                                                   |
| `src_represented_country_name`       | string    | MaxmindDBEnricher  |                                                                   |
| `src_static_ip_score`                | number    | MaxmindDBEnricher  |                                                                   |
| `src_ip_user_type`                   | string    | MaxmindDBEnricher  |                                                                   |
| `dst_asn`                            | number    | MaxmindDBEnricher  |                                                                   |
| `dst_asn_org`                        | string    | MaxmindDBEnricher  |                                                                   |
| `dst_average_income`                 | number    | MaxmindDBEnricher  |                                                                   |
| `dst_city_confidence`                | number    | MaxmindDBEnricher  |                                                                   |
| `dst_city_name`                      | string    | MaxmindDBEnricher  |                                                                   |
| `dst_connection_type`                | string    | MaxmindDBEnricher  |                                                                   |
| `dst_continent_code`                 | string    | MaxmindDBEnricher  |                                                                   |
| `dst_continent_name`                 | string    | MaxmindDBEnricher  |                                                                   |
| `dst_country_code`                   | string    | MaxmindDBEnricher  |                                                                   |
| `dst_country_confidence`             | number    | MaxmindDBEnricher  |                                                                   |
| `dst_country_eu`                     | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_country_name`                   | string    | MaxmindDBEnricher  |                                                                   |
| `dst_domain`                         | string    | MaxmindDBEnricher  |                                                                   |
| `dst_ip_risk`                        | number    | MaxmindDBEnricher  |                                                                   |
| `dst_is_anonymous`                   | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_is_anonymous_proxy`             | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_is_anonymous_vpn`               | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_is_hosting_provider`            | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_is_legitimate_proxy`            | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_isp`                            | string    | MaxmindDBEnricher  |                                                                   |
| `dst_is_public_proxy`                | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_is_residential_proxy`           | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_is_satellite_provider`          | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_is_tor_exit_node`               | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_loc_accuracy`                   | number    | MaxmindDBEnricher  |                                                                   |
| `dst_loc_lat`                        | number    | MaxmindDBEnricher  |                                                                   |
| `dst_loc_long`                       | number    | MaxmindDBEnricher  |                                                                   |
| `dst_loc_metro_code`                 | number    | MaxmindDBEnricher  |                                                                   |
| `dst_loc_postal_code`                | string    | MaxmindDBEnricher  |                                                                   |
| `dst_loc_postal_confidence`          | number    | MaxmindDBEnricher  |                                                                   |
| `dst_loc_tz`                         | string    | MaxmindDBEnricher  |                                                                   |
| `dst_organization`                   | string    | MaxmindDBEnricher  |                                                                   |
| `dst_population_density`             | number    | MaxmindDBEnricher  |                                                                   |
| `dst_registered_country_code`        | string    | MaxmindDBEnricher  |                                                                   |
| `dst_registered_country_eu`          | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_registered_country_name`        | string    | MaxmindDBEnricher  |                                                                   |
| `dst_represented_country_code`       | string    | MaxmindDBEnricher  |                                                                   |
| `dst_represented_country_name`       | string    | MaxmindDBEnricher  |                                                                   |
| `dst_static_ip_score`                | number    | MaxmindDBEnricher  |                                                                   |
| `dst_ip_user_type`                   | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_asn`                      | number    | MaxmindDBEnricher  |                                                                   |
| `src_encap_asn_org`                  | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_average_income`           | number    | MaxmindDBEnricher  |                                                                   |
| `src_encap_city_confidence`          | number    | MaxmindDBEnricher  |                                                                   |
| `src_encap_city_name`                | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_connection_type`          | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_continent_code`           | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_continent_name`           | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_country_code`             | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_country_confidence`       | number    | MaxmindDBEnricher  |                                                                   |
| `src_encap_country_eu`               | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_encap_country_name`             | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_domain`                   | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_ip_risk`                  | number    | MaxmindDBEnricher  |                                                                   |
| `src_encap_is_anonymous`             | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_encap_is_anonymous_proxy`       | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_encap_is_anonymous_vpn`         | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_encap_is_hosting_provider`      | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_encap_is_legitimate_proxy`      | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_encap_isp`                      | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_is_public_proxy`          | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_encap_is_residential_proxy`     | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_encap_is_satellite_provider`    | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_encap_is_tor_exit_node`         | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_encap_loc_accuracy`             | number    | MaxmindDBEnricher  |                                                                   |
| `src_encap_loc_lat`                  | number    | MaxmindDBEnricher  |                                                                   |
| `src_encap_loc_long`                 | number    | MaxmindDBEnricher  |                                                                   |
| `src_encap_loc_metro_code`           | number    | MaxmindDBEnricher  |                                                                   |
| `src_encap_loc_postal_code`          | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_loc_postal_confidence`    | number    | MaxmindDBEnricher  |                                                                   |
| `src_encap_loc_tz`                   | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_organization`             | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_population_density`       | number    | MaxmindDBEnricher  |                                                                   |
| `src_encap_registered_country_code`  | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_registered_country_eu`    | boolean   | MaxmindDBEnricher  |                                                                   |
| `src_encap_registered_country_name`  | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_represented_country_code` | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_represented_country_name` | string    | MaxmindDBEnricher  |                                                                   |
| `src_encap_static_ip_score`          | number    | MaxmindDBEnricher  |                                                                   |
| `src_encap_ip_user_type`             | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_asn`                      | number    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_asn_org`                  | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_average_income`           | number    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_city_confidence`          | number    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_city_name`                | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_connection_type`          | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_continent_code`           | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_continent_name`           | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_country_code`             | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_country_confidence`       | number    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_country_eu`               | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_encap_country_name`             | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_domain`                   | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_ip_risk`                  | number    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_is_anonymous`             | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_encap_is_anonymous_proxy`       | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_encap_is_anonymous_vpn`         | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_encap_is_hosting_provider`      | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_encap_is_legitimate_proxy`      | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_encap_isp`                      | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_is_public_proxy`          | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_encap_is_residential_proxy`     | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_encap_is_satellite_provider`    | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_encap_is_tor_exit_node`         | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_encap_loc_accuracy`             | number    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_loc_lat`                  | number    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_loc_long`                 | number    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_loc_metro_code`           | number    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_loc_postal_code`          | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_loc_postal_confidence`    | number    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_loc_tz`                   | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_organization`             | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_population_density`       | number    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_registered_country_code`  | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_registered_country_eu`    | boolean   | MaxmindDBEnricher  |                                                                   |
| `dst_encap_registered_country_name`  | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_represented_country_code` | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_represented_country_name` | string    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_static_ip_score`          | number    | MaxmindDBEnricher  |                                                                   |
| `dst_encap_ip_user_type`             | string    | MaxmindDBEnricher  |                                                                   |