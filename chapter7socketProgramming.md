## TCP

- クライアントの流れ
  - socketシステムコール
  - connectシステムコール
  - writeシステムコール
  - closeシステムコール
- サーバーの流れ
  - socketシステムコール
  - bindシステムコール
  - listenシステムコール
    - TCP SYNを待つ
  - acceptシステムコール
    - listenしてるソケットからやり取りするための新たなソケットを作る
  - closeシステムコール
- bind
  - サーバー側でlistenの前にbindすると、クライアント側はそのportに繋ぐ必要がある
  - クライアントはconnectの前にbindするとそのportを使ってサーバー側にconnectする
  - 特にbindしなければ、カーネルが自動でbindしてくれる
  - サーバーもbindせずにlistenできる。その場合はカーネルが自動で割り当ててくれる。
    - getsocknameでポート番号を取得できる

## UDP
送信側がいきなり送信する
- 送信側
  - socketシステムコール
  - sendtoシステムコール
    - connectせずにいきなり送れる
  - close
- 受信側
  - socket
  - bind
  - recvfrom
  - close

- 受け取って返信もできる
  - recvfromで、送信元ポートも渡されるので、それを使って返信ができる
  - DNSもこの方式で名前解決の結果を返してる
  - 偽造に気付けないので注意
    - TCPはシーケンス番号が同期してないとパケットを受け付けないが、UDPにはそう言った仕組みがない
      - 3way handshakeでサーバー側とクライアント側のシーケンス番号を同期させて、IPパケットごとにackで確認してるからシーケンス番号がわからないと、偽造はできない
    - UDPでは開発者が偽造を検知する仕組みを実装する必要がある
    - UDPセグメントのIPヘッダーに記載されている送信元IPアドレスを自由に書き換えられる。面白い。
      - DNSは素のUDPを使ってるのではなく、トランザクションIDとかつけて偽造などの対策をしてる
    - HTTPSでTLSで暗号化できるのはTCPセグメントのペイロード部分だけなので、IPヘッダー、TCPヘッダーにある、5tapleは見えてる

- IPパケットを直接送信できるソケットもある: raw socket
  - raw socketでイーサネットフレームも
  - tapではイーサネットフレームを、tunではIPパケットを送れるが、これらは仮想的なネットワークデバイスであり、これらはカーネルのネットワークスタックを通るが、raw socketは通らないらしい
- ローカルでの通信に使えるソケットもある: unix domain socket

## 名前解決
- getaddrinfoで名前解決できる
  - そのままsocket作れる
  - getaddrinfoはfopen, fclose, printfと同様にカーネルに実装されてるのではなく、libcらしい
    - 内部ではsocketシステムコールやら、fopenならopenシステムコール読んでるが、それで済んじゃうって意味
    - カーネルしかアクセスできないリソースは使わずに、システムコールだけで実現できる。
- 2011年にIPv4は在庫が枯渇した
- DNSではIPv4, IPv6の両方を問い合わせることもできる
- iOSでは両方に対応してないアプリケーションは審査に通らない
- gai.confでgetaddrinfoが返すIPアドレスの優先度を設定できる

## VPN
- VPNクライアントを用いて、VPNサーバーと通信する
- IPパケットをIPヘッダーも含めて暗号化する
  - 宛先IPアドレスをわからないようにする。実際にネットワークを流れるIPパケットの宛先IPアドレスはVPNサーバーになる。
- VPNサーバーを通じて、社内ネットワークにアクセスできる。プライベートIPアドレスが割り当てられる
