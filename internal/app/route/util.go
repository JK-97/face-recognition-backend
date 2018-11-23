package route

import (
	"net/http"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/controller"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

		if r.URL.Path != "/api/v1/login" {
			// check login cookie
            if controller.CheckLoginSession(w, r) != nil {
                return
            }
		}

		handler.ServeHTTP(w, r)
	})
}
