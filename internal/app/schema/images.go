package schema

// DBImages
type DBImages struct {
    ID          string `json:"id" bson:"_id, omitempty"`
    PID         string `json:"pid" bson:"pid, omitempty"`
    NationalID  string `json:"national_id" bson:"national_id"`
    Images      []string `json:"images" bson:"images"`
}

// DBImagesUpdater
type DBImagesUpdater struct {
    PID         string `json:"pid" bson:"pid, omitempty"`
    Images      []string `json:"images" bson:"images"`
}
