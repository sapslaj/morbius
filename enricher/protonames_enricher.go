package enricher

import "log"

type ProtonamesEnricherConfig struct {
}

type ProtonamesEnricher struct {
	Config         *ProtonamesEnricherConfig
	protoTable     map[int]string
	ethertypeTable map[int]string
}

func NewProtonamesEnricher(config *ProtonamesEnricherConfig) ProtonamesEnricher {
	if config == nil {
		config = &ProtonamesEnricherConfig{}
	}
	log.Println("[WARN] ProtonamesEnricher is deprecated. Use NetDBEnricher instead.")
	return ProtonamesEnricher{
		Config: config,
		protoTable: map[int]string{
			0:   "HOPOPT",
			1:   "ICMP",
			2:   "IGMP",
			3:   "GGP",
			4:   "IP-ENCAP",
			5:   "ST",
			6:   "TCP",
			8:   "EGP",
			9:   "IGP",
			12:  "PUP",
			17:  "UDP",
			20:  "HMP",
			22:  "XNS-IDP",
			27:  "RDP",
			29:  "ISO-TP4",
			33:  "DCCP",
			36:  "XTP",
			37:  "DDP",
			38:  "IDPR-CMTP",
			41:  "IPv6",
			43:  "IPv6-Route",
			44:  "IPv6-Frag",
			45:  "IDRP",
			46:  "RSVP",
			47:  "GRE",
			50:  "IPSEC-ESP",
			51:  "IPSEC-AH",
			57:  "SKIP",
			58:  "IPv6-ICMP",
			59:  "IPv6-NoNxt",
			60:  "IPv6-Opts",
			73:  "RSPF",
			81:  "VMTP",
			88:  "EIGRP",
			89:  "OSPFIGP",
			93:  "AX.25",
			94:  "IPIP",
			97:  "ETHERIP",
			98:  "ENCAP",
			103: "PIM",
			108: "IPCOMP",
			112: "VRRP",
			115: "L2TP",
			124: "ISIS",
			132: "SCTP",
			133: "FC",
			135: "Mobility-Header",
			136: "UDPLite",
			137: "MPLS-in-IP",
			138: "MANET",
			139: "HIP",
			140: "Shim6",
			141: "WESP",
			142: "ROHC",
		},
		ethertypeTable: map[int]string{
			0x0800: "IPv4",
			0x0806: "ARP",
			0x0842: "Wake-on-LAN",
			0x22F0: "AVTP",
			0x22F3: "TRILL",
			0x22EA: "SRP",
			0x8035: "RARP",
			0x809B: "AppleTalk",
			0x80F3: "AppleTalk AARP",
			0x8100: "C-Tag",
			0x8102: "SLPP",
			0x8103: "VLACP",
			0x8137: "IPX",
			0x8204: "QNX Qnet",
			0x86DD: "IPv6",
			0x8808: "EPON",
			0x8809: "LACP",
			0x8819: "CobraNet",
			0x8847: "MPLS unicast",
			0x8848: "MPLS multicast",
			0x8863: "PPPoE Discovery Stage",
			0x8864: "PPPoE Session Stage",
			0x888E: "802.1X",
			0x88A8: "S-Tag",
			0x88BF: "MikroTik RoMON",
			0x88CC: "LLDP",
			0x88E5: "MACsec",
			0x88E7: "PBB",
			0x9000: "Loopback",
		},
	}
}

func (e *ProtonamesEnricher) Process(msg map[string]interface{}) map[string]interface{} {
	msg = e.add(msg, "proto", "protocol_name", func(original interface{}) (result interface{}, ok bool) {
		result, ok = e.protoTable[original.(int)]
		return
	})
	msg = e.add(msg, "proto_encap", "protocol_encap_name", func(original interface{}) (result interface{}, ok bool) {
		result, ok = e.protoTable[original.(int)]
		return
	})
	msg = e.add(msg, "ethernet_type", "ethernet_type_name", func(original interface{}) (result interface{}, ok bool) {
		result, ok = e.ethertypeTable[original.(int)]
		return
	})
	msg = e.add(msg, "ethernet_type_encap", "ethernet_type_encap_name", func(original interface{}) (result interface{}, ok bool) {
		result, ok = e.ethertypeTable[original.(int)]
		return
	})
	return msg
}

func (e *ProtonamesEnricher) add(msg map[string]interface{}, originalField string, targetField string, extract func(interface{}) (interface{}, bool)) map[string]interface{} {
	original, ok := msg[originalField]
	if !ok {
		return msg
	}
	result, ok := extract(original)
	if !ok {
		return msg
	}
	msg[targetField] = result
	return msg
}
