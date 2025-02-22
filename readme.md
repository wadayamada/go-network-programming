# 概要
- [ポートとソケットがわかればインターネットがわかる](https://gihyo.jp/book/2016/978-4-7741-8570-5)を読んで、メモを書く
- 気になってるネットワーク周辺知識も調べる

# webブラウザを使って通信する流れ
- webブラウザで"http://www.example.com"と打つ

## 名前解決
- www.example.comの名前解決をDNSで行い、IPアドレスが手に入る
- IPv4ならAレコードで、IPv6はAAAAレコード
- 問い合わせはキャッシュDNSサーバーに行う
- キャッシュDNSサーバーが知らなかった場合
  - キャッシュDNSサーバーがcomの名前解決をルートDNSサーバーに問い合わせる。ルートDNSサーバーはcomの権威DNSサーバーのIPアドレスを教えてくれる
    - ルートDNSサーバーはaからmまで13系統ある
    - 系統ごとにも同じIPアドレスで複数のサーバーを立てる、IPエニーキャストしてる
    - 13系統全てで同じ情報を持っている。ルートゾーン
    - ルートゾーンはTLDの名前解決情報のこと
  - comの権威DNSサーバーにexample.comの名前解決を尋ねて、example.comの権威DNSサーバーのIPアドレスを教えてもらう
  - example.comの権威DNSサーバーにwww.example.comのIPアドレスを教えてもらう

## TCPコネクションを作る
- socketシステムコールでTCPソケットを作る。
  - プロセスがインターネット通信するためのファイル
  - これにread, writeするとインターネット通信ができる
- bindシステムコールでport=80とIPアドレス=23.208.31.49をソケットファイルにbindする
- connectシステムコールでTCPコネクションを張る
  - 3way handshakeでコネクションが張れる
    - クライアントがTCP SYNパケットを送る(seq=1)
    - サーバーがTCP SYN ACKパケットを送る(seq=1000, ack=2)
    - クライアントがTCP ACKパケットを送る(seq=2, seq=1001)
  - パケットを送る仕組みは後段のwriteシステムコールで詳しく書く
  - TCPではデータが喪失しないため、帯域を適度に使うために、順序制御、再送制御、輻輳制御(ウィンドウ制御とフロー制御)を行う
    - TCPではパケットを送るごとにTCP ACKパケットを返すことで、受信確認する
    - それぞれのパケットにseqがついてるので、順序を制御する。並び替えたり。
    - 届かなかったseqは再送してもらう
    - 帯域を使いすぎないための輻輳制御もする
      - ackなしで送信できるTCPセグメントの数を倍々にしていく。ウィンドウサイズ
      - コネクション確立直後はスロースタートだから帯域が狭い
      - パケット喪失があったら、ウィンドウサイズを1に戻す
- サーバーがacceptシステムコールでクライアントとの通信用のTCPソケットを作る
  - 今まではListenソケットでやり取りしてた(はず)

ちなみに
- 送信元IPアドレス、送信元ポート番号、受信側IPアドレス、受信側ポート番号、プロトコルの5つを5tupleという
  - これの組み合わせごとにセッションが作れる
- UDP
    - IPにはユニキャスト、マルチキャスト、ブロードキャストの仕組みがある
    - TCPはユニキャストしか使えないが、UDPは全て使える
    - 諸々の制御をしないので、パケットロスすることはあるが、再送による遅れや、受信確認や、輻輳制御がないので、速度は速い。リアルタイム性は高い。
    - そのまま使うとなりすましのリスクがある

## TCPコネクションでバイトストリームを送る

- writeシステムコールでTCPでバイトストリームを送ることができる
- HTTPリクエストである以下もTCPでバイトストリームとして送信する
  ```
  GET / HTTP/1.1
  Host: www.example.com
  ```
- TCPのバイトストリームはACKで受信確認しながら、IPパケットに分けて送信される
  - 1つのIPパケットは最大1600バイトくらいらしい。
  - 1文字2バイトだとして、800文字くらい送れる

## IPパケットの送信
- IPパケットはパソコンから家庭内LANのデフォルトゲートウェイであるNATルーターに届く
- NATルーターからインターネットを経由して、目的のIPアドレスのあるサーバーに届く

## NATルーターに届く
IPパケットはLANからWANへの入り口となるデフォルトゲートウェイに届けられる
このLANでのIPパケット送信はL2レイヤーで行われる
wifiでパソコンがルーターに繋がってるとする。wifiのルーターがデフォルトゲートウェイとなる。
`netstat -rn`でデフォルトゲートウェイのIPアドレスがわかる

OSがデフォルトゲートウェイのIPアドレスを指定して、NICドライバにIPパケットを渡す
NICドライバはARPで、IPアドレスからmacアドレスを取得する。そのmacアドレスに送信するためのwifiのネットワークインターフェースを特定する
NICドライバがIPパケットに宛先macアドレスと送信元macアドレスを追加して、wifiフレームを作る
送信元macアドレスはこのパソコンのwifiに繋がってるネットワークインターフェースのmacアドレス

NICドライバがネットワークインターフェースに対してwifiフレームを送ると、wifiを経由してデフォルトゲートウェイにwifiフレームが届く

`route -n get default`でデフォルトゲートウェイとの通信に使ってるネットワークインターフェースがわかる
`ifconfig`でネットワークインターフェースの一覧がわかる

関連
- docker0はdockerホストに作成される仮想的なネットワークインタフェース
- dockerコンテナの中からはdocker0は見えず、eth0が見える
- dockerコンテナ内ではdocker0のIPアドレスがデフォルトゲートウェイとして見えてる

スイッチは複数のコンピューターを接続し、L2レベルでイーサネットフレームの受け渡しを行う
wifiで言うと、wifiの送受信を行っているアクセスポイントがこのスイッチと同じ役割をしてる

NATルーターに届いたIPパケットは、送信元IPアドレスをプライベートなものからパブリックなものに変換される。ポート番号も変換する
- NATテーブルに書く
- TCPヘッダー, IPヘッダーのチェックサムも更新する
- サーバー側からは同じIPアドレスから通信が来たように見える

## ルーティング
NATルーターから先はインターネットに接続し、L3レベルのルーティングによって、目的のIPアドレスにIPパケットを届ける

ルーターがIPパケットを転送し続けることで、目的のIPアドレスに届く
`traceroute`で確認できる
IPパケットにはTTLというルーターの転送可能回数を設定できる。ルーターが転送するごとに、これを減らす。

- IPアドレスはASというネットワークに属する、AS間のルーティングはBGPで行う
- AS内のルーティングはASそれぞれが独自の方法で行う

- 個々のルーター同士はL2レベルでmacアドレスを指定して、イーサネットフレームを送りあう



## 物理回線
物理的な回線では、以下の流れで日本からアメリカに通信してる
- パソコン
- NATルーター
- 収容局
- 陸揚げ局
- 光海底ファイバの太平洋横断回線
- 陸揚げ局
- データセンター
- データセンター内の回線

また、NATルーターからは光回線であれば、ONUを経由してインターネットに接続してる
光回線とは光ファイバでデータを光信号として送り合う方式
ONUはデジタル信号と光信号の相互変換をしてる
ONUではL2レベルではEthernet上でPPPを行うPPPoEで通信してる
PPPはユーザー認証やセキュリティやなどが行えるプロトコルらしい

# Webサーバーに届く
- LB: 複数のアプリケーションサーバーに振り分ける
- CDN: 地理的に分散させたサーバーで主に静的なコンテンツを配布する

- NICはイーサネットフレームを受信すると、CPUにハードウェア割り込みをかける
- 割り込み頻度を調節する、割り込みコアリングとかもしてる
- CPUの周波数は10^9オーダーだが、NIC割り込みは10^6オーダー。ネットワーク負荷が無視できなくなることもあるみたい。

- アプリケーションレベルではreadシステムコールで、TCPソケットファイルからバイトストリームを取得する
- カーネルがIPパケットからTCPバイトストリームを作って、アプリケーションに渡してくれる
- カーネルはソケットバッファで一時待機させてから渡す。届いてないIPパケットがあった場合、ここで届いたやつを一時保存して待つなどをする。

## HTTPレスポンスを受け取る

- HTTPサーバーはwriteシステムコールで、HTTPレスポンスのテキストをTCPソケットファイルに書き込む
- HTTP/1.1でkeepaliveしてないなら、closeシステムコールでTCPソケットファイルを閉じる
- クライアントがTCPソケットファイルからreadして、HTTPレスポンスがTCPストリームにテキストとして返ってくる
- closeシステムコールでTCPソケットファイルを閉じる
- HTTPレスポンスをパースして、画面に表示する

# やりたい
## UDP
- UDPを使ってZoomのようにカメラ映像や音声を送受信できるサービス作ってみたい
  - 自作AbemaTVやりたい
- ブロードキャストやマルチキャスト
  - これ自体はIPレイヤーの仕組みで、UDPだと扱えるって感じ
- UDPパケットの偽装。TCPではできない。
- UDP通信のプログラミング
- UDP通信で返信
- TCP, UDPの成りすまし確認
  - VPNでどこまで防げるか？
## 色々
- インターフェースはファイルシステムで、内部ではGCP使ってるとか、HTTPリクエストしてるとかのやつやりたい
- HTTP/2, HTTP/3, QUICの実装
- メールサーバー、SMTP、POP、IMAP
- P2P通信、BitTorrent
- NAS: Network Attached Storage
- VPN
- HTTP2とHTTP3を雑に実装したい
- TCP通信のプログラミング
- getaddrinfoを使ってそのままTCP通信
- ICMPを試す
- ICMP Time Exceedも見てみたい
- ヘッダーを見て、wifi通信とethernet通信を見分けたい

## ネットワーク以外
- コンテキストスイッチは自作OSしないと具体的にどれくらいコストがあるのかイメージつかない
  - カーネルとユーザーアプリケーションの行き来がめっちゃ多そう
    - システムコール
    - スケジューリング
- フリップフロップをハードで実装したい
- Cで標準入出力のシステムコールを叩きたい

# 調べたい
- 分電盤、電気工事士、UPS、ワットとアンペア
- DHCP
- spotifyの工夫
  - https://scrapbox.io/musicsurvey/%E8%AB%96%E6%96%87_:_Spotify--large_scale,_low_latency,_P2P_music-on-demand_streaming
  - https://www.geekpage.jp/blog/?id=2017-6-20-2

# わからない
datagram型のunixドメインソケットの場合、クライアント用のUnixドメインソケットファイルが必要な理由がわからない
- サーバー側でnet.ListenPacketしてからクライアント側でnet.Dialして、conn.WriteToしたものはクライアントで受け取れたが、新しくサーバー側でnet.Dialでconnを作って、conn.Writeしたものはクライアントで受け取れなかった