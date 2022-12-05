package enricher_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sapslaj/morbius/enricher"
)

func TestMaxmindDBEnricher(t *testing.T) {
	type test struct {
		desc   string
		skip   string
		config enricher.MaxmindDBEnricherConfig
		input  map[string]interface{}
		want   map[string]interface{}
	}

	tests := []test{
		{
			desc: "Does not modify message if an address filed is not defined",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoLite2-ASN-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{"other": 69},
			want:  map[string]interface{}{"other": 69},
		},
		{
			desc: "Adds info when src_addr is set",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoLite2-ASN-Test.mmdb"},
				EnabledFields: []enricher.MaxmindDBEnricherField{
					enricher.MaxmindDBEnricherField_ASN,
				},
			},
			input: map[string]interface{}{
				"src_addr": "1.0.0.1",
			},
			want: map[string]interface{}{
				"src_addr": "1.0.0.1",
				"src_asn":  uint64(15169),
			},
		},
		{
			desc: "Adds info when dst_addr is set",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoLite2-ASN-Test.mmdb"},
				EnabledFields: []enricher.MaxmindDBEnricherField{
					enricher.MaxmindDBEnricherField_ASN,
				},
			},
			input: map[string]interface{}{
				"dst_addr": "1.0.0.1",
			},
			want: map[string]interface{}{
				"dst_addr": "1.0.0.1",
				"dst_asn":  uint64(15169),
			},
		},
		{
			desc: "Adds info when src_addr_encap is set",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoLite2-ASN-Test.mmdb"},
				EnabledFields: []enricher.MaxmindDBEnricherField{
					enricher.MaxmindDBEnricherField_ASN,
				},
			},
			input: map[string]interface{}{
				"src_addr_encap": "1.0.0.1",
			},
			want: map[string]interface{}{
				"src_addr_encap": "1.0.0.1",
				"src_encap_asn":  uint64(15169),
			},
		},
		{
			desc: "Adds info when dst_addr_encap is set",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoLite2-ASN-Test.mmdb"},
				EnabledFields: []enricher.MaxmindDBEnricherField{
					enricher.MaxmindDBEnricherField_ASN,
				},
			},
			input: map[string]interface{}{
				"dst_addr_encap": "1.0.0.1",
			},
			want: map[string]interface{}{
				"dst_addr_encap": "1.0.0.1",
				"dst_encap_asn":  uint64(15169),
			},
		},
		{
			desc: "Works with GeoIP2-Anonymous",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoIP2-Anonymous-IP-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "1.2.0.0",
			},
			want: map[string]interface{}{
				"src_addr":                 "1.2.0.0",
				"src_is_anonymous":         true,
				"src_is_anonymous_vpn":     true,
				"src_is_public_proxy":      false,
				"src_is_residential_proxy": false,
				"src_is_tor_exit_node":     false,
			},
		},
		{
			desc: "Works with GeoIP2-City",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoIP2-City-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "214.0.1.0",
			},
			want: map[string]interface{}{
				"src_addr":                    "214.0.1.0",
				"src_city_name":               "Melbourne",
				"src_continent_code":          "OC",
				"src_continent_name":          "Oceania",
				"src_country_code":            "AU",
				"src_country_name":            "Australia",
				"src_loc_accuracy":            uint64(20),
				"src_loc_lat":                 -37.8159,
				"src_loc_long":                144.9669,
				"src_loc_postal_code":         "3000",
				"src_loc_tz":                  "Australia/Melbourne",
				"src_registered_country_code": "AU",
				"src_registered_country_name": "Australia",
			},
		},
		{
			desc: "Works with GeoIP2-Connection-Type",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoIP2-Connection-Type-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "1.0.0.1",
			},
			want: map[string]interface{}{
				"src_addr":            "1.0.0.1",
				"src_connection_type": "Cable/DSL",
			},
		},
		{
			desc: "Works with GeoIP2-Country",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoIP2-City-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "2001:218::",
			},
			want: map[string]interface{}{
				"src_addr":                    "2001:218::",
				"src_continent_code":          "AS",
				"src_continent_name":          "Asia",
				"src_country_code":            "JP",
				"src_country_name":            "Japan",
				"src_loc_accuracy":            uint64(100),
				"src_loc_lat":                 35.68536,
				"src_loc_long":                139.75309,
				"src_loc_tz":                  "Asia/Tokyo",
				"src_registered_country_code": "JP",
				"src_registered_country_name": "Japan",
			},
		},
		{
			desc: "Works with GeoIP2-DensityIncome",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoIP2-DensityIncome-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "5.83.124.0",
			},
			want: map[string]interface{}{
				"src_addr":               "5.83.124.0",
				"src_average_income":     uint64(32323),
				"src_population_density": uint64(1232),
			},
		},
		{
			desc: "Works with GeoIP2-Domain",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoIP2-Domain-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "1.2.0.0",
			},
			want: map[string]interface{}{
				"src_addr":   "1.2.0.0",
				"src_domain": "maxmind.com",
			},
		},
		{
			desc: "Works with GeoIP2-Enterprise",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoIP2-Enterprise-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "2.125.160.216",
			},
			want: map[string]interface{}{
				"src_addr":                    "2.125.160.216",
				"src_city_name":               "Boxford",
				"src_city_confidence":         uint64(50),
				"src_continent_code":          "EU",
				"src_continent_name":          "Europe",
				"src_country_code":            "GB",
				"src_country_name":            "United Kingdom",
				"src_country_confidence":      uint64(95),
				"src_loc_accuracy":            uint64(100),
				"src_loc_lat":                 51.75,
				"src_loc_long":                -1.25,
				"src_loc_postal_code":         "OX1",
				"src_loc_postal_confidence":   uint64(20),
				"src_loc_tz":                  "Europe/London",
				"src_registered_country_eu":   true,
				"src_registered_country_code": "FR",
				"src_registered_country_name": "France",
				"src_static_ip_score":         0.27,
				"src_is_anonymous_proxy":      false,
				"src_is_legitimate_proxy":     false,
				"src_is_satellite_provider":   false,
			},
		},
		{
			desc: "Works with GeoIP2-IP-Risk",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoIP2-IP-Risk-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "214.2.3.0",
			},
			want: map[string]interface{}{
				"src_addr":                 "214.2.3.0",
				"src_ip_risk":              0.1,
				"src_is_anonymous":         true,
				"src_is_anonymous_vpn":     true,
				"src_is_tor_exit_node":     false,
				"src_is_hosting_provider":  false,
				"src_is_public_proxy":      false,
				"src_is_residential_proxy": false,
			},
		},
		{
			desc: "Works with GeoIP2-ISP",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoIP2-ISP-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "1.128.0.0",
			},
			want: map[string]interface{}{
				"src_addr":         "1.128.0.0",
				"src_asn":          uint64(1221),
				"src_asn_org":      "Telstra Pty Ltd",
				"src_isp":          "Telstra Internet",
				"src_organization": "Telstra Internet",
			},
		},
		{
			desc: "Works with GeoIP2-Precision-Enterprise",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoIP2-Precision-Enterprise-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "1.231.232.0",
			},
			want: map[string]interface{}{
				"src_addr":                    "1.231.232.0",
				"src_city_name":               "Dzhankoy",
				"src_city_confidence":         uint64(60),
				"src_continent_code":          "EU",
				"src_continent_name":          "Europe",
				"src_country_code":            "UA",
				"src_country_name":            "Ukraine",
				"src_country_confidence":      uint64(80),
				"src_loc_accuracy":            uint64(200),
				"src_loc_lat":                 45.7117,
				"src_loc_long":                34.3927,
				"src_loc_tz":                  "Europe/Simferopol",
				"src_registered_country_code": "UA",
				"src_registered_country_name": "Ukraine",
				"src_asn":                     uint64(28761),
				"src_asn_org":                 "CrimeaCom South LLC",
				"src_connection_type":         "Cable/DSL",
				"src_static_ip_score":         0.26,
				"src_is_anonymous_proxy":      false,
				"src_is_legitimate_proxy":     false,
				"src_is_satellite_provider":   false,
				"src_isp":                     "CrimeaCom South LLC",
				"src_organization":            "CrimeaCom South LLC",
				"src_ip_user_type":            "residential",
			},
		},
		{
			desc: "Works with GeoIP2-Static-IP-Score",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoIP2-Static-IP-Score-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "1.0.0.1",
			},
			want: map[string]interface{}{
				"src_addr":            "1.0.0.1",
				"src_static_ip_score": 0.01,
			},
		},
		{
			desc: "Works with GeoLite2-ASN",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoLite2-ASN-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "1.0.0.1",
			},
			want: map[string]interface{}{
				"src_addr":    "1.0.0.1",
				"src_asn":     uint64(15169),
				"src_asn_org": "Google Inc.",
			},
		},
		{
			desc: "Works with GeoLite2-City",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoLite2-City-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "2.125.160.216",
			},
			want: map[string]interface{}{
				"src_addr":                    "2.125.160.216",
				"src_city_name":               "Boxford",
				"src_continent_code":          "EU",
				"src_continent_name":          "Europe",
				"src_country_code":            "GB",
				"src_country_name":            "United Kingdom",
				"src_loc_accuracy":            uint64(100),
				"src_loc_lat":                 51.75,
				"src_loc_long":                -1.25,
				"src_loc_postal_code":         "OX1",
				"src_loc_tz":                  "Europe/London",
				"src_registered_country_eu":   true,
				"src_registered_country_code": "FR",
				"src_registered_country_name": "France",
			},
		},
		{
			desc: "Works with GeoLite2-Country",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoLite2-Country-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "2.125.160.216",
			},
			want: map[string]interface{}{
				"src_addr":                    "2.125.160.216",
				"src_continent_code":          "EU",
				"src_continent_name":          "Europe",
				"src_country_code":            "GB",
				"src_country_name":            "United Kingdom",
				"src_registered_country_eu":   true,
				"src_registered_country_code": "FR",
				"src_registered_country_name": "France",
			},
		},
		{
			desc: "Uses configured locale for names",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoLite2-Country-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
				Locale:        "ja",
			},
			input: map[string]interface{}{
				"src_addr": "2001:218::",
			},
			want: map[string]interface{}{
				"src_addr":                    "2001:218::",
				"src_continent_code":          "AS",
				"src_continent_name":          "アジア",
				"src_country_code":            "JP",
				"src_country_name":            "日本",
				"src_registered_country_code": "JP",
				"src_registered_country_name": "日本",
			},
		},
		{
			desc: "Falls back to English if configured locale doesn't contain a name",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoLite2-City-Test.mmdb"},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
				Locale:        "ja",
			},
			input: map[string]interface{}{
				"src_addr": "2.125.160.216",
			},
			want: map[string]interface{}{
				"src_addr":                    "2.125.160.216",
				"src_city_name":               "Boxford",
				"src_continent_code":          "EU",
				"src_continent_name":          "ヨーロッパ",
				"src_country_code":            "GB",
				"src_country_name":            "イギリス",
				"src_loc_accuracy":            uint64(100),
				"src_loc_lat":                 51.75,
				"src_loc_long":                -1.25,
				"src_loc_postal_code":         "OX1",
				"src_loc_tz":                  "Europe/London",
				"src_registered_country_eu":   true,
				"src_registered_country_code": "FR",
				"src_registered_country_name": "フランス共和国",
			},
		},
		{
			desc: "Only adds enabled fields",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{"./MaxMind-DB/test-data/GeoIP2-Enterprise-Test.mmdb"},
				EnabledFields: []enricher.MaxmindDBEnricherField{
					enricher.MaxmindDBEnricherField_ASN,
					enricher.MaxmindDBEnricherField_LocationLatitude,
					enricher.MaxmindDBEnricherField_LocationLongitude,
					enricher.MaxmindDBEnricherField_IsAnonymousProxy,
				},
			},
			input: map[string]interface{}{
				"src_addr": "2.125.160.216",
			},
			want: map[string]interface{}{
				"src_addr":               "2.125.160.216",
				"src_loc_lat":            51.75,
				"src_loc_long":           -1.25,
				"src_is_anonymous_proxy": false,
			},
		},
		{
			desc: "Merges results from multiple DBs",
			config: enricher.MaxmindDBEnricherConfig{
				DatabasePaths: []string{
					"./MaxMind-DB/test-data/GeoLite2-ASN-Test.mmdb",
					"./MaxMind-DB/test-data/GeoLite2-City-Test.mmdb",
				},
				EnabledFields: enricher.MaxmindDBEnricherFields_All,
			},
			input: map[string]interface{}{
				"src_addr": "89.160.20.112",
			},
			want: map[string]interface{}{
				"src_addr":                    "89.160.20.112",
				"src_city_name":               "Linköping",
				"src_continent_code":          "EU",
				"src_continent_name":          "Europe",
				"src_country_eu":              true,
				"src_country_code":            "SE",
				"src_country_name":            "Sweden",
				"src_loc_accuracy":            uint64(76),
				"src_loc_lat":                 58.4167,
				"src_loc_long":                15.6167,
				"src_loc_tz":                  "Europe/Stockholm",
				"src_registered_country_code": "DE",
				"src_registered_country_eu":   true,
				"src_registered_country_name": "Germany",
				"src_asn":                     uint64(29518),
				"src_asn_org":                 "Bredband2 AB",
			},
		},
	}

	for _, tc := range tests {
		if tc.input == nil {
			t.Logf("\"%s\": skip (unimplmented)", tc.desc)
			continue
		}
		if tc.skip != "" {
			t.Logf("\"%s\": skip (%s)", tc.desc, tc.skip)
			continue
		}
		e := enricher.NewMaxmindDBEnricher(&tc.config)
		got := e.Process(tc.input)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Logf("\"%s\":\n%s", tc.desc, diff)
			t.Fail()
		}
	}
}