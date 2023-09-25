package reX

import (
	"encoding/json"
	"fmt"
	"net/http"
	urlutil "net/url"
	"os"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
	SCRAPER "github.com/n0madic/twitter-scraper"
)

var FlagEOF = regexp.MustCompile("^0\\|\\d+$")
var DEFAULT_FEATURES = T{
	"rweb_lists_timeline_redesign_enabled":                                    true,
	"responsive_web_graphql_exclude_directive_enabled":                        true,
	"verified_phone_label_enabled":                                            false,
	"creator_subscriptions_tweet_preview_api_enabled":                         true,
	"responsive_web_graphql_timeline_navigation_enabled":                      true,
	"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
	"tweetypie_unmention_optimization_enabled":                                true,
	"responsive_web_edit_tweet_api_enabled":                                   true,
	"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
	"view_counts_everywhere_api_enabled":                                      true,
	"longform_notetweets_consumption_enabled":                                 true,
	"responsive_web_twitter_article_tweet_consumption_enabled":                false,
	"tweet_awards_web_tipping_enabled":                                        false,
	"freedom_of_speech_not_reach_fetch_enabled":                               true,
	"standardized_nudges_misinfo":                                             true,
	"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
	"longform_notetweets_rich_text_read_enabled":                              true,
	"longform_notetweets_inline_media_enabled":                                true,
	"responsive_web_media_download_video_enabled":                             false,
	"responsive_web_enhance_cards_enabled":                                    false,
}

type X struct {
	uname     string
	upwd      string
	scraper   *SCRAPER.Scraper
	cookies   string
	csrfToken string
}

func New(uname string, upwd string) X {
	return X{
		uname:   uname,
		upwd:    upwd,
		scraper: SCRAPER.New(),
	}
}

func (x *X) Login() error {
	return x.scraper.Login(x.uname, x.upwd)
}

func (x *X) IsLoggedIn() bool {
	return x.scraper.IsLoggedIn()
}

func (x *X) SetCookies(cookiesPath string) (err error) {
	var cookies []*http.Cookie
	var f *os.File
	if f, err = os.Open(cookiesPath); err != nil {
		return err
	}
	json.NewDecoder(f).Decode(&cookies)
	x.scraper.SetCookies(cookies)
	return nil
}

func (x *X) SaveCookies(cookiesPath string) {
	cookies := x.scraper.GetCookies()
	data, _ := json.Marshal(cookies)
	os.WriteFile(cookiesPath, data, 0644)
}

func (x *X) GetRelationsById(uid string, cursor *string, relation Relation) (resp []UserResults, nextCursor *string, err error) {
	cookies := Map(x.scraper.GetCookies(), func(field *http.Cookie) string {
		if field.Name == "ct0" {
			x.csrfToken = field.Value
		}
		return field.String()
	})
	x.cookies = strings.Join(cookies, ";")
	variables := T{
		"userId":                 uid,
		"count":                  100,
		"includePromotedContent": false,
	}
	if cursor != nil {
		variables["cursor"] = *cursor
	}
	variablesJson, _ := json.Marshal(variables)
	featuresJson, _ := json.Marshal(DEFAULT_FEATURES)
	query := fmt.Sprintf(`variables=%s&features=%s`, variablesJson, featuresJson)
	values, _ := urlutil.ParseQuery(query)
	url := fmt.Sprintf(`https://twitter.com/i/api/graphql/%s?%s`, relation.Path(), values.Encode())
	response, err := x.RequestGet(url)
	instructions := response.Data.User.Result.Timeline.Timeline.Instructions
	resp = make([]UserResults, 0)
	for _, i := range instructions {
		for _, e := range i.Entries {
			cursorType := e.Content.CursorType
			if cursorType != nil && *cursorType == "Bottom" {
				nextCursor = e.Content.Value
			}
			item := e.Content.ItemContent
			if item == nil {
				continue
			}
			resp = append(resp, item.UserResults)
		}
	}
	EOF := cursor != nil &&
		nextCursor != nil &&
		FlagEOF.Match([]byte(*cursor)) &&
		FlagEOF.Match([]byte(*nextCursor))
	if EOF {
		nextCursor = nil
	}
	return resp, nextCursor, err
}

func (x *X) RequestGet(url string) (Response, error) {
	var response Response
	var err error
	client := resty.New()
	client.
		R().
		SetHeaders(
			StringMap{
				"authority":                 "twitter.com",
				"accept":                    "*/*",
				"accept-language":           "zh-CN,zh;q=0.9,en;q=0.8",
				"authorization":             "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
				"content-type":              "application/json",
				"cookie":                    x.cookies,
				"sec-ch-ua":                 `"Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v="116"`,
				"sec-ch-ua-mobile":          "?0",
				"sec-ch-ua-platform":        `"macOS"`,
				"sec-fetch-dest":            "empty",
				"sec-fetch-mode":            "cors",
				"sec-fetch-site":            "same-origin",
				"user-agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
				"x-client-transaction-id":   "b8qgMoxBUsxfsTLOsISkGvJ9/Atx8/2g/teNmizHQONJLNnUNCKMBxHt2eRlE6jkzOuM9G/tZ/mwb0KT9PAok5bnqm6vbg",
				"x-client-uuid":             "cc8bdbff-b377-4ffd-b53e-5767f6e50ba4",
				"x-csrf-token":              x.csrfToken,
				"x-twitter-active-user":     "yes",
				"x-twitter-auth-type":       "OAuth2Session",
				"x-twitter-client-language": "en",
			},
		).
		SetResult(&response).
		SetError(&err).
		Get(url)
	return response, err
}
