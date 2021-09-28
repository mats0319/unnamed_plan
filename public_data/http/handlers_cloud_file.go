package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mats9693/unnamed_plan/admin_data/config"
	"github.com/mats9693/unnamed_plan/admin_data/db/dao"
	http_res_type "github.com/mats9693/unnamed_plan/admin_data/http/response_type"
	"github.com/mats9693/unnamed_plan/admin_data/kits"
	mhttp "github.com/mats9693/utils/toy_server/http"
)

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
