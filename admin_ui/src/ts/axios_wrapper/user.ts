import { axiosWrapper } from "@/ts/axios_wrapper/config";

class UserAxios {
  public create(operatorID: string, userName: string, password: string, permission: number) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("userName", userName);
    data.append("password", password);
    data.append("permission", permission.toString());
    return axiosWrapper.post("/api/user/create", data);
  }
}

const userAxios = new UserAxios();
export default userAxios
