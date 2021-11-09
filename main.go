package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/h2non/bimg"
	"github.com/panjf2000/ants/v2"
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

type TmpStruct struct {
	buffer      []byte
	options     bimg.Options
	newFileName string
}

func listfile(path string) {
	wg := new(sync.WaitGroup)

	files, _ := ioutil.ReadDir("in")
	options := bimg.Options{
		Quality: 60,
		Type:    bimg.ImageType(bimg.WEBP),
		// Compression: 90,
		// Speed:       8,
	}

	p, _ := ants.NewPoolWithFunc(10, func(in interface{}) {
		st := in.(TmpStruct)
		imagePress(st.buffer, st.options, st.newFileName)
		wg.Done()
	})

	defer p.Release()

	wg.Add(len(files))
	for _, file := range files {
		// if file.Name() == "README.md" {
		// 	continue
		// }
		// if file.IsDir() {
		// 	listfile(path + "/in/" + file.Name())
		// 	continue
		// }
		spew.Dump(fmt.Sprintf("run:%s", file.Name()))

		buffer, err := bimg.Read("./in/" + file.Name())
		if err != nil {
			spew.Dump(os.Stderr, err)
		}

		newfileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		p.Invoke(TmpStruct{
			buffer:      buffer,
			newFileName: newfileName,
			options:     options,
		})

		// syncCalculateSum := func() {
		// 	imagePress(buffer, options, newfileName)
		// 	wg.Done()
		// }
		// _ = ants.Submit(syncCalculateSum)
	}
	wg.Wait()
}

func imagePress(buffer []byte, options bimg.Options, newFileName string) {
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
