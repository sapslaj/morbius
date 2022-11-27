package transport

import (
	"encoding/binary"
	"net"

	goflowpb "github.com/cloudflare/goflow/v3/pb"
	"github.com/sapslaj/morbius/destination"
	"github.com/sapslaj/morbius/enricher"
)

type Transport struct {
	Destinations []destination.Destination
	Enrichers    []enricher.Enricher
}

func (s *Transport) Publish(fmsgs []*goflowpb.FlowMessage) {
	for _, fmsg := range fmsgs {
		ffmsg := s.FormatFlowMessage(fmsg)
		for _, enricher := range s.Enrichers {
			ffmsg = enricher.Process(ffmsg)
		}
		for _, destination := range s.Destinations {
			destination.Publish(ffmsg)
		}
	}
}

func (s *Transport) FormatFlowMessage(fmsg *goflowpb.FlowMessage) map[string]interface{} {
	msg := make(map[string]interface{})

	srcmac := make([]byte, 8)
	dstmac := make([]byte, 8)
	binary.BigEndian.PutUint64(srcmac, fmsg.SrcMac)
	binary.BigEndian.PutUint64(dstmac, fmsg.DstMac)
	srcmac = srcmac[2:8]
	dstmac = dstmac[2:8]

	msg["type"] = fmsg.Type.String()
	msg["time_received"] = int(fmsg.TimeReceived)
	msg["sequence_num"] = int(fmsg.SequenceNum)
	msg["sampling_rate"] = int(fmsg.SamplingRate)

	switch fmsg.Type {
	case goflowpb.FlowMessage_NETFLOW_V9, goflowpb.FlowMessage_IPFIX:
		msg["flow_direction"] = int(fmsg.FlowDirection)
	}

	msg["sampler_address"] = net.IP(fmsg.SamplerAddress).String()
	msg["time_flow_start"] = int(fmsg.TimeFlowStart)
	msg["time_flow_end"] = int(fmsg.TimeFlowEnd)
	msg["bytes"] = int(fmsg.Bytes)
	msg["packets"] = int(fmsg.Packets)
	msg["src_addr"] = net.IP(fmsg.SrcAddr).String()
	msg["dst_addr"] = net.IP(fmsg.DstAddr).String()
	msg["ethernet_type"] = int(fmsg.Etype)
	msg["proto"] = int(fmsg.Proto)
	msg["src_port"] = int(fmsg.SrcPort)
	msg["dst_port"] = int(fmsg.DstPort)
	msg["in_interface"] = int(fmsg.InIf)
	msg["out_interface"] = int(fmsg.OutIf)

	switch fmsg.Type {
	case goflowpb.FlowMessage_SFLOW_5, goflowpb.FlowMessage_NETFLOW_V9, goflowpb.FlowMessage_IPFIX:
		if fmsg.SrcMac != 0 {
			msg["src_mac"] = net.HardwareAddr(srcmac).String()
		}

		if fmsg.DstMac != 0 {
			msg["dst_mac"] = net.HardwareAddr(dstmac).String()
		}

		msg["src_vlan"] = int(fmsg.SrcVlan)
		msg["dst_vlan"] = int(fmsg.DstVlan)
		msg["vlan_id"] = int(fmsg.VlanId)
	}

	switch fmsg.Type {
	case goflowpb.FlowMessage_IPFIX:
		msg["ingress_vrf_id"] = int(fmsg.IngressVrfID)
		msg["egress_vrf_id"] = int(fmsg.EgressVrfID)
	}

	msg["ip_tos"] = int(fmsg.IPTos)

	switch fmsg.Type {
	case goflowpb.FlowMessage_NETFLOW_V9, goflowpb.FlowMessage_IPFIX:
		msg["forwarding_status"] = int(fmsg.ForwardingStatus)
	}

	switch fmsg.Type {
	case goflowpb.FlowMessage_SFLOW_5, goflowpb.FlowMessage_NETFLOW_V9, goflowpb.FlowMessage_IPFIX:
		msg["ip_ttl"] = int(fmsg.IPTTL)
	}

	msg["tcp_flags"] = int(fmsg.TCPFlags)

	switch fmsg.Type {
	case goflowpb.FlowMessage_SFLOW_5, goflowpb.FlowMessage_NETFLOW_V9, goflowpb.FlowMessage_IPFIX:
		msg["icmp_types"] = int(fmsg.IcmpType)
		msg["icmp_code"] = int(fmsg.IcmpCode)
		msg["ipv6_flow_label"] = int(fmsg.IPv6FlowLabel)
		msg["fragment_id"] = int(fmsg.FragmentId)
		msg["fragment_offset"] = int(fmsg.FragmentOffset)
	}

	switch fmsg.Type {
	case goflowpb.FlowMessage_IPFIX:
		msg["bi_flow_direction"] = int(fmsg.BiFlowDirection)
	}

	if int(fmsg.SrcAS) != 0 {
		msg["src_as"] = int(fmsg.SrcAS)
	}

	if int(fmsg.DstAS) != 0 {
		msg["dst_as"] = int(fmsg.DstAS)
	}

	if len(fmsg.NextHop) != 0 {
		msg["next_hop"] = net.IP(fmsg.NextHop).String()
	}

	switch fmsg.Type {
	case goflowpb.FlowMessage_SFLOW_5:
		msg["next_hop_as"] = int(fmsg.NextHopAS)
	}

	msg["src_net"] = int(fmsg.SrcNet)
	msg["dst_net"] = int(fmsg.DstNet)

	if fmsg.HasEncap {
		msg["has_encap"] = fmsg.HasEncap
		msg["src_addr_encap"] = net.IP(fmsg.SrcAddrEncap).String()
		msg["dst_addr_encap"] = net.IP(fmsg.DstAddrEncap).String()
		msg["proto_encap"] = int(fmsg.ProtoEncap)
		msg["ETypeEncap"] = int(fmsg.EtypeEncap)
		msg["IPTosEncap"] = int(fmsg.IPTosEncap)
		msg["IPTTLEncap"] = int(fmsg.IPTTLEncap)
		msg["IPv6FlowLabelEncap"] = int(fmsg.IPv6FlowLabelEncap)
		msg["FragmentIdEncap"] = int(fmsg.FragmentIdEncap)
		msg["FragmentOffsetEncap"] = int(fmsg.FragmentOffsetEncap)
	}

	if fmsg.HasMPLS {
		msg["HasMPLS"] = fmsg.HasMPLS
		msg["MPLSCount"] = int(fmsg.MPLSCount)
		msg["MPLS1TTL"] = int(fmsg.MPLS1TTL)
		msg["MPLS1Label"] = int(fmsg.MPLS1Label)
		msg["MPLS2TTL"] = int(fmsg.MPLS2TTL)
		msg["MPLS2Label"] = int(fmsg.MPLS2Label)
		msg["MPLS3TTL"] = int(fmsg.MPLS3TTL)
		msg["MPLS3Label"] = int(fmsg.MPLS3Label)
		msg["MPLSLastTTL"] = int(fmsg.MPLSLastTTL)
		msg["MPLSLastLabel"] = int(fmsg.MPLSLastLabel)
	}

	if fmsg.HasPPP {
		msg["HasPPP"] = fmsg.HasPPP
		msg["PPPAddressControl"] = int(fmsg.PPPAddressControl)
	}

	return msg
}
