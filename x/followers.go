package reX

func (x *X) GetFollowersByScreenName(user string, cursor *string) (resp []UserResults, nextCursor *string, err error) {
	uid, _ := x.scraper.GetUserIDByScreenName(user)
	return x.GetRelationsById(uid, cursor, Follower)
}

func (x *X) GetFollowersById(uid string, cursor *string) (resp []UserResults, nextCursor *string, err error) {
	return x.GetRelationsById(uid, cursor, Follower)
}
