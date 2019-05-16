package schema

// CheckinPeopleSet is people checked in a checkin
type CheckinPeopleSet map[string]int

// DeviceCheckinPeopleSet is people checked in a checkin at special device
type DeviceCheckinPeopleSet map[string]CheckinPeopleSet

// CheckinHistoryRecord is a checkin record
type CheckinHistoryRecord struct {
	Timestamp     int64    `bson:"timestamp"`
	PersonIDS     []string `bson:"person_ids"`
    CameraID      string   `bson:"camera_id"`
}

// DBCheckinHistoryRecord is a checkin record in db
type DBCheckinHistoryRecord struct {
	ID            string `json:"id" bson:"_id"`
	Timestamp     int64    `bson:"timestamp"`
	PersonIDS     []string `bson:"person_ids"`
    CameraID      string   `bson:"camera_id"`
}

// NewDBCheckinHistoryRecord creates DBCheckinHistoryRecord with CheckinHistoryRecord
func NewDBCheckinHistoryRecord(r *CheckinHistoryRecord, id string) *DBCheckinHistoryRecord {
    return &DBCheckinHistoryRecord{
        ID:         id,
        Timestamp:  r.Timestamp,
        PersonIDS:  r.PersonIDS,
        CameraID:   r.CameraID,
    }
}

// CheckinHistoryRecord gets checkin records info from DBCheckinHistoryRecord
func (p *DBCheckinHistoryRecord) CheckinHistoryRecord() CheckinHistoryRecord {
	return CheckinHistoryRecord{
        Timestamp:  p.Timestamp,
        PersonIDS:  p.PersonIDS,
        CameraID:   p.CameraID,
	}
}


