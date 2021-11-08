package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
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
	spew.Dump("所需時間:", timeBlock)
}

func listfile(path string) {
	wg := new(sync.WaitGroup)
	files, _ := ioutil.ReadDir(path + "/in")
	options := bimg.Options{
		Quality:     60,
		Type:        bimg.ImageType(bimg.WEBP),
		// Compression: 90,
		// Speed:       8,
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
				spew.Dump(os.Stderr, err)
			}
			newfileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			wg.Add(1)
			go imagePress(buffer, options, newfileName, wg)
			// imagePress(buffer, options, newfileName)
		}
	}
	wg.Wait()
}

func imagePress(buffer []byte, options bimg.Options, newFileName string, wg *sync.WaitGroup) {
	defer wg.Done()

	if bimg.NewImage(buffer).Type() != "jpeg" &&
		bimg.NewImage(buffer).Type() != "heif" &&
		bimg.NewImage(buffer).Type() != "webp" &&
		bimg.NewImage(buffer).Type() != "png" {
		spew.Dump(newFileName)
		return
	}

	newImage, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		spew.Dump(os.Stderr, err)
	}

	bimg.Write("./out/"+newFileName+".webp", newImage)
}
