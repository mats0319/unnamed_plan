package backup

import (
	"log"
	"log/slog"
	"os"
	"testing"

	mdb "github.com/mats0319/unnamed_plan/server/internal/db"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

// 测试数据库备份/恢复功能，如果想要更清楚的分辨一条数据库记录是预设的还是恢复的，可以按照如下步骤操作：
// 1. 取消注释'DoBackupRecover'接口中的测试方法'DoSomeChangeForTest'
// 2. 取消注释'Recover'函数中的测试代码
// 这样再执行本程序，就可以区分一条数据是预设的还是恢复的了。
// 我们做了什么：恢复功能执行到写数据库之前，我们调用了修改方法修改数据（具体修改内容见接口方法的实现）
func TestBackupRecover(t *testing.T) {
	mlog.InitializeTest()
	defer mlog.Close()

	initDB()

	brm := NewBRManager(&UserBR{}, &NoteBR{})

	// test backup
	brm.Backup()

	mlog.Info("> Backup done.")

	// prepare recover data
	err := prepareRecoverData()
	if err != nil {
		os.Exit(1)
	}

	// test recover
	brm.Recover()

	mlog.Info("> Recover done.")
}

func prepareRecoverData() error {
	err := os.RemoveAll("./recover/")
	if err != nil {
		mlog.Error("remove dir failed", slog.Any("error", err))
		return err
	}

	err = os.Rename("./backup/", "./recover/")
	if err != nil {
		mlog.Error("rename folder failed", slog.Any("error", err))
		return err
	}

	return nil
}

func initDB() {
	db := mdb.InitTestDB()

	err := db.Migrator().DropTable(model.ModelList...)
	if err != nil {
		log.Fatalln("drop db table failed, error: ", err)
	}

	err = db.Migrator().CreateTable(model.ModelList...)
	if err != nil {
		log.Fatalln("create db table failed, error: ", err)
	}

	// preset data
	db.Create(presetUser)
	db.Create(presetNote)
}

var presetUser = []*model.User{
	{
		UserName: "user 1",
		Nickname: "user 1",
		Password: "user 1 pwd",
		IsAdmin:  true,
	},
	{
		UserName:  "user 2",
		Nickname:  "user 2",
		Password:  "user 2 pwd",
		EnableMFA: true,
		TOTPKey:   "NVQXE2LP",
	},
}

var presetNote = []*model.Note{
	{
		Writer:      "writer 1",
		WriterName:  "writer 1 name",
		IsAnonymous: false,
		Title:       "title 1",
		Content:     "content 1",
	},
	{
		Writer:      "writer 2",
		WriterName:  "writer 2 name",
		IsAnonymous: true,
		Title:       "title 2",
		Content:     "content 2",
	},
}
