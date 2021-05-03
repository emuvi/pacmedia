package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/pborman/getopt/v2"
)

var (
	recordParam = false
	bodyParam   = ""
	feedParam   = ""
	cleanParam  = false
	digestParam = false
	searchParam = ""
	giveParam   = ""
	lendParam   = ""
	openParam   = false
	speedParam  = 8
	helpParam   = false
	versParam   = false
)

var waiter sync.WaitGroup

var enabledTypes = map[string]bool{
	".pdf": true, ".pdb": true,
	".epub": true, ".htmlz": true,
	".mobi": true, ".azw3": true,
	".rtf": true, ".odt": true,
	".doc": true, ".docx": true,
}

func main() {
	getopt.FlagLong(&recordParam, "rec", 'r', "Records all the logs.") // TODO Record logs.
	getopt.FlagLong(&bodyParam, "body", 'b', "Where all files I eat ends up.")
	getopt.FlagLong(&feedParam, "feed", 'f', "Yami! More files for me to eat.")
	getopt.FlagLong(&cleanParam, "clean", 'c', "Removes the folders after eating in them.")
	getopt.FlagLong(&digestParam, "digest", 'd', "This makes the food in my belly becomes my body.")
	getopt.FlagLong(&searchParam, "search", 's', "Do you wanna me to search in my body and belly?")
	getopt.FlagLong(&lendParam, "lend", 'l', "This copies the founds inside me on the destination.")
	getopt.FlagLong(&giveParam, "give", 'g', "This moves the founds out of me to the destination.")
	getopt.FlagLong(&openParam, "open", 'o', "This opens the founds inside me.")
	getopt.FlagLong(&speedParam, "speed", 'e', "How fast I should go.")
	getopt.FlagLong(&helpParam, "help", 'h', "Makes this conversation.")
	getopt.FlagLong(&versParam, "version", 'v', "Show the current version.")

	getopt.Parse()
	if helpParam {
		fmt.Println("PacMedia - Eats all the files you feed and keeps them organized,")
		fmt.Println("first in the belly, after in the body, for future searchs.")
		getopt.Usage()
		return
	}
	if versParam {
		fmt.Println("PacMedia - Version: 0.1.5")
		return
	}

	if bodyParam == "" {
		bodyParam = "./pacbody"
	}
	sts, err := os.Stat(bodyParam)
	if os.IsNotExist(err) {
		fmt.Println("My body does not exists on: " + bodyParam)
		bodyParam = ""
	} else if !sts.IsDir() {
		fmt.Println("My body is not a directory on: " + bodyParam)
		bodyParam = ""
	}
	if bodyParam == "" {
		panic("You let me as an errant soul, where is my body?")
	}
	bodyParam = fixPath(bodyParam)

	fmt.Println("Body:", bodyParam)
	fmt.Println("Speed:", speedParam)
	runtime.GOMAXPROCS(speedParam)

	if feedParam != "" {
		doFeed()
	}
	if digestParam {
		doDigest()
	}
	if searchParam != "" {
		doSearch()
	}
	if lendParam != "" {
		doLend()
	}
	if giveParam != "" {
		doGive()
	}
	if openParam {
		doOpen()
	}
	fmt.Println("Pacmedia finished this round.")
}
