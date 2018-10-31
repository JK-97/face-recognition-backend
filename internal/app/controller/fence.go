package controller

import (
	"encoding/json"
	"net/http"

	"gitlab.jiangxingai.com/luyor/tf-fence-backend/internal/app/model"
)

// FenceGetPost handles Get, Post methods.
// Get: get fence position
// Post: Set fence position
func FenceGetPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fence := model.GetFence()
		respondJSON(w, r, fence)

	case http.MethodPost:
		var fence FencePos
		if r.Body == nil {
			Error500(w, r)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&fence)
		if err != nil {
			Error500(w, r)
			return
		}
		model.SetFence([4]float32{fence.ymin, fence.xmin, fence.ymax, fence.xmax})
	default:
		Error404(w, r)
	}
}

// FencePos is a fence position
type FencePos struct {
	ymin, xmin, ymax, xmax float32
}
