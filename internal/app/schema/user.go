package schema

// User is user in db
type User struct {
	ID       string `json:"id" bson:"_id"`
	UserName string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
