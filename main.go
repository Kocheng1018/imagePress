package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type FileInfo interface {
	Name() string
	Size() int64
	Mode() os.FileMode
	ModTime() time.Time
	IsDir() bool
	Sys() interface{}
}

func main() {
	startTime := time.Now()
	path, err := os.Getwd()
	// cmd1 := exec.Command("rm", "-rf", "out/*")
	// cmd1.Run()
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
	for idx, file := range files {
		fmt.Println("")
		fmt.Printf("剩餘檔案:")
		fmt.Println(len(files) - idx)
		if file.IsDir() {
			listfile(path + "/in/" + file.Name())
		} else {
			contentType, err := checkType(path + "/in/" + file.Name())
			if err != nil {
				panic(err)
			}
			// if (contentType != "image/jpeg") {
			// }
			fmt.Println("contentType: ", contentType)
			fmt.Println("目前檔案:" + file.Name())
			img2webp(file.Name())
		}
	}
}

func img2webp(inPath string) {
	nameArr := strings.Split(inPath, ".")
	name := strings.Join(nameArr[:len(nameArr)-1], "")
	args := []string{"./in/" + inPath, "./out/" + name + ".webp"}
	cmd := exec.Command("./libwebp/bin/cwebp", "-q", "60", args[0], "-o", args[1])
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

func checkType(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	contentType, err := GetFileContentType(f)
	return contentType, err 

}

func GetFileContentType(out *os.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	return contentType, nil
}
