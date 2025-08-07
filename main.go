package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
	"github.com/sergeymakinen/go-ico"
)

// 比率を保持しながら正方形のキャンバスに画像を配置する関数
func resizeWithAspectRatio(img image.Image, size uint) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 元画像の比率を計算
	var newWidth, newHeight uint
	if width > height {
		// 横長の場合
		newWidth = size
		newHeight = uint(float64(height) * float64(size) / float64(width))
	} else {
		// 縦長または正方形の場合
		newHeight = size
		newWidth = uint(float64(width) * float64(size) / float64(height))
	}

	// 比率を保持してリサイズ
	resized := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	// 正方形のキャンバスを作成（透明な背景）
	canvas := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))

	// キャンバスの中央に配置するための座標を計算
	x := (int(size) - int(newWidth)) / 2
	y := (int(size) - int(newHeight)) / 2

	// リサイズした画像をキャンバスの中央に描画
	draw.Draw(canvas, image.Rect(x, y, x+int(newWidth), y+int(newHeight)), resized, image.Point{}, draw.Src)

	return canvas
}

func main() {
	// コマンドライン引数の定義
	inputFile := flag.String("i", "", "入力ファイル名")
	outputFile := flag.String("o", "", "出力ファイル名")
	flag.Parse()

	// 入力ファイル名が指定されているか確認
	if *inputFile == "" {
		fmt.Fprintln(os.Stderr, "エラー: 入力ファイル名を指定してください。")
		os.Exit(1)
	}

	// 入力ファイルの形式を検証
	ext := filepath.Ext(*inputFile)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		fmt.Fprintln(os.Stderr, "エラー: サポートされていないファイル形式です。png, jpg, jpegのみがサポートされています。")
		os.Exit(1)
	}

	// 画像を読み込む
	imgFile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: 入力ファイルを開けません: %v\n", err)
		os.Exit(1)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: 画像をデコードできません: %v\n", err)
		os.Exit(1)
	}

	// 出力ファイル名を決定
	if *outputFile == "" {
		baseName := filepath.Base(*inputFile)
		ext := filepath.Ext(baseName)
		*outputFile = baseName[:len(baseName)-len(ext)] + ".ico"
	}

	// 出力ファイルを作成
	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: 出力ファイルを作成できません: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// アイコンに変換して保存（複数サイズ対応）
	// Windowsエクスプローラーで適切に表示されるよう、複数のサイズを含める
	sizes := []uint{16, 32, 48, 64, 128, 256}

	var images []image.Image
	for _, size := range sizes {
		// 比率を保持しながら画像をリサイズ
		resized := resizeWithAspectRatio(img, size)
		images = append(images, resized)
	}

	err = ico.EncodeAll(outFile, images)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: アイコン形式にエンコードできません: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("'%s' を '%s' に変換しました。\n", *inputFile, *outputFile)
}
