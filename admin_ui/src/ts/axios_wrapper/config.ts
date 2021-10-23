import axios, { AxiosInstance } from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
  baseURL: process.env.VUE_APP_axios_base_url as string,
  timeout: 3000
});
