package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var logFile *os.File
var logWriter *csv.Writer
var logChan chan []string
var logSendWaiter *sync.WaitGroup
var logWriterWaiter *sync.WaitGroup

func startLogWriter() {
	logFolder := path.Join(bodyParam, "(logs)")
	os.MkdirAll(logFolder, os.ModePerm)
	logName := path.Join(logFolder, time.Now().Format("2006-01-02-15-04-05")+".csv")
	logFile, err := os.OpenFile(logName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	logWriter = csv.NewWriter(logFile)
	logChan = make(chan []string, 10*speedParam)
	logSendWaiter = &sync.WaitGroup{}
	logWriterWaiter = &sync.WaitGroup{}
	logWriterWaiter.Add(1)
	go writeLog()
}

func writeLog() {
	defer logWriterWaiter.Done()
	for lines := range logChan {
		logWriter.Write(lines)
	}
}

func closeLogWriter() {
	fmt.Println("Closing Log Writer...")
	logSendWaiter.Wait()
	close(logChan)
	logWriterWaiter.Wait()
	logWriter.Flush()
	logFile.Close()
	fmt.Println("Closed Log Writer.")
}

func sendToWrite(lines ...string) {
	logChan <- lines
	logSendWaiter.Done()
}

func pacLog(lines ...string) {
	lines = append([]string{time.Now().Format("15-04-05.000")}, lines...)
	if recordParam {
		logSendWaiter.Add(1)
		go sendToWrite(lines...)
	}
	lines = append(lines, "----------------")
	fmt.Println(strings.Join(lines, "\n"))
}

func fixPath(pathToFix string) string {
	pathToFix = path.Clean(pathToFix)
	if path.IsAbs(pathToFix) {
		return pathToFix
	}
	homeDirInitial := "~" + string(os.PathSeparator)
	if strings.HasPrefix(pathToFix, homeDirInitial) {
		uhd, err := os.UserHomeDir()
		if err != nil {
			return pathToFix
		}
		pathToFix = strings.TrimPrefix(pathToFix, homeDirInitial)
		return path.Join(uhd, pathToFix)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return pathToFix
		}
		return path.Join(cwd, pathToFix)
	}
}

func includeInName(ofFile string, theName string) {
	dir := path.Dir(ofFile)
	base := path.Base(ofFile)
	ext := path.Ext(base)
	name := strings.TrimSuffix(base, ext)
	index := 1
	for {
		newName := name + theName
		if index > 1 {
			newName = newName + " (" + strconv.Itoa(index) + ")"
		}
		newName = newName + ext
		destinyFile := path.Join(dir, newName)
		_, err := os.Stat(destinyFile)
		if os.IsNotExist(err) {
			os.Rename(ofFile, destinyFile)
			break
		}
		index++
	}
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}

var sourceRandom = rand.NewSource(time.Now().UnixNano())
var lettersRandom = []rune("0123456789abcdefghijklmnopqrstuvwxyz")

func getRandomPassword(ofSize int) string {
	runes := make([]rune, ofSize)
	for i := range runes {
		runes[i] = lettersRandom[sourceRandom.Int63()%int64(len(lettersRandom))]
	}
	return string(runes)
}
