# go-traphandle

Net-Snmp の `snmptrapd` の `traphandle` で使用するためのプログラムです。

## Install

```sh
go get github.com/ngyuki/go-traphandle/cmd/go-traphandle
```

## Example

`go-traphandle` の設定ファイルを下記のように作成します。
ファイル名はなんでもよいです（仮に `go-traphandle.yml`）。
メールアドレス、IPアドレス、コミュニティ名などは適宜変更してください。

```yaml
actions:
  oreore:
    emails:
      - host: localhost
        port : 25
        from: root@example.com
        to:  ore@example.com

defaults:
  url: http://example.com/

matches:
  - trap: IF-MIB::linkDown
    ipaddr: 192.168.0.0/16
    community: public
    bindings:
      index: RFC1213-MIB::ifIndex.*
      status: RFC1213-MIB::ifOperStatus.*
    conditions:
      index: { eq: [1, 2, 3, 4, 5] }
      status: { regexp: [down] }
    formats:
      subject: |
        {{.ipaddr}} Interface {{.index}} linkdown
      body: |
        Date: {{.date}}
        Ipaddr: {{.ipaddr}}
        Message: Interface {{.index}} linkdown
        ---
        {{.url}}
    action: oreore
    fallthrough: True
```

`go-traphandle` を次のように実行します。

```sh
./go-traphandle -server 127.0.0.1:9999 -config go-traphandle.yml
```

`snmptrapd.conf` を次のように設定します。

```
[snmp]
noDisplayHint yes

[snmptrapd]
disableAuthorization yes
ignoreAuthFailure yes
outputOption n
traphandle default /usr/bin/nc 127.0.0.1 9999
```

snmptrapd をリスタートします。

```sh
systemctl restart snmptrapd
```

`IF-MIB::linkDown` のインタフェース番号 1-5 のいずれかのリンクダウンのトラップを受信するとメールで通知されます。

## コマンドライン

```
go-traphandle -server <bind-address> -config <config-file>
```

`go-traphandle` はサーバとして動作します。`snmptrapd` の `traphandle` で `go-traphandle` がリッスンしているポートへトラップの内容を送信する必要があります。↑の例では `nc` コマンドでトラップを送信しています。

`-server` でリッスンするアドレス・ポートを指定します。

`-config` で設定ファイルを指定します。設定ファイルでは受けたトラップをどのように処理するか記述します。

## 設定ファイル

受信したトラップは設定ファイルの `matches` の上から順番にチェックされ、最初にマッチしたエントリの処理が実行されます。

### matches

**trap**

トラップの OID をチェックします（完全一致）。省略された場合はトラップの OID はチェックされません。

**ipaddr**

ソースアドレスをチェックします（CIDR形式または完全一致）。
省略された場合はソースアドレスはチェックされません。

**community**

トラップのコミュニティ名をチェックします（完全一致）。
省略された場合はトラップのコミュニティ名はチェックされません。

**bindings**

トラップのデータバインディングに名前を付けて `conditions` や `formats` で変数として使用できるようにします。

例えば `val: "XXX-MIB::sample"` と設定した場合、トラップのデータバインディングに `XXX-MIB::sample` という OID が含まれていればその値を `conditions` や `formats` で `val` という変数として使用できます。データバインディングに指定した OID が含まれていない場合、変数には空文字が入ります。

OID の最後のセグメントにはワイルドカードが使用可能です。
`*` は任意の数の数字に一致し、`**` は任意の数の `.` および数字に一致します。例えば下記のようになります。

```
XXX-MIB::sample.*   XXX-MIB::sample.1   ... OK
XXX-MIB::sample.*   XXX-MIB::sample.1.2 ... miss
XXX-MIB::sample.**  XXX-MIB::sample.1   ... OK
XXX-MIB::sample.**  XXX-MIB::sample.1.2 ... OK
```

**conditions**

`bindings` で名前を付けたデータバインディングの値に対して様々な条件でマッチします。

次のような形式で指定します。

```yaml
    conditions:
      index:
        eq: [1, 2, 3, 4, 5]
      status:
        regexp: [down]
```

`conditions` の直下のキーは変数名です。この変数名は `bindings` で名前付けられている必要があります。複数の変数が指定された場合、それらすべての変数の値が条件にマッチする必要があります。

その下のキーは条件名です。次のものが使用可能です。

- `eq`
- `not_eq`
- `regexp`
- `not_regexp`

複数の条件が指定された場合、それらすべての条件に変数の値がマッチする必要があります。

条件の値は配列で指定します。これは OR 条件なのでいずれかの値にマッチすればよいです。

**formats**

メール送信やスクリプト実行のためにメッセージを書式化します。

```yaml
    formats:
      subject: |
        {{.ipaddr}} Interface {{.index}} linkdown
      body: |
        Date: {{.date}}
        Ipaddr: {{.ipaddr}}
        Message: Interface {{.index}} linkdown
        ---
        {{.url}}
```

`formats` の直下の `subject` と `body` がメールの件名と本文になります。

テンプレートエンジンには [template/text](https://golang.org/pkg/text/template/) を使用しています。

テンプレートでは、後述の `default` で記述された変数、`bindings` で名前付けられた変数、`date` および `ipaddr` が使用できます。`date` は現在日時に、`ipaddr` はトラップのソースアドレスに展開されます。

なお、`formats` の直下のキーには `subject` と `body` 以外も使用可能です。メール送信では `subject` と `body` しか使用されないためその他のキーには意味がありませんが、スクリプトの実行では `subject` と `body` 以外も環境変数で参照できます。

**action**

条件にマッチした場合に実行する、`actions` のキー名を指定します。後述の `actions` の説明を参照してください。

**fallthrough**

条件にマッチした場合に、以降のマッチを続けるかどうかを指定子ます。

`False` だと条件にマッチした時点で当該トラップの処理は終了します（デフォルト）。

`True` だと条件にマッチした後も以降のマッチがチェックされます。

### actions

`matches` の条件にマッチしたときに実行されるアクションを定義します。

```yaml
actions:
  oreore:
    emails:
      - host: localhost
        port : 25
        from: root@example.com
        to:  ore@example.com
      - host: localhost
        port : 25
        from: root@example.com
        to:  are@example.com
    scripts:
      - env | sort | logger -t snmptrap
```

`actions` の直下のキーが `matches` の `action` で指定する値です。その下のキーが実行するアクションの種類を表しており、下記のいずれかが指定できます。

- `emails`
- `scripts`

`emails` では `host,port,from,to` の４項目の連想配列の配列を指定します。`host` と `port` は省略可能で、省略された場合はそれぞれ `localhost` と `25` になります。

`scripts` には実行するコマンドを複数指定します。コマンドは `/bin/sh` 経由で実行されます。スクリプト実行時に、`bindings` で名前付けられたデータバインディングの変数や `formats` で書式化された変数が、プレフィックス `TH_` を付与された上で設定されます。例えば `index` という変数は `TH_index` という環境変数になります。

### default

`formats` で使用可能な変数のデフォルト値です。

## 変数の詳細

`matches` の `trap`, `ipaddr`, `community` のチェックにパスすると下記の変数が定義されます。

- `ipaddr`
- `date`
- `bindings` で指定したデータバインディングの値

さらに `conditions` の条件にパスすると `defaults` で記述された値が追加で定義されます。`defaults` で記述された変数名が既に定義されている場合（`bindings` に同じキー名が記述されているなど）、上書きはされません。

次に `formats` が適用され、その結果が追加で定義されます。このとき同じ変数名が既に定義されていれば上書きされます。

最後に、メール送信で `subject` と `body` という変数が件名と本文に、スクリプト実行ですべての変数がプレフィックス `TH_` を付けて環境変数に設定されます。

## snmptrapd の設定

`snmptrapd` の `traphandle` に渡されるトラップの内容は OID が生の数値になっている前提です。なので `snmptrapd.conf` で `outputOption n` を指定して生の値が使われる必要があります。

また、マルチバイト文字を含むトラップを受信したときに情報が欠損しないようにするために `[snmp]` セクションで `noDisplayHint yes` もオススメです。
