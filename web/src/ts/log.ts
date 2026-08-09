import { ElNotification } from "element-plus"

// Log print to console and display in webpage
class Log {
    success(behavior: string): void {
        console.log("> " + behavior + " Success.")
        ElNotification({ title: behavior, type: "success", customClass: "color-bg-1" })
    }

    fail(behavior: string, errStr: string = ""): void {
        console.log("> " + behavior + " Failed.", errStr)
        ElNotification({ title: behavior, message: errStr, type: "error", customClass: "color-bg-1" })
    }

    httpError(code: number): void {
        let str = "Unknown HTTP Status Code: " + code
        if (code == 401) {
            str = "Unauthorized"
        } else if (code == 500) {
            str = "Server Internal Error"
        } else if (code == 400) {
            str = "Bad Request"
        }

        console.log("> " + str + ".")
        ElNotification({ title: str, type: "error", customClass: "color-bg-1" })
    }
}

export const log = new Log()
