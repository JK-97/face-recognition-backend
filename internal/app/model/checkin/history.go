package checkin

import (
	"sync"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/people"
)

type seal struct {
	startTime int64
	endTime   int64
}

// History is a checkin record
type History struct {
	StartTime     int64
	EndTime       int64
	ExpectedCount int
	ActualCount   int
	Record        checkRecordSet
}

var histories = map[int64]*History{}
var hisLock sync.RWMutex

func saveCheckin(s seal) {
	hisLock.Lock()
	defer hisLock.Unlock()

	histories[s.startTime] = &History{
		StartTime:     s.startTime,
		EndTime:       s.endTime,
		ExpectedCount: people.CountPeople(),
		ActualCount:   len(currentRecord),
		Record:        currentRecord,
	}
	currentRecord = checkRecordSet{}
}

// HistoryTimestamps returns all available timestamps for history query
func HistoryTimestamps() []int64 {
	hisLock.RLock()
	defer hisLock.RUnlock()

	timestamps := make([]int64, 0, len(histories))
	for key := range histories {
		timestamps = append(timestamps, key)
	}
	return timestamps
}

// GetHistory returns a checkin record with timestamp
func GetHistory(timestamp int64) *History {
	hisLock.RLock()
	defer hisLock.RUnlock()

	return histories[timestamp]
}
