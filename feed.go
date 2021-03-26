package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func feedFolder(folder string) {
	fmt.Println("Feeding folder: " + folder)
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		panic(err)
	}
	for _, inside := range files {
		doing := path.Join(folder, inside.Name())
		if inside.IsDir() {
			feedFolder(doing)
		} else {
			feedFile(doing)
		}
	}
	os.Remove(folder)
}

func feedFile(file string) {
	fmt.Println("Feeding file: " + file)
	sts, err := os.Stat(file)
	if os.IsNotExist(err) {
		panic("The feed file does not exists: " + file)
	}
	if sts.Size() == 0 {
		fmt.Println("The file is empty: " + file)
		return
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	check := fmt.Sprintf("%x", md5.Sum(data))
	fmt.Println("File Check Sum: " + check)
	root := path.Join(body, check[0:2], check[2:4])
	destiny := path.Join(root, check+path.Ext(file))
	fmt.Println("Destiny: " + destiny)
	_, err = os.Stat(destiny)
	if os.IsNotExist(err) {
		os.MkdirAll(root, os.ModePerm)
		err = os.Rename(file, destiny)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("This file is already on my belly.")
		return
	}
}

func doFeed() {
	sts, err := os.Stat(feed)
	if os.IsNotExist(err) {
		panic("The feed path does not exists: " + feed)
	}
	if sts.IsDir() {
		feedFolder(feed)
	} else {
		feedFile(feed)
	}
}
