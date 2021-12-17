import { axiosWrapper } from "./config";
import { calcSHA256 } from "../utils";

class CloudFileAxios {
  public upload(
    operatorID: string,
    fileName: string,
    extensionName: string,
    lastModifiedTime: number,
    isPublic: boolean,
    file: File
  ) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("fileName", fileName);
    data.append("extensionName", extensionName);
    data.append("lastModifiedTime", lastModifiedTime.toString())
    data.append("isPublic", isPublic.toString());
    data.append("file", file);
    return axiosWrapper.post("/api/cloudFile/upload", data);
  }

  public listByUploader(operatorID: string, pageSize: number, pageNum: number) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("pageSize", pageSize.toString());
    data.append("pageNum", pageNum.toString());
    return axiosWrapper.post("/api/cloudFile/listByUploader", data);
  }

  public listPublic(operatorID: string, pageSize: number, pageNum: number) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("pageSize", pageSize.toString());
    data.append("pageNum", pageNum.toString());
    return axiosWrapper.post("/api/cloudFile/listPublic", data);
  }

  public modify(
    operatorID: string,
    fileID: string,
    password: string,
    fileName: string,
    isPublic: boolean,
    lastModifiedTime: number,
    file?: File,
    extensionName?: string
  ) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("fileID", fileID);
    data.append("password", calcSHA256(password));
    data.append("fileName", fileName);
    data.append("isPublic", isPublic.toString());
    data.append("lastModifiedTime", lastModifiedTime.toString());
    if (file && extensionName) {
      data.append("file", file);
      data.append("extensionName", extensionName);
    }
    return axiosWrapper.post("/api/cloudFile/modify", data);
  }

  public delete(operatorID: string, fileID: string, password: string) {
    const data: FormData = new FormData();
    data.append("operatorID", operatorID);
    data.append("fileID", fileID);
    data.append("password", calcSHA256(password));
    return axiosWrapper.post("/api/cloudFile/delete", data);
  }
}

const cloudFileAxios = new CloudFileAxios();
export default cloudFileAxios;
