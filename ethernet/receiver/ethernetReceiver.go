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
		// fmt.Println("Received a packet")
		// fmt.Println(packet)
		// Ethernet フレームの解析
		ethLayer := packet.Layer(layers.LayerTypeEthernet)
		if ethLayer != nil {
			eth, _ := ethLayer.(*layers.Ethernet)

			// カスタム EtherType のパケットを抽出
			if eth.EthernetType == layers.EthernetType(0x88B5) {
				fmt.Printf("Received custom Ethernet packet: %x\n", packet.Data())
			}
		}
	}
}
