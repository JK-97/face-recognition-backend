package route

import (
	"net/http"

	"gitlab.jiangxingai.com/luyor/tf-fence-backend/internal/app/controller"
)

// Routes adds routes to http
func Routes() {
	http.HandleFunc("/ping", controller.PingGet)
	http.HandleFunc("/img/", controller.ImgGET)
	http.HandleFunc("/fence_pos", controller.FenceGetPost)
	http.HandleFunc("/setting", controller.SettingGETPOST)

	http.HandleFunc("/", controller.IndexGET)
}
