package mhttp

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"net/http"
	"strings"
)

type ResponseData struct {
	UserID string `json:"userID,omitempty"`
	Token  string `json:"token,omitempty"`

	HasError bool   `json:"hasError"`
	Data     string `json:"data"`
}

// Response 'userID' is required when you want to refresh multi-login params
func Response(data interface{}, userID ...string) *ResponseData {
	res := &ResponseData{}
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

func ResponseWithError(errMsg ...string) *ResponseData {
	return &ResponseData{
		HasError: true,
		Data:     strings.Join(errMsg, ""),
	}
}

func response(writer http.ResponseWriter, data *ResponseData) {
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
