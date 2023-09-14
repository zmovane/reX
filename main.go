package main

import (
	"fmt"
	reX "github.com/amovane/reX/x"
	"os"
)

func main() {
	uname := os.Getenv("USER_NAME")
	upwd := os.Getenv("PASSWORD")
	x := reX.New(uname, upwd)
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
	var tweets []reX.Legacy
	for {
		tweets, cursor, err = x.GetFollowingsByScreenName("shareverse_", cursor)
		for _, tweet := range tweets {
			println(tweet.ScreenName)
		}
		if cursor == nil || err != nil {
			break
		}
	}
}
