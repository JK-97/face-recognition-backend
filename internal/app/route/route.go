package route

import (
	"fmt"
	"net/http"

	"gitlab.jiangxingai.com/luyor/tf-pose-backend/internal/app/controller"
)

// Routes adds routes to http
func Routes() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})
	http.HandleFunc("/", controller.IndexGET)
}
