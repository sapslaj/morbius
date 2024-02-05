package enricher

import (
	"strings"

	"github.com/thediveo/netdb"
)

// HACK: upstream doesn't support all of the EtherTypes we want.
var NetDBEnricherBuiltinEtherTypes = []netdb.EtherType{
	{
		Name:   "cobranet",
		Number: 0x8819,
	},
	{
		Name:   "mikrotik-romon",
		Number: 0x88bf,
	},
	{
		Name:   "avtp",
		Number: 0x22f0,
	},
	{
		Name:   "vlacp",
		Number: 0x8103,
	},
	{
		Name:   "lacp",
		Number: 0x8809,
	},
	{
		Name:   "wake-on-lan",
		Number: 0x0842,
	},
	{
		Name:   "srp",
		Number: 0x22ea,
	},
	{
		Name:   "qnx-qnet",
		Number: 0x8204,
	},
	{
		Name:   "loopback",
		Number: 0x9000,
	},
	{
		Name:   "slpp",
		Number: 0x8102,
	},
	{
		Name:   "epon",
		Number: 0x8808,
	},
}

type NetDBEnricherConfigConfig struct {
	BuiltIn      bool              `yaml:"built_in"`
	SourceFiles  []string          `yaml:"source_files"`
	SourceInline []string          `yaml:"source_inline"`
	NameAliases  map[string]string `yaml:"name_aliases"`
}

type NetDBEnricherConfig struct {
	EtherTypes *NetDBEnricherConfigConfig `yaml:"ethertypes"`
	Protocols  *NetDBEnricherConfigConfig `yaml:"protocols"`
	Services   *NetDBEnricherConfigConfig `yaml:"services"`
}

type NetDBEnricher struct {
	Config         *NetDBEnricherConfig
	etherTypeIndex *netdb.EtherTypeIndex
	protocolIndex  *netdb.ProtocolIndex
	serviceIndex   *netdb.ServiceIndex
}

func NewNetDBEnricher(config *NetDBEnricherConfig) NetDBEnricher {
	if config == nil {
		config = &NetDBEnricherConfig{}
	}
	if config.EtherTypes == nil {
		config.EtherTypes = &NetDBEnricherConfigConfig{
			BuiltIn: true,
		}
	}
	if config.Protocols == nil {
		config.Protocols = &NetDBEnricherConfigConfig{
			BuiltIn: true,
		}
	}
	if config.Services == nil {
		config.Services = &NetDBEnricherConfigConfig{
			BuiltIn: true,
		}
	}
	etherTypeIndex := netdb.NewEtherTypeIndex([]netdb.EtherType{})
	if config.EtherTypes.BuiltIn {
		etherTypeIndex.Merge(NetDBEnricherBuiltinEtherTypes)
		etherTypeIndex.Merge(netdb.BuiltinEtherTypes)
	}
	for _, sourceFile := range config.EtherTypes.SourceFiles {
		eti, err := netdb.LoadEtherTypes(sourceFile)
		if err != nil {
			panic(err)
		}
		etherTypeIndex.MergeIndex(eti)
	}
	for _, inline := range config.EtherTypes.SourceInline {
		buf := strings.NewReader(inline)
		etherTypes, err := netdb.ParseEtherTypes(buf)
		if err != nil {
			panic(err)
		}
		etherTypeIndex.Merge(etherTypes)
	}
	protocolIndex := netdb.NewProtocolIndex([]netdb.Protocol{})
	if config.Protocols.BuiltIn {
		protocolIndex.Merge(netdb.BuiltinProtocols)
	}
	for _, sourceFile := range config.Protocols.SourceFiles {
		pi, err := netdb.LoadProtocols(sourceFile)
		if err != nil {
			panic(err)
		}
		protocolIndex.MergeIndex(pi)
	}
	for _, inline := range config.Protocols.SourceInline {
		buf := strings.NewReader(inline)
		protocols, err := netdb.ParseProtocols(buf)
		if err != nil {
			panic(err)
		}
		protocolIndex.Merge(protocols)
	}
	serviceIndex := netdb.NewServiceIndex([]netdb.Service{})
	if config.Services.BuiltIn {
		serviceIndex.Merge(netdb.BuiltinServices)
	}
	for _, sourceFile := range config.Services.SourceFiles {
		si, err := netdb.LoadServices(sourceFile, protocolIndex)
		if err != nil {
			panic(err)
		}
		serviceIndex.MergeIndex(si)
	}
	for _, inline := range config.Services.SourceInline {
		buf := strings.NewReader(inline)
		services, err := netdb.ParseServices(buf, protocolIndex)
		if err != nil {
			panic(err)
		}
		serviceIndex.Merge(services)
	}
	e := NetDBEnricher{
		Config:         config,
		etherTypeIndex: &etherTypeIndex,
		protocolIndex:  &protocolIndex,
		serviceIndex:   &serviceIndex,
	}
	return e
}

func (e *NetDBEnricher) maybeAliased(config *NetDBEnricherConfigConfig, name string, names ...string) string {
	if config.NameAliases == nil {
		return name
	}
	names = append([]string{name}, names...)
	for _, n := range names {
		if alias, ok := config.NameAliases[n]; ok {
			return alias
		}
	}
	return name
}

func (e *NetDBEnricher) Process(msg map[string]interface{}) map[string]interface{} {
	var refService *netdb.Service
	if proto, ok := msg["proto"].(int); ok {
		protocolNumber := uint8(proto)
		if protocol, ok := e.protocolIndex.Numbers[protocolNumber]; ok {
			msg["protocol_name"] = e.maybeAliased(e.Config.Protocols, protocol.Name, protocol.Aliases...)
			if srcPort, ok := msg["src_port"].(int); ok {
				srcService := e.serviceIndex.ByPort(srcPort, protocol.Name)
				if srcService != nil {
					msg["src_service_name"] = e.maybeAliased(e.Config.Services, srcService.Name, srcService.Aliases...)
					refService = srcService
				}
			}
			if dstPort, ok := msg["dst_port"].(int); ok {
				dstService := e.serviceIndex.ByPort(dstPort, protocol.Name)
				if dstService != nil {
					msg["dst_service_name"] = e.maybeAliased(e.Config.Services, dstService.Name, dstService.Aliases...)
					if refService == nil || refService.Port > dstService.Port {
						refService = dstService
					}
				}
			}
		}
	}

	if refService != nil {
		msg["service_name"] = e.maybeAliased(e.Config.Services, refService.Name, refService.Aliases...)
	}

	if protoEncap, ok := msg["proto_encap"].(int); ok {
		if protocol, ok := e.protocolIndex.Numbers[uint8(protoEncap)]; ok {
			msg["protocol_encap_name"] = e.maybeAliased(e.Config.Protocols, protocol.Name, protocol.Aliases...)
		}
	}

	if etype, ok := msg["ethernet_type"].(int); ok {
		etherTypeNumber := uint16(etype)
		if etherType, ok := e.etherTypeIndex.Numbers[etherTypeNumber]; ok {
			msg["ethernet_type_name"] = e.maybeAliased(e.Config.EtherTypes, etherType.Name, etherType.Aliases...)
		}
	}

	if etype, ok := msg["ethernet_type_encap"].(int); ok {
		etherTypeNumber := uint16(etype)
		if etherType, ok := e.etherTypeIndex.Numbers[etherTypeNumber]; ok {
			msg["ethernet_type_encap_name"] = e.maybeAliased(e.Config.EtherTypes, etherType.Name, etherType.Aliases...)
		}
	}

	return msg
}
