package prof

type Pager struct {
	HasNextPage    bool        `joson:"has_next_page"`
	Elements       interface{} `json:"elements"`
	NextPageCursur int64       `json:"next_page_cursur"`
}

type UpdateProfileForm struct {
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
	Location string `json:"location"`
}

type getProfileResponse struct {
	profile
	FollowedByViewer bool  `json:"followed_by_viewer"`
	RequestByOwner   bool  `json:"request_by_owner"`
	Pager            Pager `json:"pager"`
}

type PostForm struct {
	Title string   `json:"title"`
	Desc  string   `json:"desc"`
	Tags  []string `json:"tags"`
}

type CursurForm struct {
	Cursur int64 `json:"cursur"`
}
