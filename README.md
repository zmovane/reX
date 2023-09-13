# reX

Reverse Engineered Twitter API

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
var FlagEOF = regexp.MustCompile("^0\\|\\d+$")
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
```
