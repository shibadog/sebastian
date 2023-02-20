# sebastian

## oerview

`golang` のプロジェクトを開始するためには、以下のように、自身のモジュール名をつける必要がある様子。  
また、このときのモジュール名はリポジトリのパスを指定するものの様子。

``` bash
# go mod init github.com/shibadog/sebastian
```

新たに、必要なimportを行う最は、以下のようなコマンドを実行すると、依存の解決が行われる。

``` bash
# go mod tidy
```

これは、このコマンドを実行する際に、存在する最新バージョンのimport対象のバージョンを含めた依存情報を `go.mod` に追加し、 `go.sum` にそのバージョンのHASHを保持する。

## memberlist

hashiCorpが作ったクラスタ構築ライブラリを使ってみた。

[memberlistを使ってクラスタリング - ノードの参加と削除](https://qiita.com/octu0/items/808299d232bc003d5e99)

``` bash
# docker compose up app1
```

上記のようなコマンドでserver1台が起動する。

その後、別shellでnode2を起動する

``` bash
# docker compose up app2
```

こうすると、二代目に起動した app2 が app1 のクラスタに自動参加される。

app2 を `ctrl+c` で終了すると app1 が検知して離れたことを認識する。

これでお互いに死活監視が行えるようになった。