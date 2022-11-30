import {axiosWrapper} from "./config/config";
import { calcSHA256 } from "./utils/utils";

class TaskAxios {
  public create(operatorID: string, taskName: string, description: string, preTasks: Array<string>) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("taskName", taskName);
    data.append("description", description);
    data.append("preTaskIDs", preTasks.join(","));
    return axiosWrapper.post("/api/task/create", data);
  }

  public list(operatorID: string) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    return axiosWrapper.post("/api/task/list", data);
  }

  public modify(
    operatorID: string,
    taskID: string,
    password: string,
    taskName: string,
    description: string,
    preTaskIDs: Array<string>,
    taskStatus: number,
  ) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("taskID", taskID);
    data.append("password", calcSHA256(password));
    data.append("taskName", taskName);
    data.append("description", description);
    data.append("preTaskIDs", preTaskIDs.join(","));
    data.append("status", taskStatus.toString());
    return axiosWrapper.post("/api/task/modify", data);
  }
}

const taskAxios = new TaskAxios();
export default taskAxios;
