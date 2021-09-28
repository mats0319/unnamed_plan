package dao

import (
	"time"

	"github.com/go-pg/pg/v10/orm"

	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	mdb "github.com/mats9693/utils/toy_server/db"
)

type CloudFile model.CloudFile

var cloudFileIns = &CloudFile{}

func GetCloudFile() *CloudFile {
	return cloudFileIns
}

func (cf *CloudFile) Insert(cloudFile *model.CloudFile) error {
	if len(cloudFile.ID) < 1 {
		cloudFile.Common = model.NewCommon()
	}

	return mdb.WithTx(func(conn orm.DB) error {
		_, err := conn.Model(cloudFile).Insert()
		return err
	})
}

// QueryPageByUploader 不查询已删除的记录，按照更新时间降序
func (cf *CloudFile) QueryPageByUploader(
	pageSize int,
	pageNum int,
	userID string,
) (files []*model.CloudFile, count int, err error) {
	err = mdb.WithNoTx(func(conn orm.DB) error {
		count, err = conn.Model(&files).Where(model.CloudFile_IsDeleted+" = ?", false).
			Where(model.CloudFile_UploadedBy+" = ?", userID).
			Order(model.Common_UpdateTime + " DESC").
			Offset((pageNum - 1) * pageSize).Limit(pageSize).SelectAndCount()

		return err
	})
	if err != nil {
		files = nil
		count = 0
	}

	return
}

// QueryPageInPublic 获取公开文件列表，要求上传者权限等级不高于指定用户（通过userID指定），分页，不查询已删除的记录，按照更新时间降序
/**
Core: sub-query
	select *
	from cloud_files cf
	where cf.uploaded_by in (
		select "user_id"
		from users u
		where "permission" <= (
			select "permission"
			from users u2
			where user_id = 'user id'
		)
	);
*/
func (cf *CloudFile) QueryPageInPublic(
	pageSize int,
	pageNum int,
	userID string,
) (files []*model.CloudFile, count int, err error) {
	err = mdb.WithNoTx(func(conn orm.DB) error {
		permission := conn.Model((*model.User)(nil)).Column(model.User_Permission).Where(model.User_UserID+" = ?", userID)
		userIDs := conn.Model((*model.User)(nil)).Column(model.User_UserID).Where(model.User_Permission+" <= (?)", permission)

		count, err = conn.Model(&files).Where(model.CloudFile_IsDeleted+" = ?", false).
			Where(model.CloudFile_IsPublic+" = ?", true).
			Where(model.CloudFile_UploadedBy+" in (?)", userIDs).
			Order(model.Common_UpdateTime + " DESC").
			Offset((pageNum - 1) * pageSize).Limit(pageSize).SelectAndCount()

		return err
	})
	if err != nil {
		files = nil
		count = 0
	}

	return
}

// QueryPageInDeleted 查询已删除的文件列表，要求上传者权限等级不高于指定用户（通过userID指定），分页，按照更新时间降序
// Core: sub-query, same as QueryPageInPublic
func (cf *CloudFile) QueryPageInDeleted(
	pageSize int,
	pageNum int,
	userID string,
) (files []*model.CloudFile, count int, err error) {
	err = mdb.WithNoTx(func(conn orm.DB) error {
		permission := conn.Model((*model.User)(nil)).Column(model.User_Permission).Where(model.User_UserID+" = ?", userID)
		userIDs := conn.Model((*model.User)(nil)).Column(model.User_UserID).Where(model.User_Permission+" <= (?)", permission)

		count, err = conn.Model(&files).Where(model.CloudFile_IsDeleted+" = ?", true).
			Where(model.CloudFile_UploadedBy+" in (?)", userIDs).
			Order(model.Common_UpdateTime + " DESC").
			Offset((pageNum - 1) * pageSize).Limit(pageSize).SelectAndCount()

		return err
	})
	if err != nil {
		files = nil
		count = 0
	}

	return
}

// UpdateColumnsByFileID 仅可操作未删除的文件
func (cf *CloudFile) UpdateColumnsByFileID(data *model.CloudFile, columns ...string) (err error) {
	data.UpdateTime = time.Duration(time.Now().Unix())

	return mdb.WithTx(func(conn orm.DB) error {
		query := conn.Model(data).Column(model.Common_UpdateTime)
		for i := range columns {
			query.Column(columns[i])
		}

		_, err = query.Where(model.CloudFile_IsDeleted+" = ?", false).
			Where(model.CloudFile_FileID + " = ?" + model.CloudFile_FileID).Update()

		return err
	})
}