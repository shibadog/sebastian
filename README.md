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
