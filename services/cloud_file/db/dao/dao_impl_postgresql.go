package dao

import (
    "fmt"
    "github.com/mats9693/unnamed_plan/services/shared/db/model"
    "github.com/mats9693/unnamed_plan/services/shared/utils"
    "github.com/mats9693/utils/toy_server/db"
    "time"
)

type CloudFilePostgresql model.CloudFile

var _ CloudFileDao = (*CloudFilePostgresql)(nil)

func (cf *CloudFilePostgresql) Insert(cloudFile *model.CloudFile) error {
    if len(cloudFile.ID) < 1 {
        cloudFile.Common = model.NewCommon()
    }

    return mdb.DB().WithTx(func(conn mdb.Conn) error {
        _, err := conn.PostgresqlConn.Model(cloudFile).Insert()
        return err
    })
}

func (cf *CloudFilePostgresql) QueryOne(cloudFileID string) (file *model.CloudFile, err error) {
    file = &model.CloudFile{}

    condition := fmt.Sprintf("%s = ? and %s = ?", model.CloudFile_IsDeleted, model.CloudFile_FileID)

    err = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
        return conn.PostgresqlConn.Model(file).Where(condition, false, cloudFileID).First()
    })
    if err != nil {
        file = nil
    }

    return
}

func (cf *CloudFilePostgresql) QueryPageByUploader(
    pageSize int,
    pageNum int,
    userID string,
) (files []*model.CloudFile, count int, err error) {
    err = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
        count, err = conn.PostgresqlConn.Model(&files).
            Where(model.CloudFile_IsDeleted+" = ?", false).
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

// QueryPageInPublic
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
func (cf *CloudFilePostgresql) QueryPageInPublic(
    pageSize int,
    pageNum int,
    userID string,
) (files []*model.CloudFile, count int, err error) {
    err = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
        permission := conn.PostgresqlConn.Model((*model.User)(nil)).Column(model.User_Permission).Where(model.User_UserID+" = ?", userID)
        userIDs := conn.PostgresqlConn.Model((*model.User)(nil)).Column(model.User_UserID).Where(model.User_Permission+" <= (?)", permission)

        count, err = conn.PostgresqlConn.Model(&files).
            Where(model.CloudFile_IsDeleted+" = ?", false).
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

func (cf *CloudFilePostgresql) UpdateColumnsByFileID(data *model.CloudFile, columns ...string) error {
    data.UpdateTime = time.Duration(time.Now().Unix())

    return mdb.DB().WithTx(func(conn mdb.Conn) error {
        query := conn.PostgresqlConn.Model(data).Column(model.Common_UpdateTime)
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
