package backup

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"gorm.io/gorm/clause"
)

func Recover[T any](t doBackupRecover[T]) {
	dir := "./recover/" + t.Dir()
	entry, err := os.ReadDir(dir)
	if err != nil {
		mlog.Error(fmt.Sprintf("read dir '%s' failed, error: %v", dir, err))
		return
	}

	for i := range entry {
		if entry[i].IsDir() {
			continue // ignore folder
		}

		fileInfo, err := entry[i].Info()
		if err != nil {
			mlog.Error("get file info failed", slog.Any("error", err))
			continue
		}

		if !strings.HasSuffix(fileInfo.Name(), ".json") {
			continue // ignore not json files
		}

		err = recoverFile(dir+fileInfo.Name(), t)
		if err != nil {
			return
		}
	}
}

func recoverFile[T any](path string, t doBackupRecover[T]) error {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		mlog.Error("read file failed", slog.Any("error", err))
		return err
	}

	fileData := t.EmptySlice()
	err = json.Unmarshal(fileBytes, &fileData)
	if err != nil {
		mlog.Error("unmarshal file failed", slog.Any("error", err))
		return err
	}

	// 仅测试使用，为了方便看出一条数据库记录是不是通过该函数恢复的
	//t.DoSomeChangesForTest(fileData)

	clauseSkipAutoTime := clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(t.ColumnNames()), // all fields
	}

	err = dal.DB().Model(t.Model()).Clauses(clauseSkipAutoTime).Create(fileData).Error
	if err != nil {
		mlog.Error("save file failed", slog.Any("error", err))
		return err
	}

	return nil
}
