package handlers

import (
	"github.com/mats9693/unnamed_plan/admin_data/config"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/http/structure_defination"
	"github.com/mats9693/unnamed_plan/admin_data/utils"
	"github.com/mats9693/utils/toy_server/http"
	"github.com/mats9693/utils/toy_server/utils"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

const backupFileSuffix = ".backup"

func ListCloudFileByUploader(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListCloudFileByUploaderReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || params.PageSize < 1 || params.PageNum < 1 {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			Int("page size", params.PageSize),
			Int("page num", params.PageNum))
	}

	files, count, err := dao.GetCloudFile().QueryPageByUploader(params.PageSize, params.PageNum, params.OperatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	fileListRes := make([]*structure.FileRes, 0, len(files))
	for i := range files {
		fileListRes = append(fileListRes, fileDBToHTTPRes(files[i]))
	}

	return mhttp.Response(structure.MakeListCloudFileByUploaderRes(count, fileListRes))
}

func ListPublicCloudFile(r *http.Request) *mhttp.ResponseData {
	params := &structure.ListPublicCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || params.PageSize < 1 || params.PageNum < 1 {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			Int("page size", params.PageSize),
			Int("page num", params.PageNum))
	}

	files, count, err := dao.GetCloudFile().QueryPageInPublic(params.PageSize, params.PageNum, params.OperatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	fileListRes := make([]*structure.FileRes, 0, len(files))
	for i := range files {
		fileListRes = append(fileListRes, fileDBToHTTPRes(files[i]))
	}

	return mhttp.Response(structure.MakeListPublicCloudFileRes(count, fileListRes))
}

func UploadCloudFile(r *http.Request) *mhttp.ResponseData {
	params := &structure.UploadCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}
	defer func() {
		if params.File != nil {
			_ = params.File.Close()
		}
	}()

	if len(params.OperatorID) < 1 || len(params.FileName) < 1 || len(params.ExtensionName) < 1 || params.LastModifiedTime < 1 {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			String("file name", params.FileName),
			String("extension name", params.ExtensionName),
			Int("last modified time", params.LastModifiedTime))
	}

	// make sure target directory structure exist
	dir := spliceFileDir(params.IsPublic, params.OperatorID)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	// save file
	fileID := utils.CalcSHA256(params.OperatorID + time.Now().GoString())
	absolutePath := spliceFilePath(dir, fileID, params.ExtensionName)

	err = saveFile(params.File, absolutePath, params.FileHeader.Size)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	// save db, if failed, remove file
	err = dao.GetCloudFile().Insert(&model.CloudFile{
		UploadedBy:       params.OperatorID,
		FileID:           fileID,
		FileName:         params.FileName,
		ExtensionName:    params.ExtensionName,
		LastModifiedTime: time.Duration(params.LastModifiedTime),
		FileSize:         params.FileHeader.Size,
		IsPublic:         params.IsPublic,
	})
	if err != nil {
		return mhttp.ResponseWithError(utils.ErrorsToString(err, os.Remove(absolutePath)))
	}

	return mhttp.Response(structure.MakeUploadCloudFileRes(true))
}

func ModifyCloudFile(r *http.Request) *mhttp.ResponseData {
	params := &structure.ModifyCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}
	defer func() {
		if params.File != nil {
			_ = params.File.Close()
		}
	}()

	if len(params.OperatorID) < 1 || len(params.FileID) < 1 || len(params.Password) < 1 {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			String("file id", params.FileID),
			String("password", params.Password))
	}

	_, err := verifyPwdByUserID(params.Password, params.OperatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	fileRecord, err := dao.GetCloudFile().QueryFirst(model.CloudFile_FileID+" = ?", params.FileID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	if len(params.FileName)+len(params.ExtensionName) < 1 && fileRecord.IsPublic == params.IsPublic && params.File == nil {
		return mhttp.ResponseWithError(error_NoValidModification)
	}

	updateColumns := make([]string, 0, 5)

	// update file, backup old file in private folder
	if params.File != nil {
		// make sure private dir exist, only when old file is public
		if fileRecord.IsPublic {
			err = os.MkdirAll(spliceFileDir(false, params.OperatorID), 0755)
			if err != nil {
				return mhttp.ResponseWithError(err.Error())
			}
		}

		// backup old file, may implicitly include move file
		oldPath := spliceFilePath(spliceFileDir(fileRecord.IsPublic, params.OperatorID), fileRecord.FileID, params.ExtensionName)
		privPath := spliceFilePath(spliceFileDir(false, params.OperatorID), fileRecord.FileID, params.ExtensionName)

		err = os.Rename(oldPath, getValidBackupFileName(privPath))
		if err != nil {
			return mhttp.ResponseWithError(err.Error())
		}

		// save new file
		absolutePath := spliceFilePath(spliceFileDir(params.IsPublic, params.OperatorID), fileRecord.FileID, params.ExtensionName)
		err = saveFile(params.File, absolutePath, params.FileHeader.Size)
		if err != nil {
			return mhttp.ResponseWithError(err.Error())
		}

		// update db record: file last modified time and file size
		fileRecord.LastModifiedTime = time.Duration(params.LastModifiedTime)
		updateColumns = append(updateColumns, model.CloudFile_LastModifiedTime)

		fileRecord.FileSize = params.FileHeader.Size
		updateColumns = append(updateColumns, model.CloudFile_FileSize)
	}

	if len(params.FileName) > 0 {
		fileRecord.FileName = params.FileName
		updateColumns = append(updateColumns, model.CloudFile_FileName)
	}
	if len(params.ExtensionName) > 0 && fileRecord.ExtensionName != params.ExtensionName {
		fileRecord.ExtensionName = params.ExtensionName
		updateColumns = append(updateColumns, model.CloudFile_ExtensionName)
	}
	if fileRecord.IsPublic != params.IsPublic {
		fileRecord.IsPublic = params.IsPublic
		updateColumns = append(updateColumns, model.CloudFile_IsPublic)
	}

	// update db record, if update failed and have new file, remove file
	err = dao.GetCloudFile().UpdateColumnsByFileID(fileRecord, updateColumns...)
	if err != nil {
		errMsg := err.Error()
		if params.File != nil {
			dir := spliceFileDir(params.IsPublic, params.OperatorID)
			absolutePath := spliceFilePath(dir, fileRecord.FileID, params.ExtensionName)
			errMsg += utils.ErrorsToString(os.Remove(absolutePath))
		}
		return mhttp.ResponseWithError(errMsg)
	}

	return mhttp.Response(structure.MakeModifyCloudFileRes(true))
}

func DeleteCloudFile(r *http.Request) *mhttp.ResponseData {
	params := &structure.DeleteCloudFileReqParams{}
	if errMsg := params.Decode(r); len(errMsg) > 0 {
		return mhttp.ResponseWithError(errMsg)
	}

	if len(params.OperatorID) < 1 || len(params.Password) < 1 || len(params.FileID) < 1 {
		return mhttp.ResponseWithError(error_InvalidParams,
			String("operator id", params.OperatorID),
			String("password", params.Password),
			String("file id", params.FileID))
	}

	_, err := verifyPwdByUserID(params.Password, params.OperatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	err = dao.GetCloudFile().UpdateColumnsByFileID(&model.CloudFile{
		FileID:    params.FileID,
		IsDeleted: true,
	}, model.CloudFile_IsDeleted)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	return mhttp.Response(structure.MakeDeleteCloudFileRes(true))
}

func saveFile(file multipart.File, absolutePath string, fileSize int64) (err error) {
	fileContent := make([]byte, fileSize) // require enough length before read
	_, err = file.Read(fileContent)
	if err != nil {
		return
	}

	err = os.WriteFile(absolutePath, fileContent, 0744)
	if err != nil {
		return
	}

	return
}

func spliceFileURL(isPublic bool, operatorID string) string {
	url := ""

	if isPublic {
		url += system_config.GetConfiguration().CloudFilePublicDir
	} else {
		url += operatorID
	}

	return mutils.FormatDirSuffix(url)
}

func spliceFileDir(isPublic bool, operatorID string) string {
	dir := mutils.FormatDirSuffix(system_config.GetConfiguration().CloudFileRootPath)
	dir += spliceFileURL(isPublic, operatorID)

	return mutils.FormatDirSuffix(dir)
}

// spliceFilePath return absolute path
func spliceFilePath(dir string, fileName string, extensionName string) string {
	return dir + fileName + "." + extensionName
}

func getValidBackupFileName(absolutePath string) string {
	absolutePath += backupFileSuffix

	for i := 0; ; i++ {
		fileInfo, err := os.Stat(absolutePath + strconv.Itoa(i))
		if !(err == nil && !fileInfo.IsDir()) {
			absolutePath += strconv.Itoa(i)
			break
		}
	}

	return absolutePath
}

func fileDBToHTTPRes(data *model.CloudFile) *structure.FileRes {
	if data == nil {
		return &structure.FileRes{}
	}

	return &structure.FileRes{
		FileID:           data.FileID,
		FileName:         data.FileName,
		LastModifiedTime: data.LastModifiedTime,
		FileURL:          spliceFilePath(spliceFileURL(data.IsPublic, data.UploadedBy), data.FileID, data.ExtensionName),
		IsPublic:         data.IsPublic,
		UpdateTime:       data.UpdateTime,
		CreatedTime:      data.CreatedTime,
	}
}
