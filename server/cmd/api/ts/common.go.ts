// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.5

export class Pagination {
    size: number = 0
    num: number = 0
}

// Response 写给ts使用，等gocts支持导入其他包的数据后，这里改成*mhttp.Response
export class Response {
    is_success: boolean = false
    code: number = 0
    err: string = ""
    data: Object = {}
}
