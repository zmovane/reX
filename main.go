package main

import (
	"fmt"
	"os"
	"regexp"
)

var FlagEOF = regexp.MustCompile("^0\\|\\d+$")

func main() {
	uname := os.Getenv("USER_NAME")
	upwd := os.Getenv("PASSWORD")
	reversedX := New(uname, upwd)
	reversedX.login()
	if !reversedX.IsLoggedIn() {
		panic("You must login first")
	}
	var cursor *string
	var EOF bool
	for {
		tweets, nextCursor := reversedX.GetFollowingsByScreenName("shareverse_", cursor)
		EOF = cursor != nil &&
			nextCursor != nil &&
			FlagEOF.Match([]byte(*cursor)) &&
			FlagEOF.Match([]byte(*nextCursor))
		cursor = nextCursor
		if cursor == nil || EOF {
			println("=====================END=======================")
			break
		}
		for _, tweet := range tweets {
			println(tweet.ScreenName)
		}
		fmt.Printf("count<%v>====================cursor<%s>=======================\n", len(tweets), *cursor)
	}
}
