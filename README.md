# JSON to Go Struct Converter CLI

このCLIツールは、JSONファイルをGoの構造体（struct）に変換するためのものです。指定したJSONファイルを解析し、Goで使用できる構造体を自動的に生成します。

## 特徴

- JSONファイルからGoの構造体を自動生成
- CamelCase変換によるフィールド名の最適化
- 自動的にjsonタグを生成
- コマンドライン引数で構造体名や出力ファイルの指定が可能

## インストール

1. Goがインストールされていることを確認してください（バージョン1.23以降推奨）。
2. このリポジトリをクローン、またはソースコードをダウンロードします。
3. CLIツールをビルドします：

   ```bash
   go build -o json2struct json2struct.go
