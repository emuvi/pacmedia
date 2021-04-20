package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func doDigest() {
	digestFolder(body)
}

func digestFolder(folder string) {
	fmt.Println("Digesting: " + folder + "\nStarting..." + "\n-------")
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println("Digesting: " + folder + "\nError: " + err.Error() + "\n-------")
		return
	}
	for _, inside := range files {
		doing := path.Join(folder, inside.Name())
		if inside.IsDir() {
			digestFolder(doing)
		} else {
			waiter.Add(1)
			go digestFile(doing)
		}
	}
}

func digestFile(origin string) {
	defer waiter.Done()
	extension := path.Ext(origin)
	exType := strings.TrimSpace(strings.ToLower(extension))
	if !enabledTypes[exType] {
		return
	}
	fmt.Println("Digesting: " + origin + "\nStarting..." + "\n-------")
	nameBase := strings.TrimSuffix(path.Base(origin), extension)
	destinyTxt := path.Join(path.Dir(origin), nameBase+".txt")
	_, err := os.Stat(destinyTxt)
	if os.IsNotExist(err) {
		cmd := exec.Command("ebook-convert", origin, destinyTxt)
		if err := cmd.Run(); err != nil {
			fmt.Println("Digesting: " + origin + "\nError: " + err.Error() + "\n-------")
			return
		}
		err := cleanText(destinyTxt)
		if err != nil {
			fmt.Println("Digesting: " + origin + "\nError: " + err.Error() + "\n-------")
			return
		}
		fmt.Println("Digesting: " + origin + "\nResult: Successfully transformed." + "\n-------")
	}
	destinyCnt := path.Join(path.Dir(origin), nameBase+".cnt")
	_, err = os.Stat(destinyCnt)
	if os.IsNotExist(err) {
		if err := countWords(destinyTxt, destinyCnt); err != nil {
			fmt.Println("Digesting: " + origin + "\nError: " + err.Error() + "\n-------")
			return
		} else {
			fmt.Println("Digesting: " + origin + "\nResult: Successfully counted." + "\n-------")
		}
	}
	fmt.Println("Digesting: " + origin + "\nResult: Successfully digested." + "\n-------")
}

func cleanText(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	text := string(content)
	re, err := regexp.Compile(`\b\-\s+\b`)
	if err != nil {
		return err
	}
	text = re.ReplaceAllString(text, "")
	writer, err := os.Create(file)
	if err != nil {
		return err
	}
	defer writer.Close()
	writer.WriteString(text)
	writer.Sync()
	return nil
}

func countWords(origin string, destiny string) error {
	file, err := os.Create(destiny)
	if err != nil {
		return err
	}
	defer file.Close()
	words, err := getWords(origin)
	if err != nil {
		return err
	}
	counted := map[string]int{}
	for _, word := range words {
		counted[word]++
	}
	for word, count := range counted {
		file.WriteString(word + "=" + strconv.Itoa(count) + "\n")
	}
	file.Sync()
	return nil
}

func getWords(origin string) ([]string, error) {
	re, err := regexp.Compile(`[\s\x20\xA0\=\,\.\:\;\!\?]+`)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadFile(origin)
	if err != nil {
		return nil, err
	}
	text := string(content)
	split := re.Split(text, -1)
	result := []string{}
	for i := range split {
		candidate := cleanWord(split[i])
		if isValidWord(candidate) {
			result = append(result, candidate)
		}
	}
	return result, nil
}

func cleanWord(word string) string {
	var builder strings.Builder
	for _, r := range word {
		if unicode.IsLetter(r) || r == '-' {
			builder.WriteRune(r)
		}
	}
	result := builder.String()
	result = strings.ToLower(result)
	result = strings.TrimSpace(result)
	return result
}

func isValidWord(word string) bool {
	return len(word) > 0
}
