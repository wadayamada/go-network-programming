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
	dstIP := net.ParseIP("2406:da14:d95:5100:b380:ccde:1ef2:3b2c") // 送信先の IPv6 アドレス（適宜変更）
	dstPort := uint16(1)

	// 送信元 MAC & 宛先 MAC（適宜変更）
	srcMAC, _ := net.ParseMAC("10:B5:88:4F:0C:2D") // 自分の MAC アドレス
	dstMAC, _ := net.ParseMAC("a2:78:2d:2c:4d:64") // 送信先のルーター等の MAC（適宜変更）

	srcPort := uint16(60877)

	// Ethernet ヘッダー
	eth := &layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeIPv6,
	}

	// IPv6 ヘッダー
	ipv6 := &layers.IPv6{
		Version:  6,
		SrcIP:    net.ParseIP("240a:61:7282:d8:e4af:9b1:1ed9:eba0"), // 送信元 IPv6 アドレス（適宜変更）
		DstIP:    dstIP,
		HopLimit: 64,
		Length:   8 + uint16(len("hi")), // ヘッダーサイズ（8 バイト） + ペイロードの長さ
		// プロトコルのこと
		NextHeader: layers.IPProtocolUDP,
	}

	// UDP ヘッダー
	udp := &layers.UDP{
		SrcPort: layers.UDPPort(srcPort),
		DstPort: layers.UDPPort(dstPort),
		Length:  8 + uint16(len("hi")), // ヘッダーサイズ（8 バイト） + ペイロードの長さ
	}

	// v6のチェックサムを計算するために、IPv6のヘッダをセットする
	// これに気づくのに時間かかった
	udp.SetNetworkLayerForChecksum(ipv6)

	// シリアライズ
	buffer := gopacket.NewSerializeBuffer()

	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	err = gopacket.SerializeLayers(buffer, opts, eth, ipv6, udp, gopacket.Payload([]byte("hi")))
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
