import { axiosWrapper } from "@/ts/axios_wrapper/config";

class HomeAxios {
  public login(userName: string, password: string) {
    const data: FormData = new FormData();
    data.append("userName", userName);
    data.append("password", password);
    return axiosWrapper.post("/api/login", data);
  }
}

const homeAxios = new HomeAxios();
export default homeAxios;
