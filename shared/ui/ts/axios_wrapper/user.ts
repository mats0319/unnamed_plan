import { axiosWrapper } from "./config";
import { calcSHA256 } from "../utils";

class UserAxios {
  public create(operatorID: string, userName: string, password: string, permission: number) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("userName", userName);
    data.append("password", password);
    data.append("permission", permission.toString());
    return axiosWrapper.post("/api/user/create", data);
  }

  public modifyInfo(operatorID: string, userID: string, currPwd: string, nickname: string, password: string) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("userID", userID);
    data.append("currPwd", calcSHA256(currPwd));
    data.append("nickname", nickname);
    data.append("password", calcSHA256(password))
    return axiosWrapper.post("/api/user/modifyInfo", data);
  }

  public list(operatorID: string, pageSize: number, pageNum: number) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("pageSize", pageSize.toString());
    data.append("pageNum", pageNum.toString());
    return axiosWrapper.post("/api/user/list", data);
  }

  public lock(operatorID: string, userID: string) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("userID", userID);
    return axiosWrapper.post("/api/user/lock", data);
  }

  public unlock(operatorID: string, userID: string) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("userID", userID);
    return axiosWrapper.post("/api/user/unlock", data);
  }

  public modifyPermission(operatorID: string, userID: string, permission: number) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("userID", userID);
    data.append("permission", permission.toString());
    return axiosWrapper.post("/api/user/modifyPermission", data);
  }
}

const userAxios = new UserAxios();
export default userAxios
