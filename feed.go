package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var foldersToClean []string

func doFeed() {
	feedParam = fixPath(feedParam)
	if strings.HasPrefix(bodyParam, feedParam) {
		fmt.Println("Feeding: " + feedParam + "\nError: The body can't be inside the feed." + "\n-------")
		return
	}
	if strings.HasPrefix(feedParam, bodyParam) {
		fmt.Println("Feeding: " + feedParam + "\nError: The feed can't be inside the body." + "\n-------")
		return
	}
	sts, err := os.Stat(feedParam)
	if os.IsNotExist(err) {
		fmt.Println("Feeding: " + feedParam + "\nError: The path does not exists." + "\n-------")
		return
	}
	if sts.IsDir() {
		feedFolder(feedParam)
	} else {
		waiter.Add(1)
		go feedFile(feedParam)
	}
	waiter.Wait()
	if cleanParam {
		for _, folder := range foldersToClean {
			os.Remove(folder)
		}
	}
}

func feedFolder(folder string) {
	fmt.Println("Feeding: " + folder + "\nStarting..." + "\n-------")
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println("Feeding: " + folder + "\nError: " + err.Error() + "\n-------")
		return
	}
	for _, inside := range files {
		doing := path.Join(folder, inside.Name())
		if inside.IsDir() {
			feedFolder(doing)
		} else {
			waiter.Add(1)
			go feedFile(doing)
		}
	}
	if cleanParam {
		foldersToClean = append(foldersToClean, folder)
	}
}

func feedFile(origin string) {
	defer waiter.Done()
	fmt.Println("Feeding: " + origin + "\nStarting..." + "\n-------")
	exType := strings.TrimSpace(strings.ToLower(path.Ext(origin)))
	if !enabledTypes[exType] {
		fmt.Println("Feeding: " + origin + "\nError: It's not an enabled type." + "\n-------")
		return
	}
	sts, err := os.Stat(origin)
	if os.IsNotExist(err) {
		fmt.Println("Feeding: " + origin + "\nError: The file does not exists." + "\n-------")
		return
	}
	if sts.Size() == 0 {
		fmt.Println("Feeding: " + origin + "\nError: The file is empty." + "\n-------")
		return
	}
	file, err := os.Open(origin)
	if err != nil {
		fmt.Println("Feeding: " + origin + "\nError: " + err.Error() + "\n-------")
		return
	}
	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		fmt.Println("Feeding: " + origin + "\nError: " + err.Error() + "\n-------")
		return
	}
	check := fmt.Sprintf("%x", hash.Sum(nil))
	root := path.Join(bodyParam, check[0:2], check[2:4])
	destiny := path.Join(root, check+path.Ext(origin))
	fmt.Println("Feeding: " + origin + "\nDestiny: " + destiny + "\n-------")
	_, err = os.Stat(destiny)
	if os.IsNotExist(err) {
		os.MkdirAll(root, os.ModePerm)
		err = os.Rename(origin, destiny)
		if err != nil {
			fmt.Println("Feeding: " + origin + "\nError: " + err.Error() + "\n-------")
			return
		}
	} else {
		fmt.Println("Feeding: " + origin + "\nError: This file is already on my belly." + "\n-------")
		return
	}
	fmt.Println("Feeding: " + origin + "\nResult: Successfully eaten." + "\n-------")
}
