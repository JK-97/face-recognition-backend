package schema

// DBImages
type DBImages struct {
    ID          string `json:"id" bson:"_id, omitempty"`
    PID         string `json:"pid" bson:"pid, omitempty"`
    NationalID  string `json:"national_id" bson:"national_id"`
    Images      []string `json:"images" bson:"images"`
    Image       string `json:"image" bson:"image"`
    ImageID     string `json:"image_id" bson:"image_id"`
}

// DBImagesOnlyID
type DBImagesOnlyID struct {
    ID          string `json:"id" bson:"_id, omitempty"`
    PID         string `json:"pid" bson:"pid, omitempty"`
    NationalID  string `json:"national_id" bson:"national_id"`
    ImageID     string `json:"image_id" bson:"image_id"`
}

// DBImagesUpdater
type DBImagesUpdater struct {
    PID         string `json:"pid" bson:"pid, omitempty"`
    Images      []string `json:"images" bson:"images"`
    Image       string `json:"image" bson:"image"`
}
