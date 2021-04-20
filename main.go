package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/pborman/getopt/v2"
)

var (
	body   = ""
	feed   = ""
	clean  = false
	digest = false
	search = ""
	give   = ""
	lend   = ""
	open   = false
	speed  = 8
	help   = false
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
	getopt.FlagLong(&body, "body", 'b', "Where all files I eat ends up.")
	getopt.FlagLong(&feed, "feed", 'f', "Yami! More files for me to eat.")
	getopt.FlagLong(&clean, "clean", 'c', "Removes the folders after eating in them.")
	getopt.FlagLong(&digest, "digest", 'd', "This makes the food in my belly becomes my body.")
	getopt.FlagLong(&search, "search", 's', "Do you wanna me to search in my body and belly?")
	getopt.FlagLong(&lend, "lend", 'l', "This copies the founds inside me on the destination.")
	getopt.FlagLong(&give, "give", 'g', "This moves the founds out of me to the destination.")
	getopt.FlagLong(&open, "open", 'o', "This opens the founds inside me.")
	getopt.FlagLong(&speed, "speed", 'e', "How fast I should go.")
	getopt.FlagLong(&help, "help", 'h', "Makes this conversation.")

	getopt.Parse()
	if help {
		fmt.Println("PacMedia - Eats all the files you feed and keeps them organized,")
		fmt.Println("first in the belly, after in the body, for future searchs.")
		getopt.Usage()
		return
	}

	if body == "" {
		body = "./pacbody"
	}
	sts, err := os.Stat(body)
	if os.IsNotExist(err) {
		fmt.Println("My body does not exists on: " + body)
		body = ""
	} else if !sts.IsDir() {
		fmt.Println("My body is not a directory on: " + body)
		body = ""
	}
	if body == "" {
		panic("You let me as an errant soul, where is my body?")
	}

	fmt.Println("Body:", body)
	fmt.Println("Speed:", speed)
	runtime.GOMAXPROCS(speed)

	if feed != "" {
		doFeed()
		waiter.Wait()
	}
	if digest {
		doDigest()
		waiter.Wait()
	}
	if search != "" {
		doSearch()
		waiter.Wait()
	}
	if lend != "" {
		doLend()
		waiter.Wait()
	}
	if give != "" {
		doGive()
		waiter.Wait()
	}
	if open {
		doOpen()
		waiter.Wait()
	}
	fmt.Println("Pacmedia finished this round.")
}
