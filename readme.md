# 概要
- [ポートとソケットがわかればインターネットがわかる](https://gihyo.jp/book/2016/978-4-7741-8570-5)を読んで、メモを書く
- 気になってるネットワーク周辺知識も調べる

# やりたい
- UDPを使ってZoomのようにカメラ映像や音声を送受信できるサービス作ってみたい
  - 自作AbemaTVやりたい
- インターフェースはファイルシステムで、内部ではGCP使ってるとか、HTTPリクエストしてるとかのやつやりたい
- ブロードキャストやマルチキャスト
  - これ自体はIPレイヤーの仕組みで、UDPだと扱えるって感じ
- UDPパケットの偽装。TCPではできない。
- HTTP/2, HTTP/3, QUICの実装
- メールサーバー、SMTP、POP、IMAP
- P2P通信、BitTorrent
- NAS: Network Attached Storage
- VPN

## ネットワーク以外
- コンテキストスイッチは自作OSしないと具体的にどれくらいコストがあるのかイメージつかない
  - カーネルとユーザーアプリケーションの行き来がめっちゃ多そう
    - システムコール
    - スケジューリング
- フリップフロップをハードで実装したい

# 調べたい
- 分電盤、電気工事士、UPS、ワットとアンペア
- DHCP
- spotifyの工夫
  - https://scrapbox.io/musicsurvey/%E8%AB%96%E6%96%87_:_Spotify--large_scale,_low_latency,_P2P_music-on-demand_streaming
  - https://www.geekpage.jp/blog/?id=2017-6-20-2

# わからない
datagram型のunixドメインソケットの場合、クライアント用のUnixドメインソケットファイルが必要な理由がわからない
- サーバー側でnet.ListenPacketしてからクライアント側でnet.Dialして、conn.WriteToしたものはクライアントで受け取れたが、新しくサーバー側でnet.Dialでconnを作って、conn.Writeしたものはクライアントで受け取れなかった
