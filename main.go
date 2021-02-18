package main

import (
	"fmt"

	"github.com/pborman/getopt/v2"
)

var (
	body   = ""
	feed   = ""
	digest = false
	search = ""
	give   = ""
	lend   = ""
	open   = false
	speed  = 8
	help   = false
)

func main() {
	getopt.FlagLong(&body, "body", 'b', "Where all files I eat ends up.")
	getopt.FlagLong(&feed, "feed", 'f', "Yami! More files for me to eat.")
	getopt.FlagLong(&digest, "digest", 'd', "This makes the food in my belly becomes my body.")
	getopt.FlagLong(&search, "search", 's', "Do you wanna me to search in my body and belly?")
	getopt.FlagLong(&give, "give", 'g', "This moves the founds out of me to the destination.")
	getopt.FlagLong(&lend, "lend", 'l', "This copies the founds inside me on the destination.")
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
		panic("You let me as an errant soul, where is my body?")
	}
	fmt.Println("Body:", body)
	fmt.Println("Speed:", speed)
}
