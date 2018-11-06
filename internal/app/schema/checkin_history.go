package schema

// CheckinPeopleSet is people checked in a checkin
type CheckinPeopleSet map[string]struct{}

// CheckinHistory is a checkin record
type CheckinHistory struct {
	StartTime     int64    `bson:"start_time"`
	EndTime       int64    `bson:"end_time"`
	ExpectedCount int      `bson:"expected_count"`
	ActualCount   int      `bson:"actual_count"`
	Record        []string `bson:"record"`
}

// List converts CheckinPeopleSet to list
func (s *CheckinPeopleSet) List() []string {
	l := make([]string, 0, len(*s))
	for k := range *s {
		l = append(l, k)
	}
	return l
}
