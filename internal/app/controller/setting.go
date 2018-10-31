package controller

import (
	"net/http"
	"strconv"

	"gitlab.jiangxingai.com/luyor/tf-fence-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/tf-fence-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/tf-fence-backend/log"
)

// SettingGETPOST handles Get, Post methods.
// GET: display setting page.
// POST: set fence position
func SettingGETPOST(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fence := model.GetFence()
		err := serveTemplate("web/template/setting.tmpl", fence, w, r)
		if err != nil {
			log.Info(err)
		}

	case http.MethodPost:
		r.ParseForm()
		Ymin, _ := strconv.ParseFloat(r.FormValue("Ymin"), 32)
		Xmin, _ := strconv.ParseFloat(r.FormValue("Xmin"), 32)
		Ymax, _ := strconv.ParseFloat(r.FormValue("Ymax"), 32)
		Xmax, _ := strconv.ParseFloat(r.FormValue("Xmax"), 32)
		fence := schema.FencePos{
			Ymin: float32(Ymin),
			Xmin: float32(Xmin),
			Ymax: float32(Ymax),
			Xmax: float32(Xmax),
		}
		model.SetFence(fence)
		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		Error404(w, r)
	}
}
