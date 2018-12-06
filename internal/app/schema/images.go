package schema

// DBImages
type DBImages struct {
    ID          string `json:"id" bson:"_id, omitempty"`
    NationalID  string `json:"national_id" bson:"national_id"`
    Images      []string `json:"images" bson:"images"`
}
