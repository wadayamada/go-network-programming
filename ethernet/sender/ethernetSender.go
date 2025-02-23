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

	// Ethernet ヘッダの作成
	// MACアドレスは48bit
	eth := &layers.Ethernet{
		SrcMAC:       []byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}, // 送信元MAC
		DstMAC:       []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, // ブロードキャスト
		EthernetType: layers.EthernetType(0x88B5),                // カスタム EtherType
	}

	// カスタムデータを作成
	customPayload := []byte("This is custom raw data")

	// パケットをシリアライズ
	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{ComputeChecksums: false, FixLengths: false}
	gopacket.SerializeLayers(buffer, options, eth, gopacket.Payload(customPayload))

	// パケットを送信
	err = handle.WritePacketData(buffer.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Custom Ethernet packet sent successfully")
}
