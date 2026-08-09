package backup

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

func Backup[T any](t doBackupRecover[T]) {
	dir := "./backup/" + t.Dir()
	err := os.MkdirAll(dir, 0644)
	if err != nil {
		mlog.Error("mkdir failed", slog.Any("error", err))
		return
	}

	// if it has data need backup
	var count int64
	err = dal.DB().Unscoped().Model(t.Model()).Where(t.Condition()).Count(&count).Error
	if err != nil {
		mlog.Error("db count failed", slog.Any("error", err))
		return
	}
	if count < 1 { // no data need backup
		return
	}

	// do backup in page
	pageSize := 100
	timestamp := time.Now().UnixMilli()
	for range int(count)/pageSize + 1 {
		dbRecords := t.EmptySlice()
		err := dal.DB().Unscoped().Model(t.Model()).Where(t.Condition()).Limit(pageSize).Find(&dbRecords).Error
		if err != nil {
			mlog.Error("get data need to backup failed", slog.Any("error", err))
			return
		}

		for _, record := range dbRecords {
			// gen file path
			index := uuidToIndex(t.ID(record), 16) // hard code: max 16 files for each table
			filePath := fmt.Sprintf("%s%d.json", dir, index)

			// read data from file
			fileData := t.EmptySlice()
			fileBytes, err := os.ReadFile(filePath) // error if file not exist
			if err == nil {
				err = json.Unmarshal(fileBytes, &fileData)
				if err != nil {
					mlog.Error("unmarshal file failed", slog.Any("error", err))
					return
				}
			}

			// set 'backupAt' and update file data
			// 因为这里更改了查询条件涉及的列（备份时间），所以每次分页查询均查询第一页
			t.Update(record, timestamp)

			isExist := false
			for i := range fileData {
				if t.ID(fileData[i]) == t.ID(record) {
					fileData[i] = record
					isExist = true
					break
				}
			}
			if !isExist { // new data which is first time do backup
				fileData = append(fileData, record)
			}

			// write file and update db record
			// 检查：写文件成功但是写数据库失败，下一次会重新尝试备份，而备份函数具有幂等性，所以可以不写在一个事务里
			fileBytes, err = json.Marshal(fileData)
			if err != nil {
				mlog.Error("marshal file failed", slog.Any("error", err))
				return
			}

			err = os.WriteFile(filePath, fileBytes, 0644)
			if err != nil { // implicit create file at first time
				mlog.Error("write file failed", slog.Any("error", err))
				return
			}

			err = dal.DB().Unscoped().Model(record).UpdateColumns(record).Error
			if err != nil { // UpdateColumns skip hooks and auto-updateTime
				mlog.Error("update db data failed", slog.Any("error", err))
				return
			}
		}
	}
}

func uuidToIndex(id uuid.UUID, max int) int {
	var v int
	for i := range id {
		v += int(id[i])
	}

	return v % max
}
