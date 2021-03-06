package prof

import (
	"gitlab.com/NagByte/Palette/service/common"
)

var (
	responseUserNotFound = common.ErrorJSONResponse{ErrorDescription: "userNotFound"}
	responsePostNotFound = common.ErrorJSONResponse{ErrorDescription: "postNotFound"}
)

type Pager struct {
	HasNextPage    bool        `joson:"has_next_page"`
	Elements       interface{} `json:"elements"`
	NextPageCursur int64       `json:"next_page_cursur"`
}

type UpdateProfileForm struct {
	FullName  string `json:"full_name"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Wallpaper string `json:"wallpaper"`
	Avatar    string `json:"avatar"`
}

type getProfileResponse struct {
	profile
	FollowedByViewer bool  `json:"followed_by_viewer"`
	RequestByOwner   bool  `json:"request_by_owner"`
	Pager            Pager `json:"pager"`
}

type PostForm struct {
	Title     string   `json:"title"`
	Desc      string   `json:"desc"`
	Tags      []string `json:"tags"`
	FileToken string   `json:"fileToken"`
}

type CursurForm struct {
	Cursur int64 `json:"cursur"`
}

type ArtTokenForm struct {
	ArtID string `json:"artID"`
}
