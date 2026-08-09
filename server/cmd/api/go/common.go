package api

type Pagination struct {
	Size int `json:"size"`
	Num  int `json:"num"`
}

// Response 写给ts使用，等gocts支持导入其他包的数据后，这里改成*mhttp.Response
type Response struct {
	IsSuccess bool   `json:"is_success"`
	Code      int    `json:"code"`
	Err       string `json:"err"`
	Data      any    `json:"data"`
}
