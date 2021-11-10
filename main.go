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
	bimg.VipsCacheSetMaxMem(2048)
	bimg.VipsCacheSetMax(2048)
	if err != nil {
		panic(err)
	}
	listfile(path)
	timeBlock := time.Since(startTime)
	spew.Dump("所需時間:", timeBlock)
}

type TmpStruct struct {
	fileName    string
	newFileName string
}

var options bimg.Options = bimg.Options{
	Quality: 60,
	Type:    bimg.ImageType(bimg.WEBP),
}

// var optionsQuality bimg.Options = bimg.Options{
// 	Quality: 60,
// }

func listfile(path string) {
	wg := new(sync.WaitGroup)

	files, _ := ioutil.ReadDir("in")

	p, _ := ants.NewPoolWithFunc(4, func(in interface{}) {
		st := in.(TmpStruct)
		imagePress(st.fileName, st.newFileName)
		wg.Done()
	})

	defer p.Release()

	wg.Add(len(files))
	for _, file := range files {
		spew.Dump(fmt.Sprintf("run:%s", file.Name()))

		newFileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		p.Invoke(TmpStruct{
			fileName:    file.Name(),
			newFileName: newFileName,
		})
	}
	wg.Wait()
}

func imagePress(fileName string, newFileName string) {

	buffer, err := bimg.Read("./in/" + fileName)
	if err != nil {
		spew.Dump(os.Stderr, err)
	}

	image := bimg.NewImage(buffer)
	var imageByte []byte

	switch image.Type() {
	case "jpeg", "png", "heif", "webp":
		imageByte, _ = image.Process(options)
	default:
		return
	}

	bimg.Write("./out/"+newFileName+".webp", imageByte)
}
