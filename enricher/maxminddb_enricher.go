package enricher

import (
	"log"
	"net"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/oschwald/maxminddb-golang"
	"github.com/prometheus/client_golang/prometheus"
)

type MaxmindDBEnricherField string

const (
	MaxmindDBEnricherField_ASN                      MaxmindDBEnricherField = "asn"
	MaxmindDBEnricherField_ASNOrganization          MaxmindDBEnricherField = "asn_org"
	MaxmindDBEnricherField_AverageIncome            MaxmindDBEnricherField = "average_income"
	MaxmindDBEnricherField_CityConfidence           MaxmindDBEnricherField = "city_confidence"
	MaxmindDBEnricherField_CityName                 MaxmindDBEnricherField = "city_name"
	MaxmindDBEnricherField_ConnectionType           MaxmindDBEnricherField = "connection_type"
	MaxmindDBEnricherField_ContinentCode            MaxmindDBEnricherField = "continent_code"
	MaxmindDBEnricherField_ContinentName            MaxmindDBEnricherField = "continent_name"
	MaxmindDBEnricherField_CountryCode              MaxmindDBEnricherField = "country_code"
	MaxmindDBEnricherField_CountryConfidence        MaxmindDBEnricherField = "country_confidence"
	MaxmindDBEnricherField_CountryIsInEU            MaxmindDBEnricherField = "country_eu"
	MaxmindDBEnricherField_CountryName              MaxmindDBEnricherField = "country_name"
	MaxmindDBEnricherField_Domain                   MaxmindDBEnricherField = "domain"
	MaxmindDBEnricherField_IPRisk                   MaxmindDBEnricherField = "ip_risk"
	MaxmindDBEnricherField_IsAnonymous              MaxmindDBEnricherField = "is_anonymous"
	MaxmindDBEnricherField_IsAnonymousProxy         MaxmindDBEnricherField = "is_anonymous_proxy"
	MaxmindDBEnricherField_IsAnonymousVPN           MaxmindDBEnricherField = "is_anonymous_vpn"
	MaxmindDBEnricherField_IsHostingProvider        MaxmindDBEnricherField = "is_hosting_provider"
	MaxmindDBEnricherField_IsLegitimateProxy        MaxmindDBEnricherField = "is_legitimate_proxy"
	MaxmindDBEnricherField_ISP                      MaxmindDBEnricherField = "isp"
	MaxmindDBEnricherField_IsPublicProxy            MaxmindDBEnricherField = "is_public_proxy"
	MaxmindDBEnricherField_IsResidentialProxy       MaxmindDBEnricherField = "is_residential_proxy"
	MaxmindDBEnricherField_IsSatelliteProvider      MaxmindDBEnricherField = "is_satellite_provider"
	MaxmindDBEnricherField_IsTorExitNode            MaxmindDBEnricherField = "is_tor_exit_node"
	MaxmindDBEnricherField_LocationAccuracyRadius   MaxmindDBEnricherField = "loc_accuracy"
	MaxmindDBEnricherField_LocationLatitude         MaxmindDBEnricherField = "loc_lat"
	MaxmindDBEnricherField_LocationLongitude        MaxmindDBEnricherField = "loc_long"
	MaxmindDBEnricherField_LocationMetroCode        MaxmindDBEnricherField = "loc_metro_code"
	MaxmindDBEnricherField_LocationPostalCode       MaxmindDBEnricherField = "loc_postal_code"
	MaxmindDBEnricherField_LocationPostalConfidence MaxmindDBEnricherField = "loc_postal_confidence"
	MaxmindDBEnricherField_LocationTimeZone         MaxmindDBEnricherField = "loc_tz"
	MaxmindDBEnricherField_Organization             MaxmindDBEnricherField = "organization"
	MaxmindDBEnricherField_PopulationDensity        MaxmindDBEnricherField = "population_density"
	MaxmindDBEnricherField_RegisteredCountryCode    MaxmindDBEnricherField = "registered_country_code"
	MaxmindDBEnricherField_RegisteredCountryIsInEU  MaxmindDBEnricherField = "registered_country_eu"
	MaxmindDBEnricherField_RegisteredCountryName    MaxmindDBEnricherField = "registered_country_name"
	MaxmindDBEnricherField_RepresentedCountryCode   MaxmindDBEnricherField = "represented_country_code"
	MaxmindDBEnricherField_RepresentedCountryName   MaxmindDBEnricherField = "represented_country_name"
	MaxmindDBEnricherField_StaticIPScore            MaxmindDBEnricherField = "static_ip_score"
	MaxmindDBEnricherField_UserType                 MaxmindDBEnricherField = "ip_user_type"
)

var (
	MaxmindDBEnricherFields = []MaxmindDBEnricherField{
		MaxmindDBEnricherField_ASN,
		MaxmindDBEnricherField_ASNOrganization,
		MaxmindDBEnricherField_AverageIncome,
		MaxmindDBEnricherField_CityConfidence,
		MaxmindDBEnricherField_CityName,
		MaxmindDBEnricherField_ConnectionType,
		MaxmindDBEnricherField_ContinentCode,
		MaxmindDBEnricherField_ContinentName,
		MaxmindDBEnricherField_CountryCode,
		MaxmindDBEnricherField_CountryConfidence,
		MaxmindDBEnricherField_CountryIsInEU,
		MaxmindDBEnricherField_CountryName,
		MaxmindDBEnricherField_Domain,
		MaxmindDBEnricherField_IPRisk,
		MaxmindDBEnricherField_IsAnonymous,
		MaxmindDBEnricherField_IsAnonymousProxy,
		MaxmindDBEnricherField_IsAnonymousVPN,
		MaxmindDBEnricherField_IsHostingProvider,
		MaxmindDBEnricherField_IsLegitimateProxy,
		MaxmindDBEnricherField_ISP,
		MaxmindDBEnricherField_IsPublicProxy,
		MaxmindDBEnricherField_IsResidentialProxy,
		MaxmindDBEnricherField_IsSatelliteProvider,
		MaxmindDBEnricherField_IsTorExitNode,
		MaxmindDBEnricherField_LocationAccuracyRadius,
		MaxmindDBEnricherField_LocationLatitude,
		MaxmindDBEnricherField_LocationLongitude,
		MaxmindDBEnricherField_LocationMetroCode,
		MaxmindDBEnricherField_LocationPostalCode,
		MaxmindDBEnricherField_LocationPostalConfidence,
		MaxmindDBEnricherField_LocationTimeZone,
		MaxmindDBEnricherField_Organization,
		MaxmindDBEnricherField_PopulationDensity,
		MaxmindDBEnricherField_RegisteredCountryCode,
		MaxmindDBEnricherField_RegisteredCountryIsInEU,
		MaxmindDBEnricherField_RegisteredCountryName,
		MaxmindDBEnricherField_RepresentedCountryCode,
		MaxmindDBEnricherField_RepresentedCountryName,
		MaxmindDBEnricherField_StaticIPScore,
		MaxmindDBEnricherField_UserType,
	}
	MaxmindDBEnricherFields_All           = MaxmindDBEnricherFields
	MaxmindDBEnricherFields_MaximumMisery = MaxmindDBEnricherFields

	MaxmindDBEnricherFields_AnonymousIP = []MaxmindDBEnricherField{
		MaxmindDBEnricherField_IsAnonymous,
		MaxmindDBEnricherField_IsAnonymousVPN,
		MaxmindDBEnricherField_IsTorExitNode,
		MaxmindDBEnricherField_IsHostingProvider,
		MaxmindDBEnricherField_IsPublicProxy,
		MaxmindDBEnricherField_IsResidentialProxy,
	}

	MaxmindDBEnricherFields_ASN = []MaxmindDBEnricherField{
		MaxmindDBEnricherField_ASN,
		MaxmindDBEnricherField_ASNOrganization,
	}

	MaxmindDBEnricherFields_City = []MaxmindDBEnricherField{
		MaxmindDBEnricherField_ContinentCode,
		MaxmindDBEnricherField_ContinentName,
		MaxmindDBEnricherField_CountryCode,
		MaxmindDBEnricherField_CountryName,
		MaxmindDBEnricherField_LocationAccuracyRadius,
		MaxmindDBEnricherField_LocationLatitude,
		MaxmindDBEnricherField_LocationLongitude,
		MaxmindDBEnricherField_LocationPostalCode,
		MaxmindDBEnricherField_LocationTimeZone,
		MaxmindDBEnricherField_RegisteredCountryCode,
		MaxmindDBEnricherField_RegisteredCountryName,
		MaxmindDBEnricherField_RepresentedCountryCode,
		MaxmindDBEnricherField_RepresentedCountryName,
	}

	MaxmindDBEnricherFields_ConnectionType = []MaxmindDBEnricherField{
		MaxmindDBEnricherField_ConnectionType,
	}

	MaxmindDBEnricherFields_Country = []MaxmindDBEnricherField{
		MaxmindDBEnricherField_ContinentCode,
		MaxmindDBEnricherField_ContinentName,
		MaxmindDBEnricherField_CountryCode,
		MaxmindDBEnricherField_CountryName,
		MaxmindDBEnricherField_RegisteredCountryCode,
		MaxmindDBEnricherField_RegisteredCountryName,
		MaxmindDBEnricherField_RepresentedCountryCode,
		MaxmindDBEnricherField_RepresentedCountryName,
	}

	MaxmindDBEnricherFields_DensityIncome = []MaxmindDBEnricherField{
		MaxmindDBEnricherField_AverageIncome,
		MaxmindDBEnricherField_PopulationDensity,
	}

	MaxmindDBEnricherFields_Domain = []MaxmindDBEnricherField{
		MaxmindDBEnricherField_Domain,
	}

	MaxmindDBEnricherFields_Enterprise = []MaxmindDBEnricherField{
		MaxmindDBEnricherField_CityConfidence,
		MaxmindDBEnricherField_CityName,
		MaxmindDBEnricherField_ContinentCode,
		MaxmindDBEnricherField_ContinentName,
		MaxmindDBEnricherField_CountryConfidence,
		MaxmindDBEnricherField_CountryCode,
		MaxmindDBEnricherField_CountryName,
		MaxmindDBEnricherField_CountryIsInEU,
		MaxmindDBEnricherField_LocationAccuracyRadius,
		MaxmindDBEnricherField_LocationLatitude,
		MaxmindDBEnricherField_LocationLongitude,
		MaxmindDBEnricherField_LocationMetroCode,
		MaxmindDBEnricherField_LocationTimeZone,
		MaxmindDBEnricherField_LocationPostalCode,
		MaxmindDBEnricherField_LocationPostalConfidence,
		MaxmindDBEnricherField_RegisteredCountryCode,
		MaxmindDBEnricherField_RegisteredCountryIsInEU,
		MaxmindDBEnricherField_RegisteredCountryName,
		MaxmindDBEnricherField_ASN,
		MaxmindDBEnricherField_ASNOrganization,
		MaxmindDBEnricherField_ConnectionType,
		MaxmindDBEnricherField_Domain,
		MaxmindDBEnricherField_IsAnonymousProxy,
		MaxmindDBEnricherField_IsLegitimateProxy,
		MaxmindDBEnricherField_IsSatelliteProvider,
		MaxmindDBEnricherField_ISP,
		MaxmindDBEnricherField_Organization,
		MaxmindDBEnricherField_StaticIPScore,
		MaxmindDBEnricherField_UserType,
	}

	MaxmindDBEnricherFields_IPRisk = []MaxmindDBEnricherField{
		MaxmindDBEnricherField_IPRisk,
		MaxmindDBEnricherField_IsAnonymous,
		MaxmindDBEnricherField_IsAnonymousVPN,
		MaxmindDBEnricherField_IsHostingProvider,
		MaxmindDBEnricherField_IsPublicProxy,
		MaxmindDBEnricherField_IsResidentialProxy,
		MaxmindDBEnricherField_IsTorExitNode,
	}

	MaxmindDBEnricherFields_ISP = []MaxmindDBEnricherField{
		MaxmindDBEnricherField_ASN,
		MaxmindDBEnricherField_ASNOrganization,
		MaxmindDBEnricherField_ISP,
		MaxmindDBEnricherField_Organization,
	}

	MaxmindDBEnricherFields_StaticIPScore = []MaxmindDBEnricherField{
		MaxmindDBEnricherField_StaticIPScore,
	}
)

var (
	MetricMMDBCacheSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "mmdb_cache_size",
			Help: "size of MaxMind DB enricher LRU cache",
		},
	)
	MetricMMDBCacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mmdb_cache_hits",
			Help: "Number of MaxMind DB enricher LRU cache hits",
		},
	)
	MetricMMDBCacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mmdb_cache_misses",
			Help: "Number of MaxMind DB enricher LRU cache misses",
		},
	)
	MetricMMDBLookups = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mmdb_lookups",
			Help: "Number of MaxMind DB enricher lookups",
		},
		[]string{"dbpath", "status"},
	)
)

func init() {
	prometheus.MustRegister(MetricMMDBCacheSize)
	prometheus.MustRegister(MetricMMDBCacheHits)
	prometheus.MustRegister(MetricMMDBCacheMisses)
	prometheus.MustRegister(MetricMMDBLookups)
}

type MaxmindDBEnricherConfig struct {
	EnableCache   bool
	CacheSize     int
	CacheOnly     bool
	Locale        string
	DatabasePaths []string
	EnabledFields []MaxmindDBEnricherField
}

type MaxmindDBEnricherIPData map[MaxmindDBEnricherField]interface{}

type MaxmindDBEnricher struct {
	Config            *MaxmindDBEnricherConfig
	cache             *lru.Cache[string, MaxmindDBEnricherIPData]
	cacheLookupStatus sync.Map
	readers           []*maxminddb.Reader
}

func NewMaxmindDBEnricher(config *MaxmindDBEnricherConfig) MaxmindDBEnricher {
	if config == nil {
		config = &MaxmindDBEnricherConfig{}
	}
	if config.EnableCache && config.CacheSize == 0 {
		config.CacheSize = 128
	}
	if config.Locale == "" {
		config.Locale = "en"
	}
	if config.DatabasePaths == nil {
		config.DatabasePaths = make([]string, 0)
	}
	if config.EnabledFields == nil {
		fields := make([]MaxmindDBEnricherField, 1)
		fields = append(fields, MaxmindDBEnricherFields_ASN...)
		fields = append(fields, MaxmindDBEnricherFields_City...)
		config.EnabledFields = fields
	}

	var cache *lru.Cache[string, MaxmindDBEnricherIPData]
	var err error
	if config.EnableCache {
		cache, err = lru.New[string, MaxmindDBEnricherIPData](config.CacheSize)
		if err != nil {
			panic(err)
		}
	}

	readers := make([]*maxminddb.Reader, 0)
	for _, dbPath := range config.DatabasePaths {
		reader, err := maxminddb.Open(dbPath)
		if err != nil {
			panic(err)
		}
		readers = append(readers, reader)
	}

	e := MaxmindDBEnricher{
		Config:  config,
		readers: readers,
		cache:   cache,
	}

	return e
}

func (e *MaxmindDBEnricher) Process(msg map[string]interface{}) map[string]interface{} {
	defer func() {
		if e.Config.EnableCache {
			MetricMMDBCacheSize.Set(float64(e.cache.Len()))
		}
	}()
	msg = e.add(msg, "src_addr", "src_")
	msg = e.add(msg, "dst_addr", "dst_")
	msg = e.add(msg, "src_addr_encap", "src_encap_")
	msg = e.add(msg, "dst_addr_encap", "dst_encap_")
	return msg
}

func (e *MaxmindDBEnricher) add(msg map[string]interface{}, originalField string, targetPrefix string) map[string]interface{} {
	addrRaw, ok := msg[originalField]
	if !ok {
		return msg
	}
	addr := addrRaw.(string)
	var data MaxmindDBEnricherIPData

	if e.Config.EnableCache {
		data, ok := e.cache.Get(addr)
		if ok {
			MetricMMDBCacheHits.Inc()
			msg = e.mergeDataIntoMessage(msg, data, targetPrefix)
			return msg
		}
		MetricMMDBCacheMisses.Inc()
	}

	if e.Config.CacheOnly {
		go func(addr string) {
			_, inProgress := e.cacheLookupStatus.Load(addr)
			if !inProgress {
				e.cacheLookupStatus.Store(addr, true)
				e.resolveIP(addr)
				e.cacheLookupStatus.Delete(addr)
			}
		}(addr)
		return msg
	}

	data = e.resolveIP(addr)
	msg = e.mergeDataIntoMessage(msg, data, targetPrefix)
	return msg
}

func (e *MaxmindDBEnricher) mergeDataIntoMessage(msg map[string]interface{}, data MaxmindDBEnricherIPData, prefix string) map[string]interface{} {
	for key, value := range data {
		msg[prefix+string(key)] = value
	}
	return msg
}

func (e *MaxmindDBEnricher) mergeData(d ...MaxmindDBEnricherIPData) MaxmindDBEnricherIPData {
	result := make(MaxmindDBEnricherIPData)
	for _, data := range d {
		for key, value := range data {
			result[key] = value
		}
	}
	return result
}

func (e *MaxmindDBEnricher) isFieldEnabled(field MaxmindDBEnricherField) bool {
	for _, f := range e.Config.EnabledFields {
		if f == field {
			return true
		}
	}
	return false
}

func (e *MaxmindDBEnricher) resolveIP(ip string) MaxmindDBEnricherIPData {
	result := make(MaxmindDBEnricherIPData)

	if ip == "" {
		return result
	}

	for i, reader := range e.readers {
		dbPath := e.Config.DatabasePaths[i]
		var record map[string]interface{}
		err := reader.Lookup(net.ParseIP(ip), &record)
		if err != nil {
			log.Printf("error resolving IP %s with MaxMind DB %s: %v", ip, dbPath, err)
			MetricMMDBLookups.With(prometheus.Labels{"dbpath": dbPath, "status": "error"}).Inc()
			continue
		} else if len(record) == 0 {
			MetricMMDBLookups.With(prometheus.Labels{"dbpath": dbPath, "status": "empty"}).Inc()
		} else {
			MetricMMDBLookups.With(prometheus.Labels{"dbpath": dbPath, "status": "success"}).Inc()
		}
		result = e.mergeData(result, e.flattenData(reader, record))
	}

	if e.Config.EnableCache {
		e.cache.Add(ip, result)
	}
	return result
}

func (e MaxmindDBEnricher) localizedName(v interface{}) string {
	names := v.(map[string]interface{})
	if name, ok := names[e.Config.Locale]; ok {
		if name, ok := name.(string); ok {
			return name
		}
	}
	if name, ok := names["en"]; ok {
		if name, ok := name.(string); ok {
			return name
		}
	}
	for _, name := range names {
		if name, ok := name.(string); ok {
			return name
		}
	}
	return ""
}

func (e *MaxmindDBEnricher) flattenData(reader *maxminddb.Reader, raw map[string]interface{}) MaxmindDBEnricherIPData {
	result := make(MaxmindDBEnricherIPData)

	result = e.mergeData(
		result,
		e.flattenTraitsRaw(reader, raw),
		e.flattenCity(reader, raw),
		e.flattenContinent(reader, raw),
		e.flattenCountry(reader, raw),
		e.flattenRegisteredCountry(reader, raw),
		e.flattenRepresentedCountry(reader, raw),
		e.flattenLocation(reader, raw),
		e.flattenTraits(reader, raw),
	)

	return result
}

func (e *MaxmindDBEnricher) flattenCity(reader *maxminddb.Reader, raw map[string]interface{}) MaxmindDBEnricherIPData {
	result := make(MaxmindDBEnricherIPData)
	if city, ok := raw["city"]; ok {
		if city, ok := city.(map[string]interface{}); ok {
			if val, ok := city["confidence"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_CityConfidence) {
				result[MaxmindDBEnricherField_CityConfidence] = val
			}
			if val, ok := city["names"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_CityName) {
				result[MaxmindDBEnricherField_CityName] = e.localizedName(val)
			}
		}
	}
	return result
}

func (e *MaxmindDBEnricher) flattenLocation(reader *maxminddb.Reader, raw map[string]interface{}) MaxmindDBEnricherIPData {
	result := make(MaxmindDBEnricherIPData)
	if location, ok := raw["location"]; ok {
		if location, ok := location.(map[string]interface{}); ok {
			if val, ok := location["accuracy_radius"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_LocationAccuracyRadius) {
				result[MaxmindDBEnricherField_LocationAccuracyRadius] = val
			}
			if val, ok := location["latitude"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_LocationLatitude) {
				result[MaxmindDBEnricherField_LocationLatitude] = val
			}
			if val, ok := location["longitude"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_LocationLongitude) {
				result[MaxmindDBEnricherField_LocationLongitude] = val
			}
			if val, ok := location["metro_code"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_LocationMetroCode) {
				result[MaxmindDBEnricherField_LocationMetroCode] = val
			}
			if val, ok := location["time_zone"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_LocationTimeZone) {
				result[MaxmindDBEnricherField_LocationTimeZone] = val
			}
		}
	}
	if postal, ok := raw["postal"]; ok {
		if postal, ok := postal.(map[string]interface{}); ok {
			if val, ok := postal["code"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_LocationPostalCode) {
				result[MaxmindDBEnricherField_LocationPostalCode] = val
			}
			if val, ok := postal["confidence"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_LocationPostalConfidence) {
				result[MaxmindDBEnricherField_LocationPostalConfidence] = val
			}
		}
	}
	return result
}

func (e *MaxmindDBEnricher) flattenContinent(reader *maxminddb.Reader, raw map[string]interface{}) MaxmindDBEnricherIPData {
	result := make(MaxmindDBEnricherIPData)
	if continent, ok := raw["continent"]; ok {
		if continent, ok := continent.(map[string]interface{}); ok {
			if val, ok := continent["code"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_ContinentCode) {
				result[MaxmindDBEnricherField_ContinentCode] = val
			}
			if val, ok := continent["names"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_ContinentName) {
				result[MaxmindDBEnricherField_ContinentName] = e.localizedName(val)
			}
		}
	}
	return result
}

func (e *MaxmindDBEnricher) flattenCountry(reader *maxminddb.Reader, raw map[string]interface{}) MaxmindDBEnricherIPData {
	result := make(MaxmindDBEnricherIPData)
	if country, ok := raw["country"]; ok {
		if country, ok := country.(map[string]interface{}); ok {
			if val, ok := country["code"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_CountryCode) {
				result[MaxmindDBEnricherField_CountryCode] = val
			}
			if val, ok := country["iso_code"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_CountryCode) {
				result[MaxmindDBEnricherField_CountryCode] = val
			}
			if val, ok := country["is_in_european_union"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_CountryIsInEU) {
				result[MaxmindDBEnricherField_CountryIsInEU] = val
			}
			if val, ok := country["names"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_CountryName) {
				result[MaxmindDBEnricherField_CountryName] = e.localizedName(val)
			}
			if val, ok := country["confidence"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_CountryConfidence) {
				result[MaxmindDBEnricherField_CountryConfidence] = val
			}
		}
	}
	return result
}

func (e *MaxmindDBEnricher) flattenRegisteredCountry(reader *maxminddb.Reader, raw map[string]interface{}) MaxmindDBEnricherIPData {
	result := make(MaxmindDBEnricherIPData)
	if country, ok := raw["registered_country"]; ok {
		if country, ok := country.(map[string]interface{}); ok {
			if val, ok := country["code"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_RegisteredCountryCode) {
				result[MaxmindDBEnricherField_RegisteredCountryCode] = val
			}
			if val, ok := country["iso_code"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_RegisteredCountryCode) {
				result[MaxmindDBEnricherField_RegisteredCountryCode] = val
			}
			if val, ok := country["is_in_european_union"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_RegisteredCountryIsInEU) {
				result[MaxmindDBEnricherField_RegisteredCountryIsInEU] = val
			}
			if val, ok := country["names"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_RegisteredCountryName) {
				result[MaxmindDBEnricherField_RegisteredCountryName] = e.localizedName(val)
			}
		}
	}
	return result
}

func (e *MaxmindDBEnricher) flattenRepresentedCountry(reader *maxminddb.Reader, raw map[string]interface{}) MaxmindDBEnricherIPData {
	result := make(MaxmindDBEnricherIPData)
	if country, ok := raw["represented_country"]; ok {
		if country, ok := country.(map[string]interface{}); ok {
			if val, ok := country["code"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_RepresentedCountryCode) {
				result[MaxmindDBEnricherField_RepresentedCountryCode] = val
			}
			if val, ok := country["iso_code"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_RepresentedCountryCode) {
				result[MaxmindDBEnricherField_RepresentedCountryCode] = val
			}
			if val, ok := country["names"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_RepresentedCountryName) {
				result[MaxmindDBEnricherField_RepresentedCountryName] = e.localizedName(val)
			}
		}
	}
	return result
}

func (e *MaxmindDBEnricher) flattenTraits(reader *maxminddb.Reader, raw map[string]interface{}) MaxmindDBEnricherIPData {
	result := make(MaxmindDBEnricherIPData)
	if traits, ok := raw["traits"]; ok {
		if traits, ok := traits.(map[string]interface{}); ok {
			result = e.flattenTraitsRaw(reader, traits)
		}
	}
	return result
}

func (e *MaxmindDBEnricher) flattenTraitsRaw(reader *maxminddb.Reader, traits map[string]interface{}) MaxmindDBEnricherIPData {
	result := make(MaxmindDBEnricherIPData)
	if val, ok := traits["autonomous_system_number"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_ASN) {
		result[MaxmindDBEnricherField_ASN] = val
	}
	if val, ok := traits["autonomous_system_organization"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_ASNOrganization) {
		result[MaxmindDBEnricherField_ASNOrganization] = val
	}
	if val, ok := traits["average_income"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_AverageIncome) {
		result[MaxmindDBEnricherField_AverageIncome] = val
	}
	if val, ok := traits["connection_type"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_ConnectionType) {
		result[MaxmindDBEnricherField_ConnectionType] = val
	}
	if val, ok := traits["domain"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_Domain) {
		result[MaxmindDBEnricherField_Domain] = val
	}
	if val, ok := traits["ip_risk"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_IPRisk) {
		result[MaxmindDBEnricherField_IPRisk] = val
	}
	if val, ok := traits["isp"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_ISP) {
		result[MaxmindDBEnricherField_ISP] = val
	}
	if val, ok := traits["organization"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_Organization) {
		result[MaxmindDBEnricherField_Organization] = val
	}
	if val, ok := traits["population_density"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_PopulationDensity) {
		result[MaxmindDBEnricherField_PopulationDensity] = val
	}
	if val, ok := traits["static_ip_score"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_StaticIPScore) {
		result[MaxmindDBEnricherField_StaticIPScore] = val
	}
	if reader.Metadata.DatabaseType == "GeoIP2-Static-IP-Score" {
		if val, ok := traits["score"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_StaticIPScore) {
			result[MaxmindDBEnricherField_StaticIPScore] = val
		}
	}
	if val, ok := traits["user_type"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_UserType) {
		result[MaxmindDBEnricherField_UserType] = val
	}

	// In MMDB, if a flag is not present it is presumed to be false. This is an
	// attempt to set a default value (false) for those flags depending on the
	// database type. Then later on the flag can be set if it is true.
	switch reader.Metadata.DatabaseType {
	case "GeoIP2-Anonymous-IP":
		if e.isFieldEnabled(MaxmindDBEnricherField_IsAnonymous) {
			result[MaxmindDBEnricherField_IsAnonymous] = false
		}
		if e.isFieldEnabled(MaxmindDBEnricherField_IsAnonymousVPN) {
			result[MaxmindDBEnricherField_IsAnonymousVPN] = false
		}
		if e.isFieldEnabled(MaxmindDBEnricherField_IsPublicProxy) {
			result[MaxmindDBEnricherField_IsPublicProxy] = false
		}
		if e.isFieldEnabled(MaxmindDBEnricherField_IsResidentialProxy) {
			result[MaxmindDBEnricherField_IsResidentialProxy] = false
		}
		if e.isFieldEnabled(MaxmindDBEnricherField_IsTorExitNode) {
			result[MaxmindDBEnricherField_IsTorExitNode] = false
		}
	case "DBIP-ISP (compat=Enterprise)",
		"DBIP-Location-ISP (compat=Enterprise)",
		"GeoIP2-Enterprise",
		"GeoIP2-Precision-Enterprise",
		"GeoIP2-Precision-Enterprise-Sandbox":
		if e.isFieldEnabled(MaxmindDBEnricherField_IsAnonymousProxy) {
			result[MaxmindDBEnricherField_IsAnonymousProxy] = false
		}
		if e.isFieldEnabled(MaxmindDBEnricherField_IsLegitimateProxy) {
			result[MaxmindDBEnricherField_IsLegitimateProxy] = false
		}
		if e.isFieldEnabled(MaxmindDBEnricherField_IsSatelliteProvider) {
			result[MaxmindDBEnricherField_IsSatelliteProvider] = false
		}
	case "GeoIP2-IP-Risk":
		if e.isFieldEnabled(MaxmindDBEnricherField_IsAnonymous) {
			result[MaxmindDBEnricherField_IsAnonymous] = false
		}
		if e.isFieldEnabled(MaxmindDBEnricherField_IsAnonymousVPN) {
			result[MaxmindDBEnricherField_IsAnonymousVPN] = false
		}
		if e.isFieldEnabled(MaxmindDBEnricherField_IsHostingProvider) {
			result[MaxmindDBEnricherField_IsHostingProvider] = false
		}
		if e.isFieldEnabled(MaxmindDBEnricherField_IsPublicProxy) {
			result[MaxmindDBEnricherField_IsPublicProxy] = false
		}
		if e.isFieldEnabled(MaxmindDBEnricherField_IsResidentialProxy) {
			result[MaxmindDBEnricherField_IsResidentialProxy] = false
		}
		if e.isFieldEnabled(MaxmindDBEnricherField_IsTorExitNode) {
			result[MaxmindDBEnricherField_IsTorExitNode] = false
		}
	}

	if val, ok := traits["is_anonymous"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_IsAnonymous) {
		result[MaxmindDBEnricherField_IsAnonymous] = val
	}
	if val, ok := traits["is_anonymous_proxy"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_IsAnonymous) {
		result[MaxmindDBEnricherField_IsAnonymous] = val
	}
	if val, ok := traits["is_anonymous_vpn"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_IsAnonymousVPN) {
		result[MaxmindDBEnricherField_IsAnonymousVPN] = val
	}
	if val, ok := traits["is_hosting_provider"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_IsHostingProvider) {
		result[MaxmindDBEnricherField_IsHostingProvider] = val
	}
	if val, ok := traits["is_legitimate_proxy"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_IsLegitimateProxy) {
		result[MaxmindDBEnricherField_IsLegitimateProxy] = val
	}
	if val, ok := traits["is_public_proxy"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_IsPublicProxy) {
		result[MaxmindDBEnricherField_IsPublicProxy] = val
	}
	if val, ok := traits["is_residential_proxy"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_IsResidentialProxy) {
		result[MaxmindDBEnricherField_IsResidentialProxy] = val
	}
	if val, ok := traits["is_satellite_provider"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_IsSatelliteProvider) {
		result[MaxmindDBEnricherField_IsSatelliteProvider] = val
	}
	if val, ok := traits["is_tor_exit_node"]; ok && e.isFieldEnabled(MaxmindDBEnricherField_IsTorExitNode) {
		result[MaxmindDBEnricherField_IsTorExitNode] = val
	}
	return result
}
