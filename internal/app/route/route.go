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
	mux.HandleFunc("/api/v1/checkin_people", controller.CheckinPeoplePOSTDELETEGETPUT) // ?id=xxx&full_images=1
	mux.HandleFunc("/api/v1/checkin_people_list", controller.CheckinPeopleListGET) // ?exclude=1
	mux.HandleFunc("/api/v1/start_recording", controller.StartRecordingPOST)

	// camera
	mux.HandleFunc("/api/v1/cameras", controller.CamerasGETPOST)
	// GET [{"name": "", "rtmp": "", "device_name": ""}]
	// POST {"name": "rtmp": ""}

	// exclude people
	mux.HandleFunc("/api/v1/exclude_record", controller.ExcludeRecordGETPOSTPUT)
	// POST {"people": [{}, {}"], reason: ""}
	// GET response: [{"people": [{}, {}}], exclude_time: "", reason: "", id: "", include_time: }, ] query args: all=1&skip=0&limit=10
	// PUT query args id=

	// get image
	mux.HandleFunc("/api/v1/checkin_people_image", controller.CheckinPeopleImageGET) // ?id=
	// GET response [{}]

	// login
	mux.HandleFunc("/api/v1/login", controller.LoginPOST)
	// POST {"username": "", "password": "md5passwd"}

    // settings 
	mux.HandleFunc("/api/v1/settings", controller.SettingsPOSTGET)
    // POST {"starttime": {"hour":, "minute":, "second"}, "endtime": {"hour":, "minute": "second"}, "interval": seconds}
    // GET {"starttime": {"hour":, "minute":, "second"}, "endtime": {"hour":, "minute": "second"}, "interval": seconds}

	mux.HandleFunc("/", controller.IndexGET)

	return handler
}
