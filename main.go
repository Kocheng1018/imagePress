package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	wg := new(sync.WaitGroup)
	files, _ := ioutil.ReadDir(path + "/in")
	options := bimg.Options{
		Quality: 60,
	}
	for _, file := range files {
		if file.Name() == "README.md" {
			continue
		}
		if file.IsDir() {
			listfile(path + "/in/" + file.Name())
		} else {
			buffer, err := bimg.Read("./in/" + file.Name())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			newfileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			wg.Add(1)
			go imagePress(buffer, options, newfileName, wg)
		}
	}
	wg.Wait()
}

func imagePress(buffer []byte, options bimg.Options, newFileName string, wg *sync.WaitGroup) {
	defer wg.Done()

	newImage, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	convertImage, err := bimg.NewImage(newImage).Convert(bimg.WEBP)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	bimg.Write("./out/"+newFileName+".webp", convertImage)
}
