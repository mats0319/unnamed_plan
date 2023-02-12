import { axiosWrapper } from "./config/config"
import { calcSHA256, objectToFormData } from "./utils"
import { User } from "./proto/1_user.pb"

class UserAxios {
  public login(userName: string, password: string): Promise<User.LoginRes> {
    let request: User.LoginReq = {
      user_name: userName,
      password: calcSHA256(password),
    }

    return axiosWrapper.post("/api/user/login", objectToFormData(request))
  }

  public create(
    operatorID: string,
    userName: string,
    password: string,
    permission: number,
    operatorPassword: string,
  ): Promise<User.CreateRes> {
    let request: User.CreateReq = {
      operator_id: operatorID,
      user_name: userName,
      password: calcSHA256(password),
      permission: permission,
      operator_password: operatorPassword,
    }

    return axiosWrapper.post("/api/user/create", objectToFormData(request))
  }

  public list(operatorID: string, pageSize: number, pageNum: number): Promise<User.ListRes> {
    let request: User.ListReq = {
      operator_id: operatorID,
      page: {
        page_size: pageSize,
        page_num: pageNum,
      },
    }
    
    return axiosWrapper.post("/api/user/list", objectToFormData(request))
  }

  public modifyInfo(operatorID: string, userID: string, currPwd: string, nickname: string, password: string): Promise<User.ModifyInfoRes> {
    let request: User.ModifyInfoReq = {
      operator_id: operatorID,
      user_id: userID,
      curr_pwd: calcSHA256(currPwd),
      nickname: nickname,
      password: calcSHA256(password),
    }
    
    return axiosWrapper.post("/api/user/modifyInfo", objectToFormData(request))
  }

  public lock(operatorID: string, userID: string, password: string): Promise<User.LockRes> {
    let request: User.LockReq = {
      operator_id: operatorID,
      user_id: userID,
      password: calcSHA256(password),
    }
    
    return axiosWrapper.post("/api/user/lock", objectToFormData(request))
  }

  public unlock(operatorID: string, userID: string, password: string): Promise<User.UnlockRes> {
    let request: User.UnlockReq = {
      operator_id: operatorID,
      user_id: userID,
      password: calcSHA256(password),
    }
    
    return axiosWrapper.post("/api/user/unlock", objectToFormData(request))
  }

  public modifyPermission(operatorID: string, userID: string, permission: number, password: string): Promise<User.ModifyPermissionRes> {
    let request: User.ModifyPermissionReq = {
      operator_id: operatorID,
      user_id: userID,
      permission: permission,
      password: calcSHA256(password),
    }

    return axiosWrapper.post("/api/user/modifyPermission", objectToFormData(request))
  }
}

const userAxios = new UserAxios()
export default userAxios
