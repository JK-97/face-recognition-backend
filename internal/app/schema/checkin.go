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
	Timestamp     int64                 `json:"timestamp"`
	CostTime      int64                 `json:"cost_time"`
	ExpectedCount int                   `json:"expected_count"`
	ActualCount   int                   `json:"actual_count"`
	Present       []*CheckinPerson      `json:"present"`
	Absent        []*CheckinPerson      `json:"absent"`
    ExcludeRecord []*DBExcludeRecord    `json:"exclude_record"`
    Status        CheckinStatus         `json:"status"`
}

// CheckinPerson is a person's info in a checkin
type CheckinPerson struct {
	ID           string `json:"id"`
	SerialNumber string `json:"serial_number"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	NationalID   string `json:"national_id"`
	Image        string `json:"image"`
}

// CheckinPerson gets CheckinPerson from DBPerson
func (p *DBPerson) CheckinPerson() *CheckinPerson {
	return &CheckinPerson{
		ID:           p.ID,
		SerialNumber: p.SerialNumber,
		Name:         p.Name,
		Location:     p.Location,
		NationalID:   p.NationalID,
		Image:        p.Image,
	}
}
