package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/h2non/bimg"
)

func main() {
	startTime := time.Now()
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	listfile(path)
	timeBlock := time.Since(startTime)
	fmt.Println("")
	fmt.Println("所需時間:", timeBlock)
}

func listfile(path string) {
	files, _ := ioutil.ReadDir(path + "/in")
	for _, file := range files {
		if file.Name() == "README.md" {
			continue
		}
		if file.IsDir() {
			listfile(path + "/in/" + file.Name())
		} else {
			imagePress(file.Name())
		}
	}
}

func imagePress(url string) {

	buffer, err := bimg.Read("./in/" + url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	options := bimg.Options{
		Quality: 60,
	}

	newImage, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println("beforeType:", bimg.NewImage(newImage).Type())

	convertImage, err := bimg.NewImage(newImage).Convert(bimg.WEBP)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println("afterType:", bimg.NewImage(convertImage).Type())

	bimg.Write("./out/test.webp", convertImage)
}
