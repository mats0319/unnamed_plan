package db

import (
    "github.com/go-pg/pg/v10/orm"
    "github.com/mats9693/unnamed_plan/services/shared/db/model"
    "github.com/mats9693/unnamed_plan/services/shared/utils"
    "github.com/mats9693/utils/toy_server/db"
    "time"
)

type CloudFile model.CloudFile

var cloudFileIns = &CloudFile{}

func GetCloudFile() *CloudFile {
    return cloudFileIns
}

// QueryFirst 仅选择未删除的文件，pg.ErrNoRows
func (cf *CloudFile) QueryFirst(condition string, params ...interface{}) (file *model.CloudFile, err error) {
    file = &model.CloudFile{}

    err = mdb.WithNoTx(func(conn orm.DB) error {
        return conn.Model(file).Where(model.CloudFile_IsDeleted+" = ?", false).Where(condition, params...).First()
    })
    if err != nil {
        file = nil
    }

    return
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

func (cf *CloudFile) Insert(cloudFile *model.CloudFile) error {
    if len(cloudFile.ID) < 1 {
        cloudFile.Common = model.NewCommon()
    }

    return mdb.WithTx(func(conn orm.DB) error {
        _, err := conn.Model(cloudFile).Insert()
        return err
    })
}

// UpdateColumnsByFileID 仅可操作未删除的文件
func (cf *CloudFile) UpdateColumnsByFileID(data *model.CloudFile, columns ...string) error {
    data.UpdateTime = time.Duration(time.Now().Unix())

    return mdb.WithTx(func(conn orm.DB) error {
        query := conn.Model(data).Column(model.Common_UpdateTime)
        for i := range columns {
            query.Column(columns[i])
        }

        res, err := query.Where(model.CloudFile_IsDeleted+" = ?", false).
            Where(model.CloudFile_FileID + " = ?" + model.CloudFile_FileID).Update()
        if err != nil {
            return err
        }

        if res.RowsAffected() < 0 {
            return utils.NewError(utils.Error_FileAlreadyDeleted)
        }

        return nil
    })
}
