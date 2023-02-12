package mhttp

import (
	"encoding/json"
	"net/http"
)

type ResponseData struct {
	Data []byte `json:"data"`

	Success        bool   `json:"-"`
	AfterHookParam string `json:"-"`
}

func NewResponseData(data interface{}, afterHookParam ...string) *ResponseData {
	res := formatData(data)

	res.Success = true
	if len(afterHookParam) > 0 {
		res.AfterHookParam = afterHookParam[0]
	}

	return res
}

func NewResponseDataWithError(data interface{}) *ResponseData {
	return formatData(data)
}

func response(writer http.ResponseWriter, data []byte) {
	_, err := writer.Write(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func formatData(data interface{}) *ResponseData {
	res := &ResponseData{}

	jsonBytes, err := json.Marshal(data)
	if err == nil {
		res.Data = jsonBytes
	}

	return res
}
