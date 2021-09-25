package http

import (
    "fmt"
    system_config "github.com/mats9693/unnamed_plan/admin_data/config"
    "github.com/mats9693/unnamed_plan/admin_data/kits"
    shttp "github.com/mats9693/utils/toy_server/http"
    "net/http"
    "os"
)

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
        _, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()+err2.Error()))
        return
    }
    if len(operatorID) < 1 || len(fileName) < 1 || len(extensionName) < 1 {
        _, _ = fmt.Fprintln(w, shttp.ResponseWithError(fmt.Sprintf("invalid params, operator id: %s, file name: %s, extension name: %s", operatorID, fileName, extensionName)))
        return
    }

    path := system_config.GetConfiguration().CloudFileRootPath + "/"
    if isPublic {
        path += system_config.GetConfiguration().CloudFilePublicDir + "/"
    } else {
        path += operatorID + "/"
    }

    err = os.MkdirAll(path, 0755)
    if err != nil {
        _, _ = fmt.Fprintln(w, shttp.ResponseWithError(err.Error()))
        return
    }
    // save file
    _ = file
    _ = fileHeader
    // save db, if failed, remove file
}
