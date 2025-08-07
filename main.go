package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
	"github.com/sergeymakinen/go-ico"
)

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
		// 画像をリサイズ
		resized := resize.Resize(size, size, img, resize.Lanczos3)
		images = append(images, resized)
	}
	
	err = ico.EncodeAll(outFile, images)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: アイコン形式にエンコードできません: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("'%s' を '%s' に変換しました。\n", *inputFile, *outputFile)
}
