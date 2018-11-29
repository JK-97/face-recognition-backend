package exclude_record

import (
	"context"
	"github.com/google/uuid"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/util"
)

// AddExcludeRecord add a record to db
func AddExcludeRecord(p *schema.ExcludeRecordReq, excludeTime int64) error {
	dbp := schema.NewDBExcludeRecord(p, excludeTime)
    uid, _ :=  uuid.NewUUID()
    dbp.ID = uid.String()
	_, err := collection().InsertOne(context.Background(), dbp)
	if err != nil {
		return err
	}

	return nil
}

// UpdateRecord delete a record to db
func UpdateRecord(id string) error {
	selector := map[string]string{"_id": id}
	updater := map[string]map[string]int64{"$set": {"include_time": util.NowMilli()}}
	_, err := collection().UpdateOne(context.Background(), selector, updater)
	if err != nil {
		return err
	}

	return nil
}
