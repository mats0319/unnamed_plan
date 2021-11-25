import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
  //@ts-ignore-next-line
  baseURL: process.env.VUE_APP_axios_base_url as string,
  timeout: 3000
});

//@ts-ignore-next-line
let sign = process.env.VUE_APP_axios_unnamed_plan_sign as string;

export function initInterceptors(invalidLoginHandler: () => void): void {
  axiosWrapper.interceptors.request.use(
    (value: AxiosRequestConfig): AxiosRequestConfig => {
      value.headers.common["Unnamed-Plan"] = sign;
      value.headers.common["Unnamed-Plan-User"] = sessionStorage.getItem("user"); // 使用session storage防刷新
      value.headers.common["Unnamed-Plan-Token"] = sessionStorage.getItem("token");

      return value;
    }
  );

  axiosWrapper.interceptors.response.use(
    (value: AxiosResponse): AxiosResponse => {
      if (value.data && value.data.hasOwnProperty("userID")) {
        sessionStorage.setItem("user", value.data["userID"]);
      }
      if (value.data && value.data.hasOwnProperty("token")) {
        sessionStorage.setItem("token", value.data["token"]);
      }

      return value;
    },
    error => {
      let isInvalidLoginError = false;

      if (error.response && error.response.status && error.response.status === 401) {
        isInvalidLoginError = true;
        invalidLoginHandler();
      }

      return Promise.reject(isInvalidLoginError ? "Login Information is Invalid, Please Login again." : error)
    }
  )
}
