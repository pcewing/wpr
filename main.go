package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

type Config struct {
	WallpaperDir           string
	DisplayCount, Interval int
}

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

func ReadConfigFile() Config {
	home := os.Getenv("HOME")
	content, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/wpr/wprrc.json", home))
	if err != nil {
		log.Fatal(err)
	}

	contentString := string(content)

	dec := json.NewDecoder(strings.NewReader(contentString))

	var config Config
	if err := dec.Decode(&config); err == io.EOF {

	} else if err != nil {
		log.Fatal(err)
	}

	config.WallpaperDir = os.ExpandEnv(config.WallpaperDir)

	fmt.Printf("Config:\n  Wallpaper Directory: %s\n  Display Count: %d\n  Interval: %d\n",
		config.WallpaperDir,
		config.DisplayCount,
		config.Interval)

	return config
}

func GetFiles(directory string) []string {
	fileNames := []string{}
	subDirectories := []string{}

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		path := path.Join(directory, file.Name())

		if file.IsDir() {
			subDirectories = append(subDirectories, path)
		} else {
			fileNames = append(fileNames, path)
		}
	}

	for _, subDirectory := range subDirectories {
		fileNames = append(fileNames, GetFiles(subDirectory)...)
	}

	return fileNames
}

func SetWallpapers(config Config) {
	files := GetFiles(config.WallpaperDir)
	selectedFile := files[rand.Intn(len(files))]

	img := ReadImageFromFile(selectedFile)

	SetBackgroundX11(img)
}

func main() {
	config := ReadConfigFile()

	rand.Seed(time.Now().UTC().UnixNano())
	for true {
		SetWallpapers(config)
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}
