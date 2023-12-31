package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

func main() {
	fmt.Println("- minify start -")

	inputPath := "./src"
	outputPath := "./dist"

	flag.Parse()
	args := flag.Args()
	if len(args) == 1 {
		fmt.Println("Please specify param2")
		os.Exit(1)
	}
	if len(args) == 2 {
		inputPath = args[0]
		outputPath = args[1]
	}

	if f, err := os.Stat(inputPath); os.IsNotExist(err) || !f.IsDir() {
		fmt.Printf("directory %s does not exist \n", inputPath)
		os.Exit(1)
	}

	//	if f, err := os.Stat(outputPath); os.IsNotExist(err) || !f.IsDir() {
	//		fmt.Printf("directory %s does not exist \n", outputPath)
	//		os.Exit(1)
	//	}

	err := executeMinify(inputPath, outputPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("- minify finish -")
}

func executeMinify(inputPath string, outputPath string) error {
	// setup minify
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	fmt.Println("executeMinify")
	var mediaType string
	err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			mediaType = "text/html"
		} else if strings.Contains(path, ".js") {
			mediaType = "text/javascript"
		} else if strings.Contains(path, ".css") {
			mediaType = "text/css"
		} else if strings.Contains(path, ".json") {
			mediaType = "application/json"
		} else if strings.Contains(path, ".xml") {
			mediaType = "application/xml"
		} else if strings.Contains(path, ".svg") {
			mediaType = "image/svg+xml"
		} else {
			mediaType = ""
		}

		if mediaType == "" {
			return nil
		}

		fmt.Println(path)
		if !info.IsDir() {
			// pathを読み込んでreader生成、minifyして、writerに書き込む
			reader, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func(reader *os.File) {
				_ = reader.Close()
			}(reader)

			// writerの用意
			err = os.MkdirAll(outputPath+"/"+filepath.Dir(path), 0777)
			if err != nil {
				return err
			}
			writer, err := os.Create(outputPath + "/" + path)
			if err != nil {
				return err
			}
			defer func(writer *os.File) {
				_ = writer.Close()
			}(writer)

			err = m.Minify(mediaType, writer, reader)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
