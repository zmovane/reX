package reX

import (
	"golang.org/x/exp/slices"
)

func (x *X) GetFollowingsByScreenName(user string, cursor *string) (resp []UserResults, nextCursor *string, err error) {
	uid, _ := x.scraper.GetUserIDByScreenName(user)
	return x.GetRelationsById(uid, cursor, Following)
}

func (x *X) GetFollowingsById(uid string, cursor *string) (resp []UserResults, nextCursor *string, err error) {
	return x.GetRelationsById(uid, cursor, Following)
}

func (x *X) IsFollowing(uid string, uidOfFollower string) bool {
	var err error
	var cursor *string
	var followings []string
	for {
		var pagedUsers []UserResults
		pagedUsers, cursor, err = x.GetFollowingsById(uidOfFollower, cursor)
		ids := Map(pagedUsers, func(o UserResults) string { return o.Result.RESTID })
		followings = append(followings, ids...)
		if cursor == nil || err != nil {
			break
		}
	}
	return slices.Contains(followings, uid)
}