- ping
  - 宛先までIPパケットが届いてるかどうか確認できる
  - TTL: ルーターが転送できる回数も表示される
  - ICMPプロトコルでIPパケットを送ってる
    - IPプロトコル番号は1。TCPは6
    - ICMP EchoとICMP Echo Replyを使ってる
    - EchoはEchoを返してくれというもので、Echo ReplyはEchoに対する応答
    - カーネルが勝手に返信してくれるため、アプリケーションの実装は不要

- traceroute
  - 指定した宛先までの途中ルーターがわかる
  - ルーター間のRTTもわかる
  - TTLがゼロになった時、ルーターは送信元IPアドレスにICMP Time Exceedを送る
  - tracerouteはTTLを1つずつ大きくしていき、ICMP Time Exceedの結果を確認することでどのルーターを通ってるか確認してる
  - 目的IPアドレスに着いた時はICMP Time Exceedが返ってこない
    - ICMP Echoを送ることで、Echo Replyを見て確かめる
    - UDPで適当な空いてないポートに繋ごうとして、ICMP Destination Unreachableが返ってきて確かめる

- dig
  - 名前解決の確認
  - `dig www.example.com`: キャッシュDNSサーバーでIPv4アドレス
  - `dig aaaa www.example.com`: キャッシュDNSサーバーでIPv6アドレス
  - @で権威DNSサーバーに直接問い合わせられる
  - +traceでキャッシュDNSサーバー、権威DNSサーバーの反復検索も確認できる

- wireshark
  - パケットキャプチャ
  - tcpストリームも見られる

# やりたい
- pingを試す
- ICMPを試す
- ICMP Time Exceedも見てみたい