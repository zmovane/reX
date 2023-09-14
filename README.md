# reX

[![go report card](https://goreportcard.com/badge/github.com/amovane/reX "go report card")](https://goreportcard.com/report/github.com/amovane/reX)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

<p align="center">
    <img src="./twitter.png" height=150></img>
</p>

Since twitter dev removed the API for accessing user followers and following, developers have found it difficult to obtain this data. Here, I'm sharing my reverse engineering solution. If it's helpful to you, please let me know.

## Login

```golang
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
```

## Followings

```golang
var cursor *string
var tweets []reX.Legacy
for {
    tweets, cursor, err = x.GetFollowingsByScreenName("shareverse_", cursor)
    if cursor == nil || err != nil {
        break
    }
    for _, tweet := range tweets {
        println(tweet.ScreenName)
    }
}
```

## Followers

```golang
var cursor *string
var tweets []reX.Legacy
for {
    tweets, cursor, err = x.GetFollowersByScreenName("shareverse_", cursor)
    if cursor == nil || err != nil {
        break
    }
    for _, tweet := range tweets {
        println(tweet.ScreenName)
    }
}
```
