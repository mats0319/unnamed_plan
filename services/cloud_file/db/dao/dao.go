package dao

import "github.com/mats9693/unnamed_plan/services/shared/db/model"

type CloudFileDao interface {
    // Insert insert one file
    Insert(*model.CloudFile) error

    // QueryOne query one file by id
    QueryOne(cloudFileID string) (*model.CloudFile, error)

    // QueryPageByUploader query files upload by designated user(designate by user id)
    // result not contains files marked as 'deleted', result order by 'update time' desc
    QueryPageByUploader(pageSize int, pageNum int, userID string) (files []*model.CloudFile, count int, err error)

    // QueryPageInPublic query files satisfy following requirements:
    //   1. file is 'public'
    //   2. permission of uploader is less than or equal to designated user(designate by user id)
    // result not contains files marked as 'deleted', result order by 'update time' desc
    QueryPageInPublic(pageSize int, pageNum int, userID string) (files []*model.CloudFile, count int, err error)

    // UpdateColumnsByFileID update designated 'columns' on designated undeleted file(designate by file id)
    UpdateColumnsByFileID(cloudFile *model.CloudFile, columns ...string) error
}
