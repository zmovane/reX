package main

import (
	"fmt"
	reversedX "github.com/amovane/reversedX/x"
	"os"
	"regexp"
)

var FlagEOF = regexp.MustCompile("^0\\|\\d+$")

func main() {
	uname := os.Getenv("USER_NAME")
	upwd := os.Getenv("PASSWORD")
	x := reversedX.New(uname, upwd)
	wd, _ := os.Getwd()
	cookiesPath := fmt.Sprintf("%s/cookies.json", wd)
	err := x.SetCookies(cookiesPath)
	if err != nil || !x.IsLoggedIn() {
		println("You must login first")
		x.Login()
		x.SaveCookies(cookiesPath)
	}

	// followings
	var cursor *string
	var EOF bool
	for {
		tweets, nextCursor := x.GetFollowingsByScreenName("shareverse_", cursor)
		EOF = cursor != nil &&
			nextCursor != nil &&
			FlagEOF.Match([]byte(*cursor)) &&
			FlagEOF.Match([]byte(*nextCursor))
		cursor = nextCursor
		if cursor == nil || EOF {
			break
		}
		for _, tweet := range tweets {
			println(tweet.ScreenName)
		}
	}
}
