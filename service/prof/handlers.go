package prof

/*
|* Handlers:
\***********************************/

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
	"gitlab.com/NagByte/Palette/service/auth"
	"gitlab.com/NagByte/Palette/service/common"
)

func (ps *profService) getProfileHandler(w http.ResponseWriter, r *http.Request) {

	jsonEncoder := json.NewEncoder(w)
	he := mux.Vars(r)["username"]
	var response getProfileResponse

	switch prof, err := ps.GetProfile(he); err {
	case nil:
		response.profile = prof
	case ErrNotFound:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(responseUserNotFound)
		return
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
		return
	}

	me := auth.GetUsername(r)
	switch isFollowed, err := ps.IsFollowedBy(me, he); err {
	case nil:
		response.FollowedByViewer = isFollowed
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
		return
	}

	response.RequestByOwner = he == me

	switch posts, nextCursur, hasNextPage, err := ps.GetPosts(he, 20, -1); err {
	case nil:
		response.Pager = Pager{
			Elements:       posts,
			HasNextPage:    hasNextPage,
			NextPageCursur: nextCursur,
		}
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
		return
	}

	jsonEncoder.Encode(response)
}

func (ps *profService) updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	jsonDecoder := json.NewDecoder(r.Body)
	user := auth.GetUsername(r)

	form := UpdateProfileForm{}
	if err := jsonDecoder.Decode(&form); err != nil {
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}

	switch err := ps.UpdateProfile(user, form.FullName, form.Bio, form.Location); err {
	case nil:
	case ErrNotFound:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(responseUserNotFound)
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}

}

func (ps *profService) followHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	me := auth.GetUsername(r)
	he := mux.Vars(r)["username"]

	var do func(string, string) error
	switch r.Method {
	case "PUT":
		do = ps.Follow
	case "DELETE":
		do = ps.Unfollow
	default:
		return
	}

	switch err := do(me, he); err {
	case nil:
	case ErrNotFound:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(responseUserNotFound)
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}

func (ps *profService) postHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	jsonDecoder := json.NewDecoder(r.Body)

	me := auth.GetUsername(r)

	var form PostForm
	if err := jsonDecoder.Decode(&form); err != nil {
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(common.StatusBadRequestError)
		return
	}

	switch err := ps.Post(me, form.FileToken, form.Title, form.Desc, form.Tags); err {
	case nil:
	case ErrNotFound:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(responseUserNotFound)
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}

func (ps *profService) likePostHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)

	artID := mux.Vars(r)["artID"]
	me := auth.GetUsername(r)

	var do func(string, string) error
	switch r.Method {
	case "PUT":
		do = ps.Like
	case "DELETE":
		do = ps.Dislike
	default:
		return
	}

	switch err := do(me, artID); err {
	case nil:
	case ErrNotFound:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(responsePostNotFound)
	default:
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
		logrus.Errorln(err)
	}
}

func (ps *profService) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	jsonDecoder := json.NewDecoder(r.Body)

	he := mux.Vars(r)["username"]
	var form CursurForm
	if err := jsonDecoder.Decode(&form); err != nil {
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(common.ResponseBadRequest)
		return
	}

	switch posts, nextCursur, hasNextPage, err := ps.GetPosts(he, 20, form.Cursur); err {
	case nil:
		jsonEncoder.Encode(Pager{
			Elements:       posts,
			HasNextPage:    hasNextPage,
			NextPageCursur: nextCursur,
		})
	case ErrNotFound:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(responseUserNotFound)
	default:
		log.Println(err)
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}

func (ps *profService) getTimelineHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	jsonDecoder := json.NewDecoder(r.Body)

	me := auth.GetUsername(r)
	var form CursurForm
	if err := jsonDecoder.Decode(&form); err != nil {
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(common.ResponseBadRequest)
		return
	}

	switch posts, nextCursur, hasNextPage, err := ps.GetTimeline(me, 20, form.Cursur); err {
	case nil:
		jsonEncoder.Encode(Pager{
			Elements:       posts,
			HasNextPage:    hasNextPage,
			NextPageCursur: nextCursur,
		})
	case ErrNotFound:
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(responseUserNotFound)
	default:
		log.Println(err)
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}
