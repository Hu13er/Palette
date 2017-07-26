package prof

/***********************************\
|* API:							   *|
\***********************************/

import (
	"io"
	"time"

	"gitlab.com/NagByte/Palette/helper"
)

var (
	ErrNotFound = io.EOF
)

func (ps *profService) GetProfile(username string) (profile, error) {
	query := ps.db.GetQuery("getProfile")
	result, err := ps.db.QueryOne(query, map[string]interface{}{"username": username})
	switch err {
	case nil:
	case ErrNotFound:
		return profile{}, ErrNotFound
	default:
		return profile{}, err
	}

	mp, _ := result[0].(map[string]interface{})

	return profile{
		Username:   username,
		FullName:   helper.SafeMap(mp, "fullName", "").(string),
		Bio:        helper.SafeMap(mp, "bio", "").(string),
		Location:   helper.SafeMap(mp, "location", "").(string),
		Follows:    helper.SafeMap(mp, "follows", int64(0)).(int64),
		FollowedBy: helper.SafeMap(mp, "followedBy", int64(0)).(int64),
		Wallpaper:  ps.fs.LargeDownloadURL(helper.SafeMap(mp, "wallpaper", "").(string)),
		Avatars: SmallLarge{
			Small: ps.fs.SmallDownloadURL(helper.SafeMap(mp, "avatar", "").(string)),
			Large: ps.fs.LargeDownloadURL(helper.SafeMap(mp, "avatar", "").(string)),
		},
	}, nil
}

func (ps *profService) IsFollowedBy(username1, username2 string) (bool, error) {
	query := ps.db.GetQuery("isFollowedBy")
	result, err := ps.db.QueryOne(query, map[string]interface{}{"username1": username1, "username2": username2})
	if err != nil {
		return false, err
	}

	return result[0].(bool), nil
}

func (ps *profService) UpdateProfile(username, fullName, bio, location, avatar, wallpaper string) error {
	query := ps.db.GetQuery("updateProfile")
	change := map[string]interface{}{}

	if fullName != "" {
		change["fullName"] = fullName
	}

	if bio != "" {
		change["bio"] = bio
	}

	if location != "" {
		change["location"] = location
	}

	if avatar != "" {
		change["avatar"] = avatar
	}

	if wallpaper != "" {
		change["wallpaper"] = wallpaper
	}

	_, err := ps.db.QueryOne(query, map[string]interface{}{"username": username, "change": change})
	switch err {
	case ErrNotFound:
		return ErrNotFound
	default:
		return err
	}
}

func (ps *profService) UpdateWallpaper(username string, reader io.Reader) error {
	return nil
}

func (ps *profService) UpdateAvatar(username string, reader io.Reader) error {
	return nil
}

func (ps *profService) Follow(username1, username2 string) error {
	query := ps.db.GetQuery("follow")
	_, err := ps.db.QueryOne(query, map[string]interface{}{"username1": username1, "username2": username2})
	switch err {
	case ErrNotFound:
		return ErrNotFound
	default:
		return err
	}
}

func (ps *profService) Unfollow(username1, username2 string) error {
	query := ps.db.GetQuery("unfollow")
	_, err := ps.db.QueryOne(query, map[string]interface{}{"username1": username1, "username2": username2})
	switch err {
	case ErrNotFound:
		return ErrNotFound
	default:
		return err
	}
}

func (ps *profService) Post(username, source, title, desc string, tags []string) error {
	query := ps.db.GetQuery("post")
	artID := helper.DefaultCharset.RandomStr(30)

	_, err := ps.db.QueryOne(query, map[string]interface{}{
		"username":      username,
		"artID":         artID,
		"title":         title,
		"desc":          desc,
		"tags":          tags,
		"displaySource": source,
	})
	switch err {
	case ErrNotFound:
		return ErrNotFound
	default:
		return err
	}
}

func (ps *profService) Like(username, artID string) error {
	query := ps.db.GetQuery("like")

	_, err := ps.db.QueryOne(query, map[string]interface{}{
		"username": username,
		"artID":    artID,
	})
	return err
}

func (ps *profService) Dislike(username, artID string) error {
	query := ps.db.GetQuery("dislike")

	_, err := ps.db.QueryOne(query, map[string]interface{}{
		"username": username,
		"artID":    artID,
	})
	return err
}

func (ps *profService) ensureArtID(artID string) bool {
	return true
}

func (ps *profService) GetLikes(artID string, count int, cursur int64) (usernames []string, nextCursur int64, hasNextPage bool, err error) {

	if !ps.ensureArtID(artID) {
		return nil, 0, false, ErrNotFound
	}

	if cursur <= 0 {
		cursur = time.Now().UnixNano()
	}

	query := ps.db.GetQuery("getLikes")

	result, err := ps.db.QueryAll(query, map[string]interface{}{
		"artID":  artID,
		"count":  count + 1,
		"cursur": cursur,
	})

	switch err {
	case nil:

		var lastCursur int64

		for _, v := range result {
			lastCursur = v[1].(int64)
			username := v[0].(string)
			usernames = append(usernames, username)
		}

		nextCursur = lastCursur - 1
		hasNextPage = len(usernames) > count
		if !hasNextPage {
			nextCursur = -1
		}

		return usernames, nextCursur, hasNextPage, nil
	default:
		return nil, -1, false, err
	}
}

func (ps *profService) GetPosts(username string, count int, cursur int64) (posts []post, nextCursur int64, hasNextPage bool, err error) {

	if ps.auth.IsUniqueUsername(username) {
		return nil, 0, false, ErrNotFound
	}

	if cursur <= 0 {
		cursur = time.Now().UnixNano()
	}
	query := ps.db.GetQuery("getPosts")

	result, err := ps.db.QueryAll(query, map[string]interface{}{
		"username": username,
		"count":    count + 1,
		"cursur":   cursur,
	})

	switch err {
	case nil:
		var lastCursur int64
		for _, v := range result {
			conv := v[0].(map[string]interface{})

			post := post{
				ArtID:         helper.SafeMap(conv, "artID", "").(string),
				Title:         helper.SafeMap(conv, "title", "").(string),
				Desc:          helper.SafeMap(conv, "desc", "").(string),
				CommentsCount: helper.SafeMap(conv, "comments_count", int64(0)).(int64),
				LikesCount:    helper.SafeMap(conv, "likes_count", int64(0)).(int64),
				Date:          helper.SafeMap(conv, "date", int64(0)).(int64),
				Tags: helper.ConvInterfaceSliceToStringSlice(
					conv["tags"].([]interface{})),
				DisplaySource: SmallLarge{
					Small: ps.fs.SmallDownloadURL(helper.SafeMap(conv, "displaySource", "").(string)),
					Large: ps.fs.LargeDownloadURL(helper.SafeMap(conv, "displaySource", "").(string)),
				},
			}
			lastCursur = post.Date
			posts = append(posts, post)
		}

		nextCursur = lastCursur - 1
		hasNextPage = len(posts) > count
		if !hasNextPage {
			nextCursur = -1
		}

		return posts, nextCursur, hasNextPage, nil
	default:
		return nil, nextCursur, false, err
	}
}

func (ps *profService) GetTimeline(username string, count int, cursur int64) (posts []post, nextCursur int64, hasNextPage bool, err error) {

	if ps.auth.IsUniqueUsername(username) {
		return nil, 0, false, ErrNotFound
	}

	if cursur <= 0 {
		cursur = time.Now().UnixNano()
	}
	query := ps.db.GetQuery("getTimeline")

	switch result, err := ps.db.QueryAll(query, map[string]interface{}{
		"username": username,
		"count":    count + 1,
		"cursur":   cursur,
	}); err {
	case nil:
		for _, v := range result {
			conv := v[0].(map[string]interface{})
			posts = append(posts, post{
				ArtID:         helper.SafeMap(conv, "artID", "").(string),
				Title:         helper.SafeMap(conv, "title", "").(string),
				Desc:          helper.SafeMap(conv, "desc", "").(string),
				CommentsCount: helper.SafeMap(conv, "comments_count", int64(0)).(int64),
				LikesCount:    helper.SafeMap(conv, "likes_count", int64(0)).(int64),
				Date:          helper.SafeMap(conv, "date", int64(0)).(int64),
				Tags: helper.ConvInterfaceSliceToStringSlice(
					helper.SafeMap(conv, "tags", []interface{}{}).([]interface{})),
				DisplaySource: SmallLarge{
					Small: ps.fs.SmallDownloadURL(helper.SafeMap(conv, "displaySource", "").(string)),
					Large: ps.fs.LargeDownloadURL(helper.SafeMap(conv, "displaySource", "").(string)),
				},
			})
		}

		nextCursur = 0
		if len(posts) > 0 {
			nextCursur = posts[len(posts)-1].Date - 1
		}

		return posts, nextCursur, len(posts) >= count, nil
	default:
		return nil, nextCursur, false, err
	}
}
