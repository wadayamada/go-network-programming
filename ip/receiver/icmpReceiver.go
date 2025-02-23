package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	// ネットワークインターフェースを開く
	handle, err := pcap.OpenLive("en0", 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// パケットの受信ループ
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// IP パケットをチェック（IPv4 or IPv6）
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			if ip.Protocol == layers.IPProtocolICMPv4 {
				icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
				if icmpLayer != nil {
					icmp, _ := icmpLayer.(*layers.ICMPv4)
					fmt.Printf("Received ICMPv4 packet: Type=%d, Code=%d\n", icmp.TypeCode.Type(), icmp.TypeCode.Code())
				}
			}
		}
	}
}
