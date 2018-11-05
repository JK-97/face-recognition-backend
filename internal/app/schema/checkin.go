package schema

// CheckinStatus is different status type
type CheckinStatus string

// Options of type CheckinStatus
const (
	CHECKING CheckinStatus = "checking"
	STOPPED  CheckinStatus = "stopped"
)

// CheckStatusResp is a response of CheckStatusGET
type CheckStatusResp struct {
	Status CheckinStatus `json:"status"`
}

// StopCheckinResp is a response of StopCheckinPOST
type StopCheckinResp struct {
	Timestamp int64 `json:"timestamp"`
}

// CheckinHistoryResp is a response of CheckinHistoryGET
type CheckinHistoryResp []int64

// CheckinResp is a response of CheckinGET
type CheckinResp struct {
	Timestamp     int64           `json:"timestamp"`
	CostTime      int64           `json:"cost_time"`
	ExpectedCount int             `json:"expected_count"`
	ActualCount   int             `json:"actual_count"`
	Detail        []CheckinPerson `json:"detail"`
}

// CheckinPerson is a person's info in a checkin
type CheckinPerson struct{ Person }
