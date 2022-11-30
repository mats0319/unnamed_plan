import { axiosWrapper } from "./config/config";
import { calcSHA256 } from "./utils/utils";
import { User } from "./proto/1_user.pb";

class UserAxios {
  public login(userName: string, password: string) {
    let request: User.LoginReq = {
      user_name: userName,
      password: calcSHA256(password),
    };
    return axiosWrapper.post("/api/user/login", request);
  }

  public create(operatorID: string, userName: string, password: string, permission: number) {
    let request: User.CreateReq = {
      operator_id: operatorID,
      user_name: userName,
      password: calcSHA256(password),
      permission: permission,
    }
    return axiosWrapper.post("/api/user/create", request);
  }

  public modifyInfo(operatorID: string, userID: string, currPwd: string, nickname: string, password: string) {
    let request: User.ModifyInfoReq = {
      operator_id: operatorID,
      user_id: userID,
      curr_pwd: calcSHA256(currPwd),
      nickname: nickname,
      password: calcSHA256(password),
    }
    return axiosWrapper.post("/api/user/modifyInfo", request);
  }

  public list(operatorID: string, pageSize: number, pageNum: number) {
    let request: User.ListReq = {
      operator_id: operatorID,
      page: {
        page_size: pageSize,
        page_num: pageNum,
      },
    }
    return axiosWrapper.post("/api/user/list", request);
  }

  public lock(operatorID: string, userID: string) {
    let request: User.LockReq = {
      operator_id: operatorID,
      user_id: userID,
    }
    return axiosWrapper.post("/api/user/lock", request);
  }

  public unlock(operatorID: string, userID: string) {
    let request: User.UnlockReq = {
      operator_id: operatorID,
      user_id: userID,
    }
    return axiosWrapper.post("/api/user/unlock", request);
  }

  public modifyPermission(operatorID: string, userID: string, permission: number) {
    let request: User.ModifyPermissionReq = {
      operator_id: operatorID,
      user_id: userID,
      permission: permission,
    }
    return axiosWrapper.post("/api/user/modifyPermission", request);
  }
}

const userAxios = new UserAxios();
export default userAxios
