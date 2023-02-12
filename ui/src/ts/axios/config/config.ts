import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
  baseURL: import.meta.env.Vite_axios_base_url as string,
  timeout: 3000
});

let sign = import.meta.env.Vite_axios_source_sign as string;

export function initInterceptors(invalidLoginHandler: () => void): void {
  axiosWrapper.interceptors.request.use(
    (value: AxiosRequestConfig): AxiosRequestConfig => {
      // value.headers.type: AxiosHeader
      //@ts-ignore
      value.headers.set("Unnamed-Plan-Source", sign)
      //@ts-ignore
      value.headers.set("Unnamed-Plan-User",sessionStorage.getItem("user")) // 使用session storage防刷新
      //@ts-ignore
      value.headers.set("Unnamed-Plan-Token", sessionStorage.getItem("token"))

      return value;
    }
  );

  axiosWrapper.interceptors.response.use(
    (value: AxiosResponse) => {
      //@ts-ignore
      const userID = value.headers.get("Unnamed-Plan-User-Res") as string
      if (userID && userID.length > 0) {
          sessionStorage.setItem("user", userID)
      }

      //@ts-ignore
      const token = value.headers.get("Unnamed-Plan-Token-Res") as string
      if (token && token.length > 0) {
          sessionStorage.setItem("token", token)
      }

      return value.data
    },
    (error: any) => {
      let isInvalidLoginError = false

      if (error.response && error.response.status && error.response.status === 401) {
        isInvalidLoginError = true
        invalidLoginHandler()
      }

      return Promise.reject(isInvalidLoginError ? "Login Information is Invalid, Please Login again." : error)
    }
  )
}
