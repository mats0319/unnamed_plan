package http

import "net/http"

func uploadFile(w http.ResponseWriter, r *http.Request) {
    if isDev {
        w.Header().Set("Access-Control-Allow-Origin", "*")
    }

    //operatorID := r.PostFormValue("operatorID")
    //fileName := r.PostFormValue("fileName")
    //extensionName := r.PostFormValue("extensionName")
    //isPublic := r.PostFormValue("isPublic")
    //file, fileHeader, err := r.FormFile("file")

    // params check
    // generate directory structure if not exist, put into utils
    // save file
    // save db
}
