package handlers

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/config"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/http/response_type"
	"github.com/mats9693/unnamed_plan/admin_data/kits"
	"github.com/mats9693/utils/toy_server/http"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ListCloudFileByUploader(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(errorsToString(err, err2)))
		return
	}
	if len(operatorID) < 1 || pageSize < 1 || pageNum < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, page size: %d, page num: %d", operatorID, pageSize, pageNum)))
		return
	}

	files, count, err := dao.GetCloudFile().QueryPageByUploader(pageSize, pageNum, operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	filesRes := make([]*http_res_type.HTTPResFile, 0, len(files))
	for i := range files {
		url := ""
		if files[i].IsPublic {
			url = system_config.GetConfiguration().CloudFilePublicDir
		} else {
			url = files[i].UploadedBy
		}
		url = kits.AppendDirSuffix(url) + files[i].FileID + "." + files[i].ExtensionName

		filesRes = append(filesRes, &http_res_type.HTTPResFile{
			FileID:      files[i].FileID,
			FileName:    files[i].FileName,
			FileURL:     url,
			IsPublic:    files[i].IsPublic,
			UpdateTime:  files[i].UpdateTime,
			CreatedTime: files[i].CreatedTime,
		})
	}

	resData := &struct {
		Total int                          `json:"total"`
		Files []*http_res_type.HTTPResFile `json:"files"`
	}{
		Total: count,
		Files: filesRes,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func ListPublicCloudFile(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(errorsToString(err, err2)))
		return
	}
	if len(operatorID) < 1 || pageSize < 1 || pageNum < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, page size: %d, page num: %d", operatorID, pageSize, pageNum)))
		return
	}

	files, count, err := dao.GetCloudFile().QueryPageInPublic(pageSize, pageNum, operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	filesRes := make([]*http_res_type.HTTPResFile, 0, len(files))
	for i := range files {
		url := ""
		if files[i].IsPublic {
			url = system_config.GetConfiguration().CloudFilePublicDir
		} else {
			url = files[i].UploadedBy
		}
		url = kits.AppendDirSuffix(url) + files[i].FileID + "." + files[i].ExtensionName

		filesRes = append(filesRes, &http_res_type.HTTPResFile{
			FileID:      files[i].FileID,
			FileName:    files[i].FileName,
			FileURL:     url,
			IsPublic:    files[i].IsPublic,
			UpdateTime:  files[i].UpdateTime,
			CreatedTime: files[i].CreatedTime,
		})
	}

	resData := &struct {
		Total int                          `json:"total"`
		Files []*http_res_type.HTTPResFile `json:"files"`
	}{
		Total: count,
		Files: filesRes,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func UploadCloudFile(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	fileName := r.PostFormValue("fileName")
	extensionName := r.PostFormValue("extensionName")
	lastModifiedTime, err := strconv.Atoi(r.PostFormValue("lastModifiedTime"))
	isPublic, err2 := kits.StringToBool(r.PostFormValue("isPublic"))
	file, fileHeader, err3 := r.FormFile("file")

	if err != nil || err2 != nil || err3 != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(errorsToString(err, err2, err3)))
		return
	}
	defer file.Close()

	if len(operatorID) < 1 || len(fileName) < 1 || len(extensionName) < 1 || lastModifiedTime < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, file name: %s, extension name: %s, last modified time: %d", operatorID, fileName, extensionName, lastModifiedTime)))
		return
	}

	// make sure target directory structure is exist
	dir := kits.AppendDirSuffix(system_config.GetConfiguration().CloudFileRootPath)
	if isPublic {
		dir += kits.AppendDirSuffix(system_config.GetConfiguration().CloudFilePublicDir)
	} else {
		dir += kits.AppendDirSuffix(operatorID)
	}

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	// save file
	fileContent := make([]byte, fileHeader.Size) // require enough length before read
	_, err = file.Read(fileContent)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	fileID := kits.CalcSHA256(operatorID + time.Now().GoString())
	absolutePath := dir + fileID + "." + extensionName
	err = os.WriteFile(absolutePath, fileContent, 0744)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
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
		err2 = os.Remove(absolutePath)
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()+err2.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func ModifyCloudFile(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	fileID := r.PostFormValue("fileID")
	password := r.PostFormValue("password")
	fileName := r.PostFormValue("fileName")
	extensionName := r.PostFormValue("extensionName")
	isPublic, err := kits.StringToBool(r.PostFormValue("isPublic"))
	file, fileHeader, err2 := r.FormFile("file")
	lastModifiedTime, err3 := strconv.Atoi(r.PostFormValue("lastModifiedTime"))

	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}
	if len(operatorID) < 1 || len(fileID) < 1 || len(password) < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, file id: %s, password: %s", operatorID, fileID, password)))
		return
	}

	_, err = checkPwdByUserID(password, operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	fileRecord, err := dao.GetCloudFile().QueryFirst(model.CloudFile_FileID + " = ?", fileID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	if len(fileName) + len(extensionName) < 1 && isPublic == fileRecord.IsPublic && err2 != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, not any modification received")))
		return
	}

	if err2 == nil && err3 != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, parse file last modified time failed")))
		return
	}

	absolutePath := ""
	if err2 == nil {
		dir := kits.AppendDirSuffix(system_config.GetConfiguration().CloudFileRootPath)
		if isPublic {
			dir += kits.AppendDirSuffix(system_config.GetConfiguration().CloudFilePublicDir)
		} else {
			dir += kits.AppendDirSuffix(operatorID)
		}

		fileContent := make([]byte, fileHeader.Size) // require enough length before read
		_, err = file.Read(fileContent)
		if err != nil {
			_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
			return
		}

		absolutePath = dir + fileRecord.FileID + "." + extensionName
		err = os.WriteFile(absolutePath, fileContent, 0744)
		if err != nil {
			_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
			return
		}
	}

	updateColumns := make([]string, 0, 4)
	_ = updateColumns
	_ = lastModifiedTime
}

func DeleteCloudFile(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	password := r.PostFormValue("password")
	fileID := r.PostFormValue("fileID")

	if len(operatorID) < 1 || len(password) < 1 || len(fileID) < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, password: %s, file id: %s", operatorID, password, fileID)))
		return
	}

	_, err := checkPwdByUserID(password, operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	err = dao.GetCloudFile().UpdateColumnsByFileID(&model.CloudFile{
		FileID:    fileID,
		IsDeleted: true,
	}, model.CloudFile_IsDeleted)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	resData := &struct {
		IsSuccess bool `json:"isSuccess"`
	}{
		IsSuccess: true,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}
