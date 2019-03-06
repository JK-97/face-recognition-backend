package schema

// Camera is a post /detect request to face ai service
type Camera struct {
	ID         string `json:"id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	Rtmp       string `json:"rtmp" bson:"rtmp"`
	DeviceName string `json:"device_name" bson:"device_name"`
}

// CaptureCamera POST all device to video_capture
type CaptureCamera struct {
	Device string `json:"device"`
	Rtmp string `json:"rtmp"`
}

type CamerasRespItem struct {
    CamareURL   string `json:"camera_url"`
    HostIP      string `json:"host_ip"`
    ID          string `json:"id"`
    Name        string `json:"name"`
    StreamAddr  string `json:"stream_addr"`
    Type        string `json:"type"`
}

// CamerasResp revices all device from apigw
type CamerasResp struct {
    Data        []CamerasRespItem `json:"data"`
    Desc        string  `json:"desc"`
}

type OpenCameraReq struct {
    Action      string  `json:"action"`
    Timeout     int64   `json:"timeout"`
}

type OpenCameraResp struct {
    Data        string `json:"data"`
    Desc        string `json:"desc"`
}

// CameraListResp is response of get camera list
type CameraListResp []Camera
