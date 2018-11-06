package schema

// DetectFaceReq is a post /detect request to face ai service
type DetectFaceReq struct {
	Image string `json:"image"` // jpeg (base64 encoded)
}

// DetectFaceResp is a response from post /detect request to face ai service
type DetectFaceResp struct {
	Data struct {
		Recognitions []Recognition `json:"recognitions"`
	} `json:"data"`
}

// Recognition is a face recognition result
type Recognition struct {
	ID       string `json:"id"`
	Position struct {
		Ymin int
		Xmin int
		Ymax int
		Xmax int
	} `json:"position"`
}

// RecordReq is a post /record request to face ai service
type RecordReq struct {
	ID     string   `json:"id"`
	Images []string `json:"images"`
}
