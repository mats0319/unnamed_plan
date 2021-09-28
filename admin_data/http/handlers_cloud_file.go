package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mats9693/unnamed_plan/admin_data/config"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/http/response_type"
	"github.com/mats9693/unnamed_plan/admin_data/kits"
	mhttp "github.com/mats9693/utils/toy_server/http"
)

func init() {
	root := system_config.GetConfiguration().CloudFileRootPath
	path := kits.AppendDirSuffix(root) + system_config.GetConfiguration().CloudFilePublicDir

	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println("os.MkdirAll failed, error:", err.Error())
		os.Exit(-1)
	}

	fmt.Println("> Cloud file directory init finish.")
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	fileName := r.PostFormValue("fileName")
	extensionName := r.PostFormValue("extensionName")
	isPublic, err := kits.StringToBool(r.PostFormValue("isPublic"))
	file, fileHeader, err2 := r.FormFile("file")

	if err != nil || err2 != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()+err2.Error()))
		return
	}
	defer file.Close()

	if len(operatorID) < 1 || len(fileName) < 1 || len(extensionName) < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, file name: %s, extension name: %s", operatorID, fileName, extensionName)))
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
	err = os.WriteFile(absolutePath, fileContent, 0755)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	// save db, if failed, remove file
	err = dao.GetCloudFile().Insert(&model.CloudFile{
		UploadedBy:    operatorID,
		FileID:        fileID,
		FileName:      fileName,
		ExtensionName: extensionName,
		FileSize:      strconv.FormatInt(fileHeader.Size, 10),
		IsPublic:      isPublic,
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

func listCloudFileByUploader(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()+err2.Error()))
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

	filesRes := make([]*http_res_type.HTTPResFiles, 0, len(files))
	for i := range files {
		url := ""
		if files[i].IsPublic {
			url = system_config.GetConfiguration().CloudFilePublicDir
		} else {
			url = files[i].UploadedBy
		}
		url = kits.AppendDirSuffix(url) + files[i].FileID + "." + files[i].ExtensionName

		filesRes = append(filesRes, &http_res_type.HTTPResFiles{
			FileID:      files[i].FileID,
			FileName:    files[i].FileName,
			FileURL:     url,
			IsPublic:    files[i].IsPublic,
			UpdateTime:  files[i].UpdateTime,
			CreatedTime: files[i].CreatedTime,
		})
	}

	resData := &struct {
		Total int                           `json:"total"`
		Files []*http_res_type.HTTPResFiles `json:"files"`
	}{
		Total: count,
		Files: filesRes,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func listPublicCloudFile(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()+err2.Error()))
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

	filesRes := make([]*http_res_type.HTTPResFiles, 0, len(files))
	for i := range files {
		url := ""
		if files[i].IsPublic {
			url = system_config.GetConfiguration().CloudFilePublicDir
		} else {
			url = files[i].UploadedBy
		}
		url = kits.AppendDirSuffix(url) + files[i].FileID + "." + files[i].ExtensionName

		filesRes = append(filesRes, &http_res_type.HTTPResFiles{
			FileID:      files[i].FileID,
			FileName:    files[i].FileName,
			FileURL:     url,
			IsPublic:    files[i].IsPublic,
			UpdateTime:  files[i].UpdateTime,
			CreatedTime: files[i].CreatedTime,
		})
	}

	resData := &struct {
		Total int                           `json:"total"`
		Files []*http_res_type.HTTPResFiles `json:"files"`
	}{
		Total: count,
		Files: filesRes,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func listDeleted(w http.ResponseWriter, r *http.Request) {
	if isDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	operatorID := r.PostFormValue("operatorID")
	pageSize, err := strconv.Atoi(r.PostFormValue("pageSize"))
	pageNum, err2 := strconv.Atoi(r.PostFormValue("pageNum"))

	if err != nil || err2 != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()+err2.Error()))
		return
	}
	if len(operatorID) < 1 || pageSize < 1 || pageNum < 1 {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, page size: %d, page num: %d", operatorID, pageSize, pageNum)))
		return
	}

	operator, err := dao.GetUser().QueryOne(model.User_UserID+" = ?", operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}
	if operator.Permission < system_config.GetConfiguration().SRankAdminPermission {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(fmt.Sprintf("permission denied, operator: %d, required: %d",
			operator.Permission, system_config.GetConfiguration().SRankAdminPermission)))
		return
	}

	files, count, err := dao.GetCloudFile().QueryPageInDeleted(pageSize, pageNum, operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}

	filesRes := make([]*http_res_type.HTTPResFiles, 0, len(files))
	for i := range files {
		url := ""
		if files[i].IsPublic {
			url = system_config.GetConfiguration().CloudFilePublicDir
		} else {
			url = files[i].UploadedBy
		}
		url = kits.AppendDirSuffix(url) + files[i].FileID + "." + files[i].ExtensionName

		filesRes = append(filesRes, &http_res_type.HTTPResFiles{
			FileID:      files[i].FileID,
			FileName:    files[i].FileName,
			FileURL:     url,
			IsPublic:    files[i].IsPublic,
			UpdateTime:  files[i].UpdateTime,
			CreatedTime: files[i].CreatedTime,
		})
	}

	resData := &struct {
		Total int                           `json:"total"`
		Files []*http_res_type.HTTPResFiles `json:"files"`
	}{
		Total: count,
		Files: filesRes,
	}

	_, _ = fmt.Fprintln(w, mhttp.Response(resData))

	return
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
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

	operator, err := dao.GetUser().QueryOne(model.User_UserID+" = ?", operatorID)
	if err != nil {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError(err.Error()))
		return
	}
	if operator.Password != kits.CalcSHA256(password, operator.Salt) {
		_, _ = fmt.Fprintln(w, mhttp.ResponseWithError("invalid account or password"))
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
