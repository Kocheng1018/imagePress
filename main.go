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
	options     bimg.Options
	fileName string
}

func listfile(path string) {
	wg := new(sync.WaitGroup)

	files, _ := ioutil.ReadDir("in")
	options := bimg.Options{
		Quality: 60,
		Type:    bimg.ImageType(bimg.WEBP),
	}

	p, _ := ants.NewPoolWithFunc(5, func(in interface{}) {
		st := in.(TmpStruct)
		imagePress(st.options, st.fileName)
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

		p.Invoke(TmpStruct{
			fileName: file.Name(),
			options:     options,
		})
	}
	wg.Wait()
}

func imagePress(options bimg.Options, fileName string) {

	buffer, err := bimg.Read("./in/" + fileName)
	if err != nil {
		spew.Dump(os.Stderr, err)
	}

	newFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	imageType := bimg.NewImage(buffer).Type()
	if imageType != "jpeg" &&
		imageType != "heif" &&
		imageType != "webp" &&
		imageType != "png" {
		spew.Dump(newFileName)
		return
	}

	newImage, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		spew.Dump(os.Stderr, err)
	}

	bimg.Write("./out/"+newFileName+".webp", newImage)
	bimg.VipsCacheDropAll()
}
