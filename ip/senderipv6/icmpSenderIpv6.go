package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	// 使用するネットワークインターフェースを指定
	ifaceName := "en0" // macOS の Wi-Fi の場合
	handle, err := pcap.OpenLive(ifaceName, 1600, false, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// 宛先 IPv6 アドレス
	dstIP := net.ParseIP("64:ff9b::17d9:88d0") // 送信先の IPv6 アドレス（適宜変更）

	// 送信元 MAC & 宛先 MAC（適宜変更）
	srcMAC, _ := net.ParseMAC("10:B5:88:4F:0C:2D") // 自分の MAC アドレス
	dstMAC, _ := net.ParseMAC("a2:78:2d:2c:4d:64") // 送信先のルーター等の MAC（適宜変更）

	// Ethernet ヘッダー
	eth := &layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeIPv6,
	}

	// IPv6 ヘッダー
	ipv6 := &layers.IPv6{
		Version:    6,
		NextHeader: layers.IPProtocolICMPv6,
		SrcIP:      net.ParseIP("240a:61:7282:d8:8a0:3529:a656:13c0"), // 送信元 IPv6 アドレス（適宜変更）
		DstIP:      dstIP,
		HopLimit:   64,
	}

	// ICMPv6 ヘッダー（Echo Request）
	icmpv6 := &layers.ICMPv6{
		TypeCode: layers.CreateICMPv6TypeCode(layers.ICMPv6TypeEchoRequest, 0),
		// Id:       12345,
		// Seq:      1,

	}
	// v6のチェックサムを計算するために、IPv6のヘッダをセットする
	// これに気づくのに時間かかった
	icmpv6.SetNetworkLayerForChecksum(ipv6)

	// シリアライズ
	buffer := gopacket.NewSerializeBuffer()

	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	err = gopacket.SerializeLayers(buffer, opts, eth, ipv6, icmpv6, gopacket.Payload([]byte("Hello!")))
	if err != nil {
		log.Fatal(err)
	}

	// パケット送信
	err = handle.WritePacketData(buffer.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ICMPv6 Echo Request sent!")

	// 念のため少し待つ
	time.Sleep(1 * time.Second)
}
