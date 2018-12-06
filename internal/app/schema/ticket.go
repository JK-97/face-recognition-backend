package schema

// UserTicket is user login ticket in db
type UserTicket struct {
    ID         string `json:"id" bson:"_id"`
    UserName   string `json:"username" bson:"username"`
    NonceStr   string `json:"nonce_str" bson:"nonce_str"`
    Timestamp  int64  `json:"timestamp" bson:"timestamp"`
}
