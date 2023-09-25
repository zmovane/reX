package reX

import "golang.org/x/exp/slices"

func (x *X) GetFollowingsByScreenName(user string, cursor *string) (resp []UserResults, nextCursor *string, err error) {
	uid, _ := x.scraper.GetUserIDByScreenName(user)
	return x.GetRelationsById(uid, cursor, Following)
}

func (x *X) GetFollowingsById(uid string, cursor *string) (resp []UserResults, nextCursor *string, err error) {
	return x.GetRelationsById(uid, cursor, Following)
}

func (x *X) IsFollowingByScreenName(screenname string, screennameOfFollower string) bool {
	var err error
	var cursor *string
	var followings []string
	for {
		var pagedUsers []UserResults
		pagedUsers, cursor, err = x.GetFollowingsByScreenName(screenname, cursor)
		names := Map(pagedUsers, func(o UserResults) string { return o.Result.Legacy.ScreenName })
		followings = append(followings, names...)
		if cursor == nil || err != nil {
			break
		}
	}
	return slices.Contains(followings, screennameOfFollower)
}

func (x *X) IsFollowingByUID(uid string, uidOfFollower string) bool {
	var err error
	var cursor *string
	var followings []string
	for {
		var pagedUsers []UserResults
		pagedUsers, cursor, err = x.GetFollowingsById(uid, cursor)
		ids := Map(pagedUsers, func(o UserResults) string { return o.Result.ID })
		followings = append(followings, ids...)
		if cursor == nil || err != nil {
			break
		}
	}
	return slices.Contains(followings, uidOfFollower)
}
