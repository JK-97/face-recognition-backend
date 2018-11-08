package route

import (
	"net/http"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/controller"
)

// Routes adds routes to http
func Routes() http.Handler {
	mux := http.NewServeMux()
	handler := logRequest(mux)

	mux.HandleFunc("/api/v1/ping", controller.PingGet)
	mux.HandleFunc("/api/v1/img/", controller.ImgGET)

	// checkin
	mux.HandleFunc("/api/v1/start_checkin", controller.StartCheckinPOST)
	mux.HandleFunc("/api/v1/check_status", controller.CheckStatusGET)
	mux.HandleFunc("/api/v1/stop_checkin", controller.StopCheckinPOST)

	// history
	mux.HandleFunc("/api/v1/checkin_history", controller.CheckinHistoryGET)
	mux.HandleFunc("/api/v1/checkin", controller.CheckinGET) // ?timestamp=xxxx

	// people
	mux.HandleFunc("/api/v1/face_records", controller.FaceRecordsGET)
	mux.HandleFunc("/api/v1/checkin_people", controller.CheckinPeoplePOSTDELETE)
	mux.HandleFunc("/api/v1/checkin_people_list", controller.CheckinPeopleListGET)
	mux.HandleFunc("/api/v1/start_recording", controller.StartRecordingPOST)

	mux.HandleFunc("/", controller.IndexGET)

	return handler
}
