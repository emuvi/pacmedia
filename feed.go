package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var feedWaiter *sync.WaitGroup
var feedParamLength int
var feedSuccess uint32
var feedDuplicate uint32
var feedError uint32
var filesToFeed chan string
var foldersToClean []string

func doFeed() {
	feedParam = fixPath(feedParam)
	pacLog("Feeding: "+feedParam, "Feed Starting...")
	if strings.HasPrefix(bodyParam, feedParam) {
		pacLog("Feeding: "+feedParam, "Error: The body can't be inside the feed.")
		return
	}
	if strings.HasPrefix(feedParam, bodyParam) {
		pacLog("Feeding: "+feedParam, "Error: The feed can't be inside the body.")
		return
	}
	sts, err := os.Stat(feedParam)
	if os.IsNotExist(err) {
		pacLog("Feeding: "+feedParam, "Error: The path does not exists.")
		return
	}
	feedParamLength = len(feedParam)
	feedWaiter = &sync.WaitGroup{}
	feedSuccess = 0
	feedDuplicate = 0
	feedError = 0
	filesToFeed = make(chan string, 2*speedParam)
	for i := 0; i < speedParam; i++ {
		feedWaiter.Add(1)
		go feedFile()
	}
	if sts.IsDir() {
		feedFolder(feedParam)
	} else {
		filesToFeed <- feedParam
	}
	close(filesToFeed)
	pacLog("Feed: Closed files to feed.")
	feedWaiter.Wait()
	if cleanParam {
		for _, folder := range foldersToClean {
			os.Remove(folder)
		}
	}
	pacLog("Feed: Terminated.",
		"Success: "+strconv.Itoa(int(feedSuccess)),
		"Duplicate: "+strconv.Itoa(int(feedDuplicate)),
		"Error: "+strconv.Itoa(int(feedError)))
}

func feedFolder(folder string) {
	display := "[f]" + folder[feedParamLength:]
	pacLog("Feeding: "+display, "Folder Starting...")
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		pacLog("Feeding: "+folder, "Error: "+err.Error())
		return
	}
	for _, inside := range files {
		doing := path.Join(folder, inside.Name())
		if inside.IsDir() {
			feedFolder(doing)
		} else {
			filesToFeed <- doing
		}
	}
	if cleanParam {
		foldersToClean = append(foldersToClean, folder)
	}
}

func feedFile() {
	defer feedWaiter.Done()
	for origin := range filesToFeed {
		display := "[f]" + origin[feedParamLength:]
		pacLog("Feeding: "+display, "File Starting...")
		exType := strings.TrimSpace(strings.ToLower(path.Ext(origin)))
		if !enabledTypes[exType] {
			if cleanParam {
				includeInName(origin, " (disabled)")
			}
			pacLog("Feeding: "+display, "Error: It's not an enabled type.")
			return
		}
		sts, err := os.Stat(origin)
		if os.IsNotExist(err) {
			pacLog("Feeding: "+display, "Error: The file does not exists.")
			return
		}
		if sts.Size() == 0 {
			pacLog("Feeding: "+display, "Error: The file is empty.")
			return
		}
		file, err := os.Open(origin)
		if err != nil {
			pacLog("Feeding: "+display, "Error: "+err.Error())
			return
		}
		hash := sha256.New()
		_, err = io.Copy(hash, file)
		file.Close()
		if err != nil {
			pacLog("Feeding: "+display, "Error: "+err.Error())
			return
		}
		check := fmt.Sprintf("%x", hash.Sum(nil))
		root := path.Join(bodyParam, check[0:3], check[3:6], check)
		pass := getRandomPassword(18)
		ext := path.Ext(origin)
		destiny := path.Join(root, "org-"+pass+ext)
		_, err = os.Stat(root)
		if os.IsNotExist(err) {
			os.MkdirAll(root, os.ModePerm)
			err = os.Rename(origin, destiny)
			if err != nil {
				err = moveFile(origin, destiny)
				if err != nil {
					if cleanParam {
						includeInName(origin, " (error)")
					}
					pacLog("Feeding: "+display, "Checker: "+check,
						"Error: "+err.Error())
					atomic.AddUint32(&feedError, 1)
					os.Remove(root)
					return
				}
			}
		} else {
			if cleanParam {
				includeInName(origin, " (duplicate)")
			}
			pacLog("Feeding: "+display, "Checker: "+check,
				"Error: This file is already on my belly.")
			atomic.AddUint32(&feedDuplicate, 1)
			return
		}
		if cleanParam {
			doing := path.Dir(origin)
			err = os.Remove(doing)
			for err == nil && len(doing) > feedParamLength {
				doing = path.Dir(doing)
				err = os.Remove(doing)
			}
		}
		pacLog("Feeding: "+display, "Checker: "+check, "Success!")
		atomic.AddUint32(&feedSuccess, 1)
	}
}
