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

	// 宛先 IP アドレス
	dstIP := net.ParseIP("23.218.94.19") // 送信先の IP（適宜変更）

	// 送信元 MAC & 宛先 MAC（適宜変更）
	srcMAC, _ := net.ParseMAC("00:11:22:33:44:55") // 自分の MAC アドレス
	dstMAC, _ := net.ParseMAC("66:77:88:99:AA:BB") // 送信先のルーター等の MAC

	// Ethernet ヘッダー
	eth := &layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeIPv4,
	}

	// IP ヘッダー
	ip := &layers.IPv4{
		Version:  4,
		IHL:      5,
		TTL:      64,
		Protocol: layers.IPProtocolICMPv4,
		SrcIP:    net.ParseIP("192.0.0.2"), // 自分の IP
		DstIP:    dstIP,
	}

	// ICMP ヘッダー（Echo Request）
	icmp := &layers.ICMPv4{
		TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
		Id:       12345,
		Seq:      1,
	}

	// シリアライズ
	buffer := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	err = gopacket.SerializeLayers(buffer, opts, eth, ip, icmp, gopacket.Payload([]byte("Hello!")))
	if err != nil {
		log.Fatal(err)
	}

	// パケット送信
	err = handle.WritePacketData(buffer.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ICMP Echo Request sent!")

	// 念のため少し待つ
	time.Sleep(1 * time.Second)
}
