package http

import "encoding/json"

type ResponseData struct {
	HasError bool   `json:"hasError"`
	Data     string `json:"data"`
}

func Response(data interface{}) string {
	res := &ResponseData{}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return ResponseWithError(err.Error())
	}

	res.Data = string(jsonBytes)

	jsonBytes, err = json.Marshal(res)
	if err != nil {
		return err.Error()
	}

	return string(jsonBytes)
}

func ResponseWithError(errMsg string) string {
	res := &ResponseData{
		HasError: true,
		Data:     errMsg,
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return err.Error()
	}

	return string(jsonBytes)
}
