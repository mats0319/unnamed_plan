package mhttp

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http/response"
	"net/http"
	"strings"
)

// Response 'userID' is required when you want to refresh multi-login params
func Response(data interface{}, userID ...string) *mresponse.ResponseData {
	res := &mresponse.ResponseData{}
	if len(userID) > 0 {
		res.UserID = userID[0]
	}

	if data == mconst.EmptyHTTPRes {
		return res
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return ResponseWithError(err.Error())
	}

	res.Data = string(jsonBytes)

	return res
}

func ResponseWithError(errMsg ...string) *mresponse.ResponseData {
	return &mresponse.ResponseData{
		HasError: true,
		Data:     strings.Join(errMsg, ""),
	}
}

func response(writer http.ResponseWriter, data *mresponse.ResponseData) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(jsonBytes)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
