package controller

import (
	"encoding/json"
	"net/http"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
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
		var fence schema.FencePos
		if r.Body == nil {
			Error500(w, r)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&fence)
		if err != nil {
			Error500(w, r)
			return
		}
		model.SetFence(fence)
	default:
		Error404(w, r)
	}
}
