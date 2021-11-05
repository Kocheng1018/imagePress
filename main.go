package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/h2non/filetype"
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
	for _, file := range files {
		if file.Name() == "README.md" {
			continue
		}
		if file.IsDir() {
			listfile(path + "/in/" + file.Name())
		} else {
			contentType := checkType(file.Name())
			fmt.Println("contentType:", contentType)
			fmt.Println("目前檔案:" + file.Name() + "\n")
			wg.Add(1)
			go img2webp(file.Name(), wg)
		}
	}
	wg.Wait()
}

func img2webp(inPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	nameArr := strings.Split(inPath, ".")
	name := strings.Join(nameArr[:len(nameArr)-1], "")
	args := []string{"./in/" + inPath, "./out/" + name + ".webp"}
	cmd := exec.Command("./libwebp/bin/cwebp", "-q", "60", args[0], "-o", args[1])
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

func checkType(path string) string {
	buf, _ := ioutil.ReadFile("./in/" + path)

	kind, _ := filetype.Match(buf)
	if kind == filetype.Unknown {
		fmt.Println("Unknown file type")
		return ""
	}
	return kind.Extension
}
