package main

import (
	"fmt"
	reX "github.com/amovane/reX/x"
	"os"
)

func main() {
	uname := os.Getenv("USER_NAME")
	upwd := os.Getenv("PASSWORD")
	x := reX.New()
	wd, _ := os.Getwd()
	cookiesPath := fmt.Sprintf("%s/cookies.json", wd)
	err := x.SetCookies(cookiesPath)
	if err != nil || !x.IsLoggedIn() {
		println("You must login first")
		x.Login(uname, upwd)
		x.SaveCookies(cookiesPath)
	}

	// followings
	var cursor *string
	var tweets []reX.UserResults
	for {
		tweets, cursor, err = x.GetFollowingsByScreenName("shareverse_", cursor)
		for _, tweet := range tweets {
			println(tweet.Result.Legacy.ScreenName)
		}
		if cursor == nil || err != nil {
			break
		}
	}
	uid, _ := x.GetUserIDByScreenName("jarredsumner")
	uidOfFollower, _ := x.GetUserIDByScreenName("cybermageek")
	isFollowing := x.IsFollowing(uid, uidOfFollower)
	println("isFollowing", isFollowing)
}
