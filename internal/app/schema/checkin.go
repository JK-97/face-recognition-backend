package schema

// CheckinStatus is different status type
type CheckinStatus string

// Options of type CheckinStatus
const (
	CHECKING CheckinStatus = "checking"
	STOPPED  CheckinStatus = "stopped"
)

// CheckinResp is a response of CheckinGET
type CheckinResp struct {
	StartTime     int64                 `json:"start_time"`
	EndTime       int64                 `json:"end_time"`
	Person        []*Person             `json:"person"`
}
