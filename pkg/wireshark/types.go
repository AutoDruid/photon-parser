package wireshark

import "time"

type Source struct{
	Layers Layers `json:"layers"`
}

type Data struct{
	DataData string `json:"data.data"`
	DataLen  string `json:"data.len"`
}

type Layers struct{
			Frame struct {
				FrameSectionNumber   string `json:"frame.section_number"`
				FrameInterfaceID     string `json:"frame.interface_id"`
				FrameInterfaceIDTree struct {
					FrameInterfaceName        string `json:"frame.interface_name"`
					FrameInterfaceDescription string `json:"frame.interface_description"`
				} `json:"frame.interface_id_tree"`
				FrameEncapType          string    `json:"frame.encap_type"`
				FrameTime               time.Time `json:"frame.time"`
				FrameTimeUtc            time.Time `json:"frame.time_utc"`
				FrameTimeEpoch          time.Time `json:"frame.time_epoch"`
				FrameOffsetShift        string    `json:"frame.offset_shift"`
				FrameTimeDelta          string    `json:"frame.time_delta"`
				FrameTimeDeltaDisplayed string    `json:"frame.time_delta_displayed"`
				FrameTimeRelative       string    `json:"frame.time_relative"`
				FrameNumber             string    `json:"frame.number"`
				FrameLen                string    `json:"frame.len"`
				FrameCapLen             string    `json:"frame.cap_len"`
				FrameMarked             string    `json:"frame.marked"`
				FrameIgnored            string    `json:"frame.ignored"`
				FrameProtocols          string    `json:"frame.protocols"`
				FrameEncoding           string    `json:"frame.encoding"`
				FrameColoringRuleName   string    `json:"frame.coloring_rule.name"`
				FrameColoringRuleString string    `json:"frame.coloring_rule.string"`
			} `json:"frame"`
			Eth struct {
				EthDst     string `json:"eth.dst"`
				EthDstTree struct {
					EthDstResolved     string `json:"eth.dst_resolved"`
					EthDstOui          string `json:"eth.dst.oui"`
					EthDstOuiResolved  string `json:"eth.dst.oui_resolved"`
					EthDstLg           string `json:"eth.dst.lg"`
					EthDstIg           string `json:"eth.dst.ig"`
					EthAddr            string `json:"eth.addr"`
					EthAddrResolved    string `json:"eth.addr_resolved"`
					EthAddrOui         string `json:"eth.addr.oui"`
					EthAddrOuiResolved string `json:"eth.addr.oui_resolved"`
					EthLg              string `json:"eth.lg"`
					EthIg              string `json:"eth.ig"`
				} `json:"eth.dst_tree"`
				EthSrc     string `json:"eth.src"`
				EthSrcTree struct {
					EthSrcResolved  string `json:"eth.src_resolved"`
					EthSrcOui       string `json:"eth.src.oui"`
					EthSrcLg        string `json:"eth.src.lg"`
					EthSrcIg        string `json:"eth.src.ig"`
					EthAddr         string `json:"eth.addr"`
					EthAddrResolved string `json:"eth.addr_resolved"`
					EthAddrOui      string `json:"eth.addr.oui"`
					EthLg           string `json:"eth.lg"`
					EthIg           string `json:"eth.ig"`
				} `json:"eth.src_tree"`
				EthType   string `json:"eth.type"`
				EthStream string `json:"eth.stream"`
			} `json:"eth"`
			IP struct {
				IPVersion     string `json:"ip.version"`
				IPHdrLen      string `json:"ip.hdr_len"`
				IPDsfield     string `json:"ip.dsfield"`
				IPDsfieldTree struct {
					IPDsfieldDscp string `json:"ip.dsfield.dscp"`
					IPDsfieldEcn  string `json:"ip.dsfield.ecn"`
				} `json:"ip.dsfield_tree"`
				IPLen       string `json:"ip.len"`
				IPID        string `json:"ip.id"`
				IPFlags     string `json:"ip.flags"`
				IPFlagsTree struct {
					IPFlagsRb string `json:"ip.flags.rb"`
					IPFlagsDf string `json:"ip.flags.df"`
					IPFlagsMf string `json:"ip.flags.mf"`
				} `json:"ip.flags_tree"`
				IPFragOffset     string `json:"ip.frag_offset"`
				IPTTL            string `json:"ip.ttl"`
				IPProto          string `json:"ip.proto"`
				IPChecksum       string `json:"ip.checksum"`
				IPChecksumStatus string `json:"ip.checksum.status"`
				IPSrc            string `json:"ip.src"`
				IPAddr           string `json:"ip.addr"`
				IPSrcHost        string `json:"ip.src_host"`
				IPHost           string `json:"ip.host"`
				IPDst            string `json:"ip.dst"`
				IPDstHost        string `json:"ip.dst_host"`
				IPStream         string `json:"ip.stream"`
			} `json:"ip"`
			UDP struct {
				UDPSrcport        string `json:"udp.srcport"`
				UDPDstport        string `json:"udp.dstport"`
				UDPPort           string `json:"udp.port"`
				UDPLength         string `json:"udp.length"`
				UDPChecksum       string `json:"udp.checksum"`
				UDPChecksumStatus string `json:"udp.checksum.status"`
				UDPStream         string `json:"udp.stream"`
				UDPStreamPnum     string `json:"udp.stream.pnum"`
				Timestamps        struct {
					UDPTimeRelative string `json:"udp.time_relative"`
					UDPTimeDelta    string `json:"udp.time_delta"`
				} `json:"Timestamps"`
				UDPPayload string `json:"udp.payload"`
			} `json:"udp"`
			Data Data `json:"data"`

}

type WireSharkRequest struct {
	Index  string `json:"_index"`
	Score  any    `json:"_score"`
	Source Source `json:"_source"`
}
