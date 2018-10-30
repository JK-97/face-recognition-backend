package route

import (
	"fmt"
	"net/http"

	"gitlab.jiangxingai.com/luyor/tf-fence-backend/internal/app/controller"
)

// Routes adds routes to http
func Routes() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})
	http.HandleFunc("/img/", controller.ImgGET)
	http.HandleFunc("/", controller.IndexGET)
}
