package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func Init() {
	log.Println("Configuring log output")
	configureLogger()
}

func configureLogger() {

	file := GetOrCreateFile("app")
	fmt.Println("Setting output to " + file.Name())
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
}

func createFileName(suffix string) string {

	now := time.Now()
	month := now.Month().String()
	day := strconv.Itoa(now.Day())
	dayOfWeek := now.Weekday().String()

	fileName := dayOfWeek + "(" + day + " " + month + ").log"
	return "./logs/" + suffix + "-" + fileName
}

func GetOrCreateFile(suffix string) (file *os.File) {

	fileName := createFileName(suffix)

	_, err := os.Stat(fileName)
	if err != nil && os.IsNotExist(err) {
		fmt.Println("No file found for today. Creating new one")
		file, err = os.Create(fileName)
	} else {
		file, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	}

	if err != nil {
		fmt.Println("Error while opening or creating a file")
		fmt.Println(err)
		os.Exit(3)
	}
	return file
}
