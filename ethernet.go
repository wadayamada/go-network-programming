// package temp

// import (
// 	"fmt"
// 	"net"
// 	"syscall"
// )

// // 必要な作業
// // - 通信に使うネットワークインターフェースを取得して、addressを取得する。(アドレス)
// // - raw socketを作る(ファイルディスクリプタ)
// // - ethernet frameを作る(データ)
// // - sendTo(ファイルディスクリプタ、データ、フラグ、アドレス)
// func main() {

// 	// ネットワークインターフェース
// 	iface, err := net.InterfaceByName("lo0")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	// SockaddrDatalink構造体を作成します。
// 	addr := syscall.SockaddrDatalink{
// 		Family: syscall.AF_LINK,
// 		Index:  uint16(iface.Index),
// 		Type:   syscall.AF_LINK, // イーサネットタイプ
// 		Nlen:   uint8(len(iface.HardwareAddr)),
// 		Alen:   uint8(6), // MACアドレスの長さ（6バイト）
// 		Slen:   uint8(0), // ソケットレベルでの長さ（通常は0）
// 		Data:   [12]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
// 	}

// 	// raw socketを作る
// 	// macOS では Linux の AF_PACKET に相当する仕組みがない
// 	// https://zenn.dev/satoken/articles/golang-tcpip
// 	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(fd)
// 	// ethernet frameを作る
// 	// 宛先を指定する
// 	// addr := syscall.SockaddrDatalink{}
// 	// 送信する
// 	err = syscall.Sendto(fd, []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}, 0, &addr)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	// // io.Writerで書き込む
// 	// // file := os.NewFile(uintptr(fd), "")
// 	// // writer := bufio.NewWriter(file)
// 	// // writer.Write([]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55})

// 	// fmt.Println("Hello, world!")
// }
