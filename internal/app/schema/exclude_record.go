package schema

// DBExcludeRecord is exclude_person in mongo
type DBExcludeRecord struct {
	ID          string     `json:"id" bson:"_id"`
	Peoples     []DBPerson `json:"peoples" bson:"peoples"`
	Reason      string     `json:"reason" bson:"reason"`
	ExcludeTime int64      `json:"exclude_time" bson:"exclude_time"`
	IncludeTime int64      `json:"include_time" bson:"include_time"`
}

// ExcludeRecordReq is body for POST ExcludeRecord
type ExcludeRecordReq struct {
	Peoples []DBPerson `json:"peoples"`
	Reason  string     `json:"reason"`
}

// NewDBExcludeRecord create exclude record
func NewDBExcludeRecord(p *ExcludeRecordReq, excludeTime int64) *DBExcludeRecord {
	return &DBExcludeRecord{
		Peoples:     p.Peoples,
		ExcludeTime: excludeTime,
		Reason:      p.Reason,
		IncludeTime: -1,
	}
}

// CheckinExcludeRecordListResp is response of GET ExcludeRecord
type CheckinExcludeRecordListResp []DBExcludeRecord
