package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Config struct {
	WallpaperDir           string
	DisplayCount, Interval int
}

func GetFiles(directory string) []string {
	fileNames := []string{}
	subDirectories := []string{}

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		path := fmt.Sprintf("%s/%s", directory, file.Name())

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

	cmd := "feh"
	args := []string{"--no-fehbg", "--bg-scale"}

	for i := 0; i < config.DisplayCount; i++ {
		index := rand.Int() % len(files)
		args = append(args, files[index])
	}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ReadConfigFile() Config {
	content, err := ioutil.ReadFile("/home/paul/.config/wpr/wprrc.json")
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

	fmt.Printf("Config:\n  Wallpaper Directory: %s\n  Display Count: %d\n  Interval: %d\n",
		config.WallpaperDir,
		config.DisplayCount,
		config.Interval)

	return config
}

func main() {
	config := ReadConfigFile()

	for true {
		SetWallpapers(config)
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}
