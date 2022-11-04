export interface Task {
  Data: Data,
  // ListReq: ListReq,
  // ListRule: ListRule,
}

enum ListRule {
  unSpecified,
  uploader,
  public,
}

interface Data {
  task_id: string,
}

interface ListReq {
  operator_id: string,
}

let taskIns: Task = {
  Data: {task_id: ""},
}

console.log(taskIns)
