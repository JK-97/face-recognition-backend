package schema

// Camera is a post /detect request to face ai service
type Camera struct {
	ID         string `json:"id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	Rtmp       string `json:"rtmp" bson:"rtmp"`
	DeviceName string `json:"device_name" bson:"device_name"`
}

// CameraListResp is response of get camera list
type CameraListResp []Camera
