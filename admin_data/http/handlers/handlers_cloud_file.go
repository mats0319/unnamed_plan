package handlers

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/config"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/http/structure_defination"
	"github.com/mats9693/unnamed_plan/admin_data/utils"
	"github.com/mats9693/utils/toy_server/http"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

const backupFileSuffix = ".backup"

func ListCloudFileByUploader(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		return mhttp.ResponseWithError(errorsToString(err, err2))
	}
	if len(operatorID) < 1 || pageSize < 1 || pageNum < 1 {
		return mhttp.ResponseWithError(error_InvalidParams+
			fmt.Sprintf(", operator id: %s, page size: %d, page num: %d", operatorID, pageSize, pageNum))
	}

	files, count, err := dao.GetCloudFile().QueryPageByUploader(pageSize, pageNum, operatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	fileListRes := make([]*structure.FileListRes, 0, len(files))
	for i := range files {
		fileListRes = append(fileListRes, &structure.FileListRes{
			FileID:           files[i].FileID,
			FileName:         files[i].FileName,
			LastModifiedTime: files[i].LastModifiedTime,
			FileURL:          spliceFilePath(spliceFileDir(files[i].IsPublic, files[i].UploadedBy), files[i].FileID, files[i].ExtensionName),
			IsPublic:         files[i].IsPublic,
			UpdateTime:       files[i].UpdateTime,
			CreatedTime:      files[i].CreatedTime,
		})
	}

	return mhttp.Response(structure.MakeListCloudFileByUploaderRes(count, fileListRes))
}

func ListPublicCloudFile(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		return mhttp.ResponseWithError(errorsToString(err, err2))
	}
	if len(operatorID) < 1 || pageSize < 1 || pageNum < 1 {
		return mhttp.ResponseWithError(error_InvalidParams+
			fmt.Sprintf(", operator id: %s, page size: %d, page num: %d", operatorID, pageSize, pageNum))
	}

	files, count, err := dao.GetCloudFile().QueryPageInPublic(pageSize, pageNum, operatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	fileListRes := make([]*structure.FileListRes, 0, len(files))
	for i := range files {
		fileListRes = append(fileListRes, &structure.FileListRes{
			FileID:           files[i].FileID,
			FileName:         files[i].FileName,
			LastModifiedTime: files[i].LastModifiedTime,
			FileURL:          spliceFilePath(spliceFileDir(files[i].IsPublic, files[i].UploadedBy), files[i].FileID, files[i].ExtensionName),
			IsPublic:         files[i].IsPublic,
			UpdateTime:       files[i].UpdateTime,
			CreatedTime:      files[i].CreatedTime,
		})
	}

	return mhttp.Response(structure.MakeListPublicCloudFileRes(count, fileListRes))
}

func UploadCloudFile(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	fileName := r.PostFormValue("fileName")
	extensionName := r.PostFormValue("extensionName")
	lastModifiedTime, err := strconv.Atoi(r.PostFormValue("lastModifiedTime"))
	isPublic, err2 := utils.StringToBool(r.PostFormValue("isPublic"))
	file, fileHeader, err3 := r.FormFile("file")

	if err != nil || err2 != nil || err3 != nil {
		return mhttp.ResponseWithError(errorsToString(err, err2, err3))
	}
	defer file.Close()

	if len(operatorID) < 1 || len(fileName) < 1 || len(extensionName) < 1 || lastModifiedTime < 1 {
		return mhttp.ResponseWithError(error_InvalidParams+
			fmt.Sprintf(", operator id: %s, file name: %s, extension name: %s, last modified time: %d",
				operatorID, fileName, extensionName, lastModifiedTime))
	}

	// make sure target directory structure exist
	dir := spliceFileDir(isPublic, operatorID)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	// save file
	fileID := utils.CalcSHA256(operatorID + time.Now().GoString())
	absolutePath := spliceFilePath(dir, fileID, extensionName)

	err = saveFile(file, absolutePath, fileHeader.Size)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	// save db, if failed, remove file
	err = dao.GetCloudFile().Insert(&model.CloudFile{
		UploadedBy:       operatorID,
		FileID:           fileID,
		FileName:         fileName,
		ExtensionName:    extensionName,
		LastModifiedTime: time.Duration(lastModifiedTime),
		FileSize:         fileHeader.Size,
		IsPublic:         isPublic,
	})
	if err != nil {
		return mhttp.ResponseWithError(errorsToString(err, os.Remove(absolutePath)))
	}

	return mhttp.Response(structure.MakeUploadCloudFileRes(true))
}

func ModifyCloudFile(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	fileID := r.PostFormValue("fileID")
	password := r.PostFormValue("password")
	fileName := r.PostFormValue("fileName")
	extensionName := r.PostFormValue("extensionName")
	isPublic, err := utils.StringToBool(r.PostFormValue("isPublic"))
	file, fileHeader, err2 := r.FormFile("file")
	lastModifiedTime, err3 := strconv.Atoi(r.PostFormValue("lastModifiedTime"))

	// params check
	if err != nil || err2 != nil || err3 != nil {
		return mhttp.ResponseWithError(errorsToString(err, err2, err3))
	}
	defer file.Close()

	if len(operatorID) < 1 || len(fileID) < 1 || len(password) < 1 {
		return mhttp.ResponseWithError(error_InvalidParams+
			fmt.Sprintf(", operator id: %s, file id: %s, password: %s", operatorID, fileID, password))
	}

	_, err = verifyPwdByUserID(password, operatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	fileRecord, err := dao.GetCloudFile().QueryFirst(model.CloudFile_FileID+" = ?", fileID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	if len(fileName)+len(extensionName) < 1 && fileRecord.IsPublic == isPublic && err2 != nil {
		return mhttp.ResponseWithError(error_NoValidModification)
	}

	updateColumns := make([]string, 0, 5)

	// update file, backup old file in private folder
	if err2 == nil {
		// make sure private dir exist, only when old file is public
		if fileRecord.IsPublic {
			err = os.MkdirAll(spliceFileDir(false, operatorID), 0755)
			if err != nil {
				return mhttp.ResponseWithError(err.Error())
			}
		}

		// backup old file, may implicitly include move file
		oldPath := spliceFilePath(spliceFileDir(fileRecord.IsPublic, operatorID), fileRecord.FileID, extensionName)
		privPath := spliceFilePath(spliceFileDir(false, operatorID), fileRecord.FileID, extensionName)

		err = os.Rename(oldPath, getValidBackupFileName(privPath))
		if err != nil {
			return mhttp.ResponseWithError(err.Error())
		}

		// save new file
		absolutePath := spliceFilePath(spliceFileDir(isPublic, operatorID), fileRecord.FileID, extensionName)
		err = saveFile(file, absolutePath, fileHeader.Size)
		if err != nil {
			return mhttp.ResponseWithError(err.Error())
		}

		// update db record: file last modified time and file size
		fileRecord.LastModifiedTime = time.Duration(lastModifiedTime)
		updateColumns = append(updateColumns, model.CloudFile_LastModifiedTime)

		fileRecord.FileSize = fileHeader.Size
		updateColumns = append(updateColumns, model.CloudFile_FileSize)
	}

	if len(fileName) > 0 {
		fileRecord.FileName = fileName
		updateColumns = append(updateColumns, model.CloudFile_FileName)
	}
	if len(extensionName) > 0 && fileRecord.ExtensionName != extensionName {
		fileRecord.ExtensionName = extensionName
		updateColumns = append(updateColumns, model.CloudFile_ExtensionName)
	}
	if fileRecord.IsPublic != isPublic {
		fileRecord.IsPublic = isPublic
		updateColumns = append(updateColumns, model.CloudFile_IsPublic)
	}

	// update db record, if update failed and have new file, remove file
	err = dao.GetCloudFile().UpdateColumnsByFileID(fileRecord, updateColumns...)
	if err != nil {
		errMsg := err.Error()
		if err2 == nil {
			absolutePath := spliceFilePath(spliceFileDir(isPublic, operatorID), fileRecord.FileID, extensionName)
			errMsg += errorsToString(os.Remove(absolutePath))
		}
		return mhttp.ResponseWithError(errMsg)
	}

	return mhttp.Response(structure.MakeModifyCloudFileRes(true))
}

func DeleteCloudFile(r *http.Request) string {
	operatorID := r.PostFormValue("operatorID")
	password := r.PostFormValue("password")
	fileID := r.PostFormValue("fileID")

	if len(operatorID) < 1 || len(password) < 1 || len(fileID) < 1 {
		return mhttp.ResponseWithError(error_InvalidParams+
			fmt.Sprintf(", operator id: %s, password: %s, file id: %s", operatorID, password, fileID))
	}

	_, err := verifyPwdByUserID(password, operatorID)
	if err != nil {
		return mhttp.ResponseWithError(err.Error())
	}

	err = dao.GetCloudFile().UpdateColumnsByFileID(&model.CloudFile{
		FileID:    fileID,
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

func spliceFileDir(isPublic bool, operatorID string) string {
	dir := utils.AppendDirSuffix(system_config.GetConfiguration().CloudFileRootPath)
	if isPublic {
		dir += system_config.GetConfiguration().CloudFilePublicDir
	} else {
		dir += operatorID
	}

	return utils.AppendDirSuffix(dir)
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
