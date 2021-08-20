package http

import "encoding/json"

type Response struct {
	HasError bool   `json:"hasError"`
	Data     string `json:"data"`
}

func response(data interface{}) string {
	res := &Response{}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return responseWithError(err.Error())
	}

	res.Data = string(jsonBytes)

	jsonBytes, err = json.Marshal(res)
	if err != nil {
		return err.Error()
	}

	return string(jsonBytes)
}

func responseWithError(errMsg string) string {
	res := &Response{
		HasError: true,
		Data:     errMsg,
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return err.Error()
	}

	return string(jsonBytes)
}
