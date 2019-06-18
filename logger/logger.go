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

	fileName := getFileName()
	file := getFile(fileName)
	fmt.Println("Setting output to " + file.Name())
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
}

func getFileName() string {

	now := time.Now()
	month := now.Month().String()
	day := strconv.Itoa(now.Day())
	dayOfWeek := now.Weekday().String()

	fileName := dayOfWeek + "(" + day + " " + month + ").log"
	return "./logs/" + fileName
}

func getFile(fileName string) (file *os.File) {

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
