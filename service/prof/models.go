package prof

type SmallLarge struct {
	Small string `json:"small"`
	Large string `json:"large"`
}

type profile struct {
	Username   string     `json:"username"`
	FullName   string     `json:"full_name"`
	Bio        string     `json:"bio"`
	FollowedBy int64      `json:"followed_by"`
	Follows    int64      `json:"follows"`
	Location   string     `json:"location"`
	Avatars    SmallLarge `json:"avatar"`
	Wallpaper  string     `json:"wallpaper"`
}

type post struct {
	ArtID         string     `json:"art_id"`
	Title         string     `json:"title"`
	Desc          string     `json:"desc"`
	LikesCount    int64      `json:"likes_count"`
	CommentsCount int64      `json:"comments_count"`
	Tags          []string   `json:"tags"`
	Date          int64      `json:"date"`
	DisplaySource SmallLarge `json:"display_source"`
}
