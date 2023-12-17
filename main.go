package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

func main() {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	// m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	// entries, err := os.ReadDir("./sample")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// _ = make([]fs.FileInfo, 0, len(entries))
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// for _, file := range entries {
	// 	if strings.Contains(file.Name(), ".js") ||
	// 		strings.Contains(file.Name(), ".css") {
	// 		fmt.Println(file.Name())
	// 		if file.IsDir() {
	// 			fmt.Println("dir")
	// 		} else {
	// 			fmt.Println("file")
	// 		}
	// 	}
	// }

	err := filepath.Walk("./sample", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".js") ||
			strings.Contains(path, ".css") {

			fmt.Println(path)
			if info.IsDir() {
				fmt.Println("dir")
			} else {
				fmt.Println("file")
				// pathを読み込んでreader生成、minifyして、writerに書き込む
				reader, err := os.Open(path)
				if err != nil {
					return err
				}
				defer func(reader *os.File) {
					_ = reader.Close()
				}(reader)
				// writerの用意
				err = os.MkdirAll("./dist/"+filepath.Dir(path), 0777)
				if err != nil {
					return err
				}
				writer, err := os.Create("./dist/" + path)
				if err != nil {
					return err
				}
				defer func(writer *os.File) {
					_ = writer.Close()
				}(writer)

				err = m.Minify("text/css", writer, reader)
				if err != nil {
					return err
				}

			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

}
