import { axiosWrapper } from "./config";
import { calcSHA256 } from "../utils";

class HomeAxios {
  public login(userName: string, password: string) {
    const data: FormData = new FormData();
    data.append("userName", userName);
    data.append("password", calcSHA256(password));
    return axiosWrapper.post("/api/login", data);
  }
}

const homeAxios = new HomeAxios();
export default homeAxios;
