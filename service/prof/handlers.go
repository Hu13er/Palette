package prof

/*
|* Handlers:
\***********************************/

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

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
	default:
		log.Println(he)
		log.Println(err)
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
		log.Println("BAZ")
		return
	}

	me := context.Get(r, "username").(string)
	switch isFollowed, err := ps.IsFollowedBy(me, he); err {
	case nil:
		response.FollowedByViewer = isFollowed
	default:
		log.Println(he)
		log.Println(err)
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
		log.Println("bar")
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
		log.Println(err)
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
	default:
		log.Println(err)
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}

}

func (ps *profService) followHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	me := context.Get(r, "username").(string)
	he := mux.Vars(r)["username"]

	switch err := ps.Follow(me, he); err {
	case nil:
	default:
		jsonEncoder.Encode(common.ErrorJSONResponse{ErrorDescription: "usernameDoesNotExists"})
	}
}

func (ps *profService) unfollowHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	me := context.Get(r, "username").(string)
	he := mux.Vars(r)["username"]

	switch err := ps.Unfollow(me, he); err {
	case nil:
	default:
		jsonEncoder.Encode(common.ErrorJSONResponse{ErrorDescription: "noUserWithThatUsernameInFollowings"})
	}
}

func (ps *profService) postHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	jsonDecoder := json.NewDecoder(r.Body)

	me := context.Get(r, "username").(string)

	var form PostForm
	if err := jsonDecoder.Decode(&form); err != nil {
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(common.StatusBadRequestError)
		return
	}

	switch err := ps.Post(me, form.Title, form.Desc, form.Tags); err {
	case nil:
	default:
		log.Println(err)
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
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
	default:
		log.Println(err)
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}

func (ps *profService) getTimelineHandler(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)
	jsonDecoder := json.NewDecoder(r.Body)

	me := context.Get(r, "username").(string)
	var form CursurForm
	if err := jsonDecoder.Decode(&form); err != nil {
		w.WriteHeader(common.StatusBadRequestError)
		jsonEncoder.Encode(common.ResponseBadRequest)
		return
	}

	switch posts, nextCursur, hasNextPage, err := ps.GetPosts(me, 20, form.Cursur); err {
	case nil:
		jsonEncoder.Encode(Pager{
			Elements:       posts,
			HasNextPage:    hasNextPage,
			NextPageCursur: nextCursur,
		})
	default:
		log.Println(err)
		w.WriteHeader(common.StatusInternalServerError)
		jsonEncoder.Encode(common.ResponseInternalServerError)
	}
}
