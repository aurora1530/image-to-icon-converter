# image-to-icon-converter

画像ファイル（PNG, JPG, JPEG）をWindowsアイコン（.ico）に変換するコマンドラインツール

## 使用方法

```bash
# 基本的な使用
go run main.go -i input.png

# 出力ファイル名を指定
go run main.go -i input.png -o custom.ico
```

## 特徴

- 複数サイズ（16x16～256x256）を含む高品質なアイコンファイルを生成
