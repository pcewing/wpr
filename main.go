package main

import (
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"time"
)

func ReadImageFromFile(path string) image.Image {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return m
}

func SelectBackgroundImage() string {
	dir := "/home/paul/go/src/github.com/pcewing/wpr2/data"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("Failed to read files")
	}

	rand.Seed(time.Now().UTC().UnixNano())
	fileInfo := files[rand.Intn(len(files))]

	return path.Join(dir, fileInfo.Name())
}

func main() {
	filePath := SelectBackgroundImage()
	img := ReadImageFromFile(filePath)

	SetBackgroundX11(img)
}
